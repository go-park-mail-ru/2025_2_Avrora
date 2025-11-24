package service

import (
	"context"
	"io"

	fileserverpb "github.com/go-park-mail-ru/2025_2_Avrora/proto/fileserver"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type FileServerClient struct {
	client fileserverpb.FileServerClient
	logger *log.Logger
}

func NewFileServerClient(addr string, logger *log.Logger) (*FileServerClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	
	return &FileServerClient{
		client: fileserverpb.NewFileServerClient(conn),
		logger: logger.With(zap.String("service", "fileserver_client")),
	}, nil
}

func (c *FileServerClient) Upload(ctx context.Context, data []byte, filename, contentType string) (string, error) {
	req := &fileserverpb.UploadRequest{
		Data:        data,
		Filename:    filename,
		ContentType: contentType,
	}
	
	resp, err := c.client.Upload(ctx, req)
	if err != nil {
		c.logger.Error(ctx, "upload failed", zap.Error(err), zap.String("filename", filename))
		return "", err
	}
	
	return resp.Url, nil
}

func (c *FileServerClient) Get(ctx context.Context, filename string) (io.ReadCloser, string, error) {
	req := &fileserverpb.GetRequest{
		Filename: filename,
	}
	
	stream, err := c.client.Get(ctx, req)
	if err != nil {
		c.logger.Error(ctx, "get stream failed", zap.Error(err), zap.String("filename", filename))
		return nil, "", err
	}
	
	// Create a pipe to stream the data
	pr, pw := io.Pipe()
	
	go func() {
		defer pw.Close()
		
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				s, ok := status.FromError(err)
				if ok && s.Code() == codes.Canceled {
					return
				}
				c.logger.Error(ctx, "stream recv failed", zap.Error(err), zap.String("filename", filename))
				return
			}
			
			if _, writeErr := pw.Write(resp.Chunk); writeErr != nil {
				c.logger.Error(ctx, "pipe write failed", zap.Error(writeErr))
				return
			}
		}
	}()
	
	return pr, "", nil
}

// GetFile downloads the entire file at once (for small files)
func (c *FileServerClient) GetFile(ctx context.Context, filename string) ([]byte, error) {
	req := &fileserverpb.GetRequest{
		Filename: filename,
	}
	
	stream, err := c.client.Get(ctx, req)
	if err != nil {
		c.logger.Error(ctx, "get stream failed", zap.Error(err), zap.String("filename", filename))
		return nil, err
	}
	
	var data []byte
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			c.logger.Error(ctx, "stream recv failed", zap.Error(err), zap.String("filename", filename))
			return nil, err
		}
		
		data = append(data, resp.Chunk...)
	}
	
	return data, nil
}