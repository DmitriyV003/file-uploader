package internal

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"upload-gateway/internal/models"
	"upload-gateway/internal/repository"
)

type Downloader struct {
	uploadedFileRepository *repository.UploadedFileRepository
	streamDownloader       StreamDownloaderFactory
}

type StreamDownloaderFactory func() StreamDownloader

type StreamDownloader interface {
	GetFileChunk(ctx context.Context, fileChunkName, server, hash string) error
}

func NewDownloader(
	uploadedFileRepository *repository.UploadedFileRepository,
	streamDownloader StreamDownloaderFactory,
) *Downloader {
	return &Downloader{
		uploadedFileRepository: uploadedFileRepository,
		streamDownloader:       streamDownloader,
	}
}

func (d *Downloader) DownloadFileChunks(ctx context.Context, generatedFileName string) (*models.UploadedFile, error) {
	fileWithChunks, err := d.uploadedFileRepository.GetByGeneratedName(ctx, generatedFileName)
	if err != nil {
		return nil, fmt.Errorf("error to get file by generated name: %w", err)
	}

	if len(fileWithChunks.UploadedFileChunks) == 0 {
		return nil, fmt.Errorf("file has no chunks")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	for _, fch := range fileWithChunks.UploadedFileChunks {
		fch := fch
		eg.Go(func() error {
			downloader := d.streamDownloader()

			return downloader.GetFileChunk(ctx, fch.Name, fch.Server.URL, fch.Hash)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("error to get file chunk: %w", err)
	}

	return fileWithChunks, nil
}
