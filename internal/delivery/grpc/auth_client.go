package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/go-park-mail-ru/2025_2_Avrora/proto/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authClient struct {
	client auth.AuthServiceClient
	logger *log.Logger
}

func NewAuthClient(addr string, logger *log.Logger) (*authClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	
	return &authClient{
		client: auth.NewAuthServiceClient(conn),
		logger: logger.With(zap.String("service", "auth_client")),
	}, nil
}

func (c *authClient) Register(ctx context.Context, email, password string) error {
	req := &auth.RegisterRequest{
		Email:    email,
		Password: password,
	}
	
	_, err := c.client.Register(ctx, req)
	return err
}

func (c *authClient) Login(ctx context.Context, email, password string) (string, error) {
	req := &auth.LoginRequest{
		Email:    email,
		Password: password,
	}
	
	resp, err := c.client.Login(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *authClient) Logout(ctx context.Context) (string, error) {
	resp, err := c.client.Logout(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}