package http_interface

import (
	"fmt"
	"log"
	"net/http"
	"upload-gateway/internal"
)

type UploadHandler struct {
	uploader *internal.Uploader
}

func NewUploaderHandler(uploader *internal.Uploader) *UploadHandler {
	return &UploadHandler{
		uploader: uploader,
	}
}

func (u *UploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)

		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Ошибка разбора формы: "+err.Error(), http.StatusBadRequest)

		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка получения файла: "+err.Error(), http.StatusBadRequest)

		return
	}
	log.Printf("Received file %s with size %d byte\n", fileHeader.Filename, fileHeader.Size)
	defer file.Close()

	generatedName, err := u.uploader.UploadFileByChunks(r.Context(), file, fileHeader.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("error to upload file: %w", err).Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("File uploaded: %s", generatedName)))
}
