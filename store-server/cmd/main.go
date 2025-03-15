package main

import (
	"fmt"
	uploadpb "github.com/dmitriyV003/platform/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"store-server/internal"
)

func main() {
	server := internal.Server{}
	port := os.Getenv("PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
	grpcServer := grpc.NewServer()
	uploadpb.RegisterUploadFileApiServer(grpcServer, &server)
	reflection.Register(grpcServer)

	log.Printf("Сервер запущен на %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы сервера: %v", err)
	}
}
