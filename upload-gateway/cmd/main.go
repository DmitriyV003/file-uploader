package main

import (
	"fmt"
	"log"
	"net/http"
	"upload-gateway/app"
	"upload-gateway/internal"
	"upload-gateway/internal/http_interface"
	"upload-gateway/internal/repository"
)

func main() {
	app.InitDb()

	serverRepository := repository.NewServerRepository(app.GetDb())
	uploadedFileRepository := repository.NewUploadedFileRepository(app.GetDb())
	uploadedFileChunkRepository := repository.NewUploadedFileChunkRepository(app.GetDb())

	serverSelector := internal.NewDefaultServerSelector(serverRepository)
	fileNameGenerator := internal.NewDefaultFileNameNameGenerator()
	uploader := internal.NewUploader(serverSelector, fileNameGenerator, internal.GRPCStreamSender, uploadedFileRepository, uploadedFileChunkRepository)

	uh := http_interface.NewUploaderHandler(uploader)
	dh := http_interface.NewDownloadHandler(uploadedFileRepository, internal.NewGRPCChannelDownloader())

	http.HandleFunc("/upload", uh.Handle)
	http.HandleFunc("/", dh.Handle)

	fmt.Println("Сервер запущен на порту :7000")
	log.Fatal(http.ListenAndServe(":7000", nil))
}
