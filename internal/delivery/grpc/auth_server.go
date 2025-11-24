package service

import (
	"context"

	"google.golang.org/grpc"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/handlers"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	auth "github.com/go-park-mail-ru/2025_2_Avrora/proto/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.authUsecase.Register(ctx, req.Email, req.Password); err != nil {
		s.logger.Error(ctx, "failed to register user", zap.Error(err))
		return nil, mapAuthErrorToGRPCStatus(err)
	}

	return &auth.RegisterResponse{
		Email: req.Email,
	}, nil
}

func (s *authServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	s.logger.Info(ctx, "received login request", zap.String("email", req.Email))

	if err := s.validateLoginRequest(req); err != nil {
		s.logger.Error(ctx, "invalid login request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.authUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		s.logger.Error(ctx, "failed to login user", zap.Error(err))
		return nil, mapAuthErrorToGRPCStatus(err)
	}

	return &auth.LoginResponse{
		Token: token,
		Email: req.Email,
	}, nil
}

func (s *authServer) Logout(ctx context.Context, _ *emptypb.Empty) (*auth.LogoutResponse, error) {
	expiredToken, err := s.authUsecase.Logout(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to process logout")
	}

	return &auth.LogoutResponse{
		Message: "success",
		Token:   expiredToken,
	}, nil
}

func (s *authServer) validateRegisterRequest(req *auth.RegisterRequest) error {
	if req.Email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if len(req.Password) < 8 {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}
	return nil
}

func (s *authServer) validateLoginRequest(req *auth.LoginRequest) error {
	if req.Email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func mapAuthErrorToGRPCStatus(err error) error {
	switch {
	case err == nil:
		return nil
	case usecase.ErrUserAlreadyExists.Error() == err.Error():
		return status.Error(codes.AlreadyExists, "user already exists")
	case usecase.ErrInvalidCredentials.Error() == err.Error():
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case usecase.ErrInvalidInput.Error() == err.Error():
		return status.Error(codes.InvalidArgument, "invalid input")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

// RegisterAuthServer registers the auth service with the gRPC server
func RegisterAuthServer(grpcServer *grpc.Server, authUsecase handlers.IAuthService, logger *log.Logger) {
	auth.RegisterAuthServiceServer(grpcServer, NewAuthServer(authUsecase, logger))
}