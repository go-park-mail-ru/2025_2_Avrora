package handlers

import (
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type authHandler struct {
	authService IAuthService
	logger *log.Logger
}

func NewAuthHandler(uc IAuthService, logger *log.Logger) *authHandler {
	return &authHandler{authService: uc, logger: logger}
}