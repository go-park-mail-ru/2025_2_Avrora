package handlers

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type IAuthUsecase interface {
	Register(ctx context.Context, email string, password string) error
	Login(ctx context.Context, email string, password string) (string, error)
	Logout(ctx context.Context) (string, error)
}

type authHandler struct {
	authUsecase IAuthUsecase
	logger *log.Logger
}

func NewAuthHandler(uc IAuthUsecase, logger *log.Logger) *authHandler {
	return &authHandler{authUsecase: uc, logger: logger}
}
