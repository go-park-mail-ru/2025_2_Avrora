package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
)

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.HandleError(w, err, http.StatusBadRequest, "invalid JSON")
	}
	
	if err := validateRegisterRequest(&req); err != nil {
		response.HandleError(w, err, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authUsecase.Register(req.Email, req.Password); err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserAlreadyExists):
			response.HandleError(w, err, http.StatusConflict, usecase.ErrUserAlreadyExists.Error())
		case errors.Is(err, usecase.ErrInvalidInput):
			response.HandleError(w, err, http.StatusBadRequest, usecase.ErrInvalidInput.Error())
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
		response.HandleError(w, ErrInvalidJSON, http.StatusBadRequest, ErrInvalidJSON.Error())
	}

	if err := validateLoginRequest(&req); err != nil {
		response.HandleError(w, err, http.StatusBadRequest, err.Error())
		return
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

	response.WriteJSON(w, http.StatusOK, AuthResponse{
		Email: req.Email,
		Token: token,
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	expiredToken, err := h.authUsecase.Logout()
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка генерации jwt")
	}

	response.WriteJSON(w, http.StatusOK, LogoutResponse{
		Message: "success",
		Token:   expiredToken,
	})
}
