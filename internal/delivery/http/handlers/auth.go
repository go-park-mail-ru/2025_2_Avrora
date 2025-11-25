package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	usecase "github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	"go.uber.org/zap"
)

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(r.Context(), "invalid JSON", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, ErrInvalidJSON.Error())
		return
	}
	
	if err := validateRegisterRequest(&req); err != nil {
		h.logger.Error(r.Context(), "invalid input", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authService.Register(r.Context(), req.Email, req.Password); err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserAlreadyExists):
			response.HandleError(w, err, http.StatusConflict, usecase.ErrUserAlreadyExists.Error())
			return
		case errors.Is(err, usecase.ErrInvalidInput):
			response.HandleError(w, err, http.StatusBadRequest, usecase.ErrInvalidInput.Error())
			return
		default:
			response.HandleError(w, err, http.StatusInternalServerError, usecase.ErrServerSideError.Error())
		}
		return
	}

	response.WriteJSON(w, http.StatusCreated, RegisterResponse{
		Email: req.Email,
	})
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(r.Context(), "invalid JSON", zap.Error(err))
		response.HandleError(w, ErrInvalidJSON, http.StatusBadRequest, ErrInvalidJSON.Error())
	}

	if err := validateLoginRequest(&req); err != nil {
		h.logger.Error(r.Context(), "invalid input", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			response.HandleError(w, err, http.StatusUnauthorized, usecase.ErrInvalidCredentials.Error())
		case errors.Is(err, usecase.ErrInvalidInput):
			response.HandleError(w, err, http.StatusBadRequest, usecase.ErrInvalidInput.Error())
		default:
			response.HandleError(w, err, http.StatusInternalServerError, usecase.ErrServerSideError.Error())
		}
		return
	}

	response.WriteJSON(w, http.StatusOK, AuthResponse{
		Email: req.Email,
		Token: token,
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	expiredToken, err := h.authService.Logout(r.Context())
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка генерации jwt")
	}

	response.WriteJSON(w, http.StatusOK, LogoutResponse{
		Message: "success",
		Token:   expiredToken,
	})
}
