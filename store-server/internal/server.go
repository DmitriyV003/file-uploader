package internal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	uploadpb "github.com/dmitriyV003/platform/proto"
	"io"
	"log"
	"os"
	"sync"
)

type Server struct {
	uploadpb.UnimplementedUploadFileApiServer
}

func (s *Server) UploadFile(stream uploadpb.UploadFileApi_UploadFileServer) error {
	var file *os.File
	var partNumber int32
	var chunkName string
	var once sync.Once

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	hasher := md5.New()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ошибка получения чанка: %v", err)
		}

		once.Do(func() {
			chunkName = chunk.FileId
			partNumber = chunk.PartNumber
		})

		if file == nil {
			file, err = os.Create(chunk.FileId)
			if err != nil {
				return fmt.Errorf("ошибка создания файла %s: %v", chunk.FileId, err)
			}
		}

		if _, err := file.Write(chunk.Data); err != nil {
			return fmt.Errorf("ошибка записи данных: %v", err)
		}

		_, err = hasher.Write(chunk.Data)
		if err != nil {
			return fmt.Errorf("error to write to hasher: %w", err)
		}
	}

	log.Printf("file chunk saved: %s, chunk number: %d", chunkName, partNumber)

	return stream.SendAndClose(&uploadpb.UploadFileResponse{
		Code: 200,
		Hash: hex.EncodeToString(hasher.Sum(nil)),
	})
}

func (s *Server) GetFileChunk(req *uploadpb.FileRequest, stream uploadpb.UploadFileApi_GetFileChunkServer) error {
	file, err := os.Open(req.GetFileName())
	if err != nil {
		return fmt.Errorf("error to open file %s: %w", req.GetFileName(), err)
	}
	defer file.Close()

	const chunkSize = 1024 * 50
	buffer := make([]byte, chunkSize)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error to read file: %w", err)
		}
		if n == 0 {
			break
		}

		resp := &uploadpb.FileChunk{
			Data: buffer[:n],
		}

		if err := stream.Send(resp); err != nil {
			return fmt.Errorf("error to send data: %w", err)
		}
	}

	return nil
}
