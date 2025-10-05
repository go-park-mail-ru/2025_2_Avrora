package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
)

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req usecase.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.HandleError(w, err, http.StatusBadRequest, "invalid JSON")
	}

	if err := h.authUsecase.Register(req.Email, req.Password); err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserAlreadyExists):
			response.HandleError(w, err, http.StatusConflict, "user already exists")
		case errors.Is(err, usecase.ErrInvalidInput):
			response.HandleError(w, err, http.StatusBadRequest, "invalid input")
		default:
			response.HandleError(w, err, http.StatusInternalServerError, "server side error")
		}
		return
	}

	
	response.WriteJSON(w, http.StatusCreated, usecase.RegisterResponse{
		Email: req.Email,
	})
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req usecase.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.HandleError(w, err, http.StatusBadRequest, "invalid JSON")
	}

	token, err := h.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			response.HandleError(w, err, http.StatusUnauthorized, "invalid credentials")
		case errors.Is(err, usecase.ErrInvalidInput):
			response.HandleError(w, err, http.StatusBadRequest, "invalid input")
		default:
			response.HandleError(w, err, http.StatusInternalServerError, "server side error")
		}
		return
	}

	response.WriteJSON(w, http.StatusOK, usecase.AuthResponse{
		Email: req.Email,
		Token: token,
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	expiredToken, err := h.authUsecase.Logout()
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка генерации jwt")
	}

	response.WriteJSON(w, http.StatusOK, usecase.LogoutResponse{
		Message: "success",
		Token:   expiredToken,
	})
}
