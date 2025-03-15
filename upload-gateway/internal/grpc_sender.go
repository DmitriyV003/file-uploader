package internal

import (
	"context"
	uploadpb "github.com/dmitriyV003/platform/proto"
	"google.golang.org/grpc"
)

type FileChunkUploadResponse struct {
	Code int32
	Hash string
}

type GRPCSender struct {
	conn   *grpc.ClientConn
	client uploadpb.UploadFileApiClient
	stream uploadpb.UploadFileApi_UploadFileClient
}

func newGRPCSender() *GRPCSender {
	return &GRPCSender{}
}

func (g *GRPCSender) Send(data []byte, partNumber int, fileName string) error {
	if g.stream != nil {
		if err := g.stream.Send(&uploadpb.UploadFileRequest{
			FileId:     fileName,
			PartNumber: int32(partNumber),
			Data:       data,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (g *GRPCSender) OpenConn(ctx context.Context, server string) error {
	if g.conn != nil {
		return nil
	}

	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return err
	}

	g.conn = conn
	g.client = uploadpb.NewUploadFileApiClient(g.conn)

	stream, err := g.client.UploadFile(ctx)
	if err != nil {
		g.conn.Close()
		g.conn = nil

		return err
	}
	g.stream = stream

	return nil
}

func (g *GRPCSender) CloseAndRecv() (FileChunkUploadResponse, error) {
	var (
		res FileChunkUploadResponse
	)

	if g.stream != nil {
		resp, err := g.stream.CloseAndRecv()
		g.stream = nil
		if err != nil {
			return FileChunkUploadResponse{}, err
		}

		res.Hash = resp.Hash
		res.Code = resp.Code
	}
	if g.conn != nil {
		errClose := g.conn.Close()
		g.conn = nil
		if errClose != nil {
			return FileChunkUploadResponse{}, errClose
		}
	}

	return res, nil
}
func (g *GRPCSender) Close() error {
	if g.stream != nil {
		err := g.stream.CloseSend()
		g.stream = nil
		if err != nil {
			return err
		}
	}
	if g.conn != nil {
		errClose := g.conn.Close()
		g.conn = nil
		if errClose != nil {
			return errClose
		}
	}

	return nil
}

func GRPCStreamSender() StreamSender {
	return newGRPCSender()
}
