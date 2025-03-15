package internal

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"upload-gateway/internal/models"
	"upload-gateway/internal/repository"
)

type UploadServerSelector interface {
	GetServers(ctx context.Context) ([]*models.Server, error)
}

type FileNameGenerator interface {
	Generate() (string, error)
}

type StreamSender interface {
	Send(data []byte, partNumber int, fileName string) error
	OpenConn(ctx context.Context, server string) error
	CloseAndRecv() (FileChunkUploadResponse, error)
	Close() error
}

type StreamSenderFactory func() StreamSender

type Uploader struct {
	serverChooser               UploadServerSelector
	fileNameGenerator           FileNameGenerator
	streamSenderFactory         StreamSenderFactory
	uploadedFileRepository      *repository.UploadedFileRepository
	uploadedFileChunkRepository *repository.UploadedFileChunkRepository
}

func NewUploader(
	serverChooser UploadServerSelector,
	fileNameGenerator FileNameGenerator,
	streamSenderFactory StreamSenderFactory,
	uploadedFileRepository *repository.UploadedFileRepository,
	uploadedFileChunkRepository *repository.UploadedFileChunkRepository,
) *Uploader {
	return &Uploader{
		serverChooser:               serverChooser,
		fileNameGenerator:           fileNameGenerator,
		streamSenderFactory:         streamSenderFactory,
		uploadedFileRepository:      uploadedFileRepository,
		uploadedFileChunkRepository: uploadedFileChunkRepository,
	}
}

func (u *Uploader) UploadFileByChunks(ctx context.Context, file io.ReadCloser, fileName string) (string, error) {
	ext := filepath.Ext(fileName)
	originalName := fileName[:len(fileName)-len(ext)]

	uplFile := &models.UploadedFile{
		Ext:          filepath.Ext(fileName),
		OriginalName: originalName,
	}

	tempFileName, err := u.createTempFile(file)
	if err != nil {
		return "", err
	}

	uplFile.GeneratedName = tempFileName

	if err = u.uploadedFileRepository.Save(ctx, uplFile); err != nil {
		return "", fmt.Errorf("error to save uplodaded file to d: %w", err)
	}

	tempFile, err := os.Open(tempFileName)
	if err != nil {
		return "", fmt.Errorf("error to open temp file: %w", err)
	}
	defer func() {
		tempFile.Close()
		err = os.Remove(tempFileName)
		if err != nil {
			log.Printf("error to remove file %s", tempFileName)
		}
	}()

	servers, err := u.serverChooser.GetServers(ctx)
	if err != nil {
		return "", fmt.Errorf("error to select servers: %w", err)
	}

	chunks := int64(len(servers))
	if chunks == 0 {
		return "", fmt.Errorf("no available servers")
	}

	stats, err := tempFile.Stat()
	if err != nil {
		return "", fmt.Errorf("error to get file stats: %w", err)
	}
	partSize := stats.Size() / chunks
	remainder := stats.Size() % chunks

	wg := sync.WaitGroup{}
	for i := 0; i < len(servers); i++ {
		currentPartSize := partSize
		if int64(i) == chunks-1 {
			currentPartSize += remainder
		}

		chunkFileName, err := u.createFileChunk(tempFile, i, currentPartSize)
		if err != nil {
			return "", fmt.Errorf("error to create file part: %w", err)
		}
		_ = chunkFileName

		uplChunk := &models.UploadedFileChunk{
			UploadedFileID: uplFile.ID,
			ChunkNumber:    uint(i),
			Name:           chunkFileName,
			Size:           uint64(currentPartSize),
			ServerID:       servers[i].ID,
		}
		if err = u.uploadedFileChunkRepository.Save(ctx, uplChunk); err != nil {
			return "", fmt.Errorf("error to save file chunk: %w", err)
		}

		wg.Add(1)
		go func(counter int, chunkFileName string, uplChunk *models.UploadedFileChunk) {
			defer wg.Done()
			if err := u.sendFileChunkFromPath(ctx, chunkFileName, servers[counter], uplChunk, counter); err != nil {
				log.Printf("Ошибка загрузки файла %s на сервер %s: %v", chunkFileName, servers[counter], err)
			} else {
				log.Printf("Файл %s успешно загружен на %s", chunkFileName, servers[counter])
			}
		}(i, chunkFileName, uplChunk)
	}

	wg.Wait()

	return tempFileName, nil
}

