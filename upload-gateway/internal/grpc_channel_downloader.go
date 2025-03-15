package internal

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	uploadpb "github.com/dmitriyV003/platform/proto"
	"google.golang.org/grpc"
	"log"
)

type GRPCChannelDownloader struct {
	conn   *grpc.ClientConn
	client uploadpb.UploadFileApiClient
}

func NewGRPCChannelDownloader() *GRPCChannelDownloader {
	return &GRPCChannelDownloader{}
}

func (g *GRPCChannelDownloader) GetFileChunk(ctx context.Context, outChan chan []byte, fileChunkName, server, hash string) error {
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

	hasher := md5.New()

	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}

		outChan <- resp.Data
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
