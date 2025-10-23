package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (p *profileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id := GetPathParameter(r, "/api/v1/profile/"); if id == "" {
		p.log.Error(r.Context(), "invalid or no id")
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка получения параметра")
		return
	}

	result, err := p.profileUsecase.GetProfileByID(r.Context(), id)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения профиля")
		return
	}
	response.WriteJSON(w, http.StatusOK, result)
}

func (p *profileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req ProfileUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p.log.Error(r.Context(), "invalid JSON", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "invalid JSON")
		return
	}

	id := GetPathParameter(r, "/api/v1/profile/update/"); if id == "" {
		p.log.Error(r.Context(), "invalid or no id")
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка получения параметра")
		return
	}

	if err := p.profileUsecase.UpdateProfile(r.Context(), id, &domain.ProfileUpdate{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      domain.UserRole(req.Role),
		AvatarURL: req.AvatarURL,
	}); err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка обновления профиля")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (p *profileHandler) UpdateProfileSecurityByID(w http.ResponseWriter, r *http.Request) {
	var req SecurityUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p.log.Error(r.Context(), "invalid JSON", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "invalid JSON")
		return
	}

	id := GetPathParameter(r, "/api/v1/profile/security/"); if id == "" {
		p.log.Error(r.Context(), "invalid or no id")
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка получения параметра")
		return
	}

	if err := p.profileUsecase.UpdateProfileSecurityByID(r.Context(), id, req.OldPassword, req.NewPassword); err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка обновления профиля")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (p *profileHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	var req UpdateEmail
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p.log.Error(r.Context(), "invalid JSON", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "invalid JSON")
		return
	}

	vars := mux.Vars(r)
	id := vars["user_id"]

	if err := p.profileUsecase.UpdateEmail(r.Context(), id, req.Email); err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка обновления профиля")
		return
	}
	w.WriteHeader(http.StatusOK)
}