func (u *Uploader) sendFileChunkFromPath(ctx context.Context, chunkFileName string, server *models.Server, uplFileChunk *models.UploadedFileChunk, chunkNumber int) error {
	file, err := os.Open(chunkFileName)
	if err != nil {
		return fmt.Errorf("error to open file: %w", err)
	}
	defer func() {
		file.Close()

		if err = os.Remove(chunkFileName); err != nil {
			log.Println(fmt.Errorf("error to remove file chunk: %w", err).Error())
		}
	}()

	hash, err := u.sendFileChunkByStream(ctx, file, server, chunkNumber, chunkFileName)
	if err != nil {
		return fmt.Errorf("error to send file chunk: %w", err)
	}

	uplFileChunk.Hash = hash
	if err = u.uploadedFileChunkRepository.Save(ctx, uplFileChunk); err != nil {
		return fmt.Errorf("error to save file chunk: %w", err)
	}

	return nil
}

func (u *Uploader) sendFileChunkByStream(ctx context.Context, file io.ReadCloser, server *models.Server, chunkNumber int, fileName string) (string, error) {
	const chunkSize = 1024 * 50
	buffer := make([]byte, chunkSize)

	sender := u.streamSenderFactory()
	err := sender.OpenConn(ctx, server.URL)
	if err != nil {
		return "", fmt.Errorf("error to connect to server %s: %w", server.URL, err)
	}
	defer func() {
		err = sender.Close()
		if err != nil {
			log.Printf("error ro close conn: %v", err)
		}
	}()

	hasher := md5.New()

	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}

			return "", fmt.Errorf("error to read file: %w", err)
		}

		err = sender.Send(buffer[:n], chunkNumber, fileName)
		if err != nil {
			return "", fmt.Errorf("error to send file to server %s: %w", server, err)
		}

		_, err = hasher.Write(buffer[:n])
		if err != nil {
			return "", fmt.Errorf("error to write bytes to hasher: %w", err)
		}
	}

	res, err := sender.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("error to close conn: %w", err)
	}

	calcHash := hex.EncodeToString(hasher.Sum(nil))
	if calcHash != res.Hash {
		return "", fmt.Errorf("client hash %s is not equal to server hash %s", calcHash, res.Hash)
	}

	return calcHash, nil
}

func (u *Uploader) createFileChunk(file *os.File, chunkNum int, chunkSize int64) (string, error) {
	partFileName := file.Name() + "_" + strconv.Itoa(chunkNum) + ".dat"
	fileChunk, err := os.Create(partFileName)
	if err != nil {
		return "", fmt.Errorf("error to create file part: %w", err)
	}
	defer fileChunk.Close()

	bytesCopied, err := io.CopyN(fileChunk, file, chunkSize)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error to copy file part: %w", err)
	}

	log.Printf("Создан файл %s, скопировано %d байт\n", partFileName, bytesCopied)

	return partFileName, nil
}

func (u *Uploader) createTempFile(file io.ReadCloser) (string, error) {
	name, err := u.fileNameGenerator.Generate()
	if err != nil {
		return "", fmt.Errorf("error to generate file name: %w", err)
	}

	tmpFile, err := os.Create(name)
	if err != nil {
		return "", fmt.Errorf("error to create temp file: %w", err)
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		return "", fmt.Errorf("error to copy original file to temp file: %w", err)
	}

	return name, nil
}
