package service

import (
	"context"

	"google.golang.org/grpc"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	auth "github.com/go-park-mail-ru/2025_2_Avrora/proto/auth"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authServer struct {
	auth.UnimplementedAuthServiceServer
	authUsecase handlers.IAuthService
	logger      *log.Logger
}

func NewAuthServer(authUsecase handlers.IAuthService, logger *log.Logger) *authServer {
	return &authServer{
		authUsecase: authUsecase,
		logger:      logger.With(zap.String("service", "auth_grpc")),
	}
}

func (s *authServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	s.logger.Info(ctx, "received register request", zap.String("email", req.Email))

	if err := s.validateRegisterRequest(req); err != nil {
		s.logger.Error(ctx, "invalid register request", zap.Error(err))
		return nil, err
	}

	if err := s.authUsecase.Register(ctx, req.Email, req.Password); err != nil {
		s.logger.Error(ctx, "failed to register user", zap.Error(err))
		return nil, err
	}

	return &auth.RegisterResponse{
		Email: req.Email,
	}, nil
}

func (s *authServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	s.logger.Info(ctx, "received login request", zap.String("email", req.Email))

	if err := s.validateLoginRequest(req); err != nil {
		s.logger.Error(ctx, "invalid login request", zap.Error(err))
		return nil, err
	}

	token, err := s.authUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		s.logger.Error(ctx, "failed to login user", zap.Error(err))
		return nil, err
	}

	return &auth.LoginResponse{
		Token: token,
		Email: req.Email,
	}, nil
}

func (s *authServer) Logout(ctx context.Context, _ *emptypb.Empty) (*auth.LogoutResponse, error) {
	expiredToken, err := s.authUsecase.Logout(ctx)
	if err != nil {
		return nil, err
	}

	return &auth.LogoutResponse{
		Message: "success",
		Token:   expiredToken,
	}, nil
}

func (s *authServer) validateRegisterRequest(req *auth.RegisterRequest) error {
	if req.Email == "" {
		return usecase.ErrInvalidInput
	}
	if req.Password == "" {
		return usecase.ErrInvalidInput
	}
	if len(req.Password) < 8 {
		return usecase.ErrInvalidInput
	}
	return nil
}

func (s *authServer) validateLoginRequest(req *auth.LoginRequest) error {
	if req.Email == "" {
		return usecase.ErrInvalidInput
	}
	if req.Password == "" {
		return usecase.ErrInvalidInput
	}
	return nil
}

// RegisterAuthServer registers the auth service with the gRPC server
func RegisterAuthServer(grpcServer *grpc.Server, authUsecase handlers.IAuthService, logger *log.Logger) {
	auth.RegisterAuthServiceServer(grpcServer, NewAuthServer(authUsecase, logger))
}