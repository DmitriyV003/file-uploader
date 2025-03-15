package internal

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	uploadpb "github.com/dmitriyV003/platform/proto"
	"google.golang.org/grpc"
	"log"
	"os"
)

type GRPCDownloader struct {
	conn   *grpc.ClientConn
	client uploadpb.UploadFileApiClient
}

func newGRPCDownloader() *GRPCDownloader {
	return &GRPCDownloader{}
}

func GRPCStreamDownloader() StreamDownloader {
	return newGRPCDownloader()
}

func (g *GRPCDownloader) GetFileChunk(ctx context.Context, fileChunkName, server, hash string) error {
	if g.conn != nil {
		return nil
	}

	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return err
	}

	g.conn = conn
	g.client = uploadpb.NewUploadFileApiClient(g.conn)

	stream, err := g.client.GetFileChunk(ctx, &uploadpb.FileRequest{FileName: fileChunkName})
	if err != nil {
		g.conn.Close()
		g.conn = nil

		return err
	}

	file, err := os.Create(fileChunkName)
	if err != nil {
		return fmt.Errorf("error to create flie chunk: %w", err)
	}
	defer file.Close()

	hasher := md5.New()

	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		n, err := file.Write(resp.Data)
		if err != nil {
			return fmt.Errorf("error fo write bytes to file: %w", err)
		}
		fmt.Printf("Получено: %d для чанка файла %s \n", n, fileChunkName)

		hasher.Write(resp.Data)
	}

	caclHash := hex.EncodeToString(hasher.Sum(nil))
	log.Printf("HASH: %t \n", caclHash == hash)

	fmt.Printf("Получение %s файла окончено \n", fileChunkName)

	if g.conn != nil {
		errClose := g.conn.Close()
		g.conn = nil
		if errClose != nil {
			return errClose
		}
	}

	return nil
}
