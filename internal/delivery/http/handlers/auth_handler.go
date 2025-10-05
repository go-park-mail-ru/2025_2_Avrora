package handlers

import (
	"net/http"
	
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authUsecase usecase.AuthUsecase
}

var _ AuthHandler = (*authHandler)(nil)

func NewAuthHandler(uc usecase.AuthUsecase) *authHandler {
	return &authHandler{authUsecase: uc}
}
