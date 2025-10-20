package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	"go.uber.org/zap"
)

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(r.Context(), "invalid JSON", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "invalid JSON")
	}
	
	if err := validateRegisterRequest(&req); err != nil {
		h.logger.Error(r.Context(), "invalid input", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authUsecase.Register(r.Context(), req.Email, req.Password); err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserAlreadyExists):
			h.logger.Error(r.Context(), "user already exists", zap.Error(err))
			response.HandleError(w, err, http.StatusConflict, usecase.ErrUserAlreadyExists.Error())
		case errors.Is(err, usecase.ErrInvalidInput):
			h.logger.Error(r.Context(), "invalid input", zap.Error(err))
			response.HandleError(w, err, http.StatusBadRequest, usecase.ErrInvalidInput.Error())
		default:
			h.logger.Error(r.Context(), "server side error", zap.Error(err))
			response.HandleError(w, err, http.StatusInternalServerError, usecase.ErrServerSideError.Error())
		}
		return
	}


	h.logger.Info(r.Context(), "register success")
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

	token, err := h.authUsecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			h.logger.Error(r.Context(), "invalid credentials", zap.Error(err))
			response.HandleError(w, err, http.StatusUnauthorized, "invalid credentials")
		case errors.Is(err, usecase.ErrInvalidInput):
			h.logger.Error(r.Context(), "invalid input", zap.Error(err))
			response.HandleError(w, err, http.StatusBadRequest, "invalid input")
		default:
			h.logger.Error(r.Context(), "server side error", zap.Error(err))
			response.HandleError(w, err, http.StatusInternalServerError, "server side error")
		}
		return
	}

	h.logger.Info(r.Context(), "login success")
	response.WriteJSON(w, http.StatusOK, AuthResponse{
		Email: req.Email,
		Token: token,
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	expiredToken, err := h.authUsecase.Logout(r.Context())
	if err != nil {
		h.logger.Error(r.Context(), "error generating jwt", zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка генерации jwt")
	}

	h.logger.Info(r.Context(), "logout success")
	response.WriteJSON(w, http.StatusOK, LogoutResponse{
		Message: "success",
		Token:   expiredToken,
	})
}
