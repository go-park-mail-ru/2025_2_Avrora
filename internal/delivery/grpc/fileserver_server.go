package service

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	fileserverpb "github.com/go-park-mail-ru/2025_2_Avrora/proto/fileserver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileServer struct {
	fileserverpb.UnimplementedFileServerServer
	storageDir string
	baseURL    string
	logger     *log.Logger
}

func NewFileServer(storageDir, baseURL string, logger *log.Logger) *FileServer {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		logger.Panic("failed to create storage dir", zap.Error(err))
	}
	return &FileServer{
		storageDir: storageDir,
		baseURL:    baseURL,
		logger:     logger.With(zap.String("service", "fileserver")),
	}
}

func (s *FileServer) Upload(ctx context.Context, req *fileserverpb.UploadRequest) (*fileserverpb.UploadResponse, error) {
	s.logger.Info(ctx, "uploading file", zap.String("filename", req.Filename))

	if len(req.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty data")
	}
	if req.Filename == "" {
		return nil, status.Error(codes.InvalidArgument, "filename required")
	}

	// Sanitize: only basename (no path traversal)
	filename := filepath.Base(req.Filename)
	fullPath := filepath.Join(s.storageDir, filename)

	if err := os.WriteFile(fullPath, req.Data, 0644); err != nil {
		s.logger.Error(ctx, "failed to save file", zap.Error(err))
		return nil, status.Error(codes.Internal, "storage error")
	}

	// Construct URL using baseURL and filename
	url := s.baseURL + "/" + filename
	return &fileserverpb.UploadResponse{Url: url}, nil
}

func (s *FileServer) Get(req *fileserverpb.GetRequest, stream fileserverpb.FileServer_GetServer) error {
	ctx := stream.Context()
	
	if req.Filename == "" {
		return status.Error(codes.InvalidArgument, "filename required")
	}

	filename := filepath.Base(req.Filename)
	fullPath := filepath.Join(s.storageDir, filename)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			s.logger.Error(ctx, "file not found", zap.String("filename", filename))
			return status.Error(codes.NotFound, "file not found")
		}
		s.logger.Error(ctx, "failed to open file", zap.Error(err), zap.String("filename", filename))
		return status.Error(codes.Internal, "failed to open file")
	}
	defer file.Close()

	// Get file info for better error handling
	stat, err := file.Stat()
	if err != nil {
		s.logger.Error(ctx, "failed to get file stat", zap.Error(err))
		return status.Error(codes.Internal, "failed to get file info")
	}

	if stat.IsDir() {
		return status.Error(codes.InvalidArgument, "cannot serve directory")
	}

	// Stream file in chunks (32KB per chunk)
	buf := make([]byte, 32*1024) // 32KB chunks
	for {
		n, err := file.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			
			if err := stream.Send(&fileserverpb.GetResponse{Chunk: chunk}); err != nil {
				// Client disconnected or stream error
				if status.Code(err) == codes.Canceled {
					s.logger.Info(ctx, "client disconnected during file stream", zap.String("filename", filename))
					return nil
				}
				s.logger.Error(ctx, "failed to send chunk", zap.Error(err))
				return err
			}
		}
		
		if err == io.EOF {
			// Successfully reached end of file
			return nil
		}
		
		if err != nil {
			s.logger.Error(ctx, "error reading file", zap.Error(err))
			return status.Error(codes.Internal, "error reading file")
		}
	}
}

func RegisterFileServerServer(s *grpc.Server, logger *log.Logger) {
	fileserverpb.RegisterFileServerServer(s, NewFileServer("./image", "/api/v1/image", logger))
}