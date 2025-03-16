package http_interface

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"upload-gateway/internal"
	"upload-gateway/internal/repository"
)

type DownloadHandler struct {
	uploadedFileRepository *repository.UploadedFileRepository
	grpcChannelDownloader  *internal.GRPCChannelDownloader
}

func NewDownloadHandler(
	uploadedFileRepository *repository.UploadedFileRepository,
	grpcChannelDownloader *internal.GRPCChannelDownloader,
) *DownloadHandler {
	return &DownloadHandler{
		uploadedFileRepository: uploadedFileRepository,
		grpcChannelDownloader:  grpcChannelDownloader,
	}
}

func (u *DownloadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)

		return
	}

	fileName := strings.TrimPrefix(r.URL.Path, "/")
	if fileName == "" {
		http.Error(w, "fileName is required", http.StatusBadRequest)

		return
	}

	fileWithChunks, err := u.uploadedFileRepository.GetByGeneratedName(r.Context(), fileName)
	if err != nil || len(fileWithChunks.UploadedFileChunks) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to get file"))

		log.Printf("error to get file: %v", err)

		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(fileWithChunks.OriginalName+fileWithChunks.Ext)))

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming does not support", http.StatusInternalServerError)

		return
	}

	ch := make(chan []byte, 50*1024)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for bytes := range ch {
			if _, errWrite := w.Write(bytes); errWrite != nil {
				log.Printf("error to write data: %v", errWrite)

				return
			}
			flusher.Flush()
		}
	}()

	for _, chunk := range fileWithChunks.UploadedFileChunks {
		err := u.grpcChannelDownloader.GetFileChunk(r.Context(), ch, chunk.Name, chunk.Server.URL, chunk.Hash)
		if err != nil {
			log.Printf("error to open file %s: %v", chunk.Name, err)

			return
		}
	}

	close(ch)
	wg.Wait()
}
