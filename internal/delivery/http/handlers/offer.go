package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	"go.uber.org/zap"
)

func (o *offerHandler) GetOffers(w http.ResponseWriter, r *http.Request) {
	page, err := parseIntQueryParam(r, "page", 1)
	if err != nil {
		o.logger.Error(r.Context(), "invalid or no page", zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	limit, err := parseIntQueryParam(r, "limit", 10)
	if err != nil {
		o.logger.Error(r.Context(), "invalid or no limit", zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	result, err := o.offerUsecase.List(r.Context(),page, limit)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	response.WriteJSON(w, http.StatusOK, result)
}

func (o *offerHandler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	var req CreateOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		o.logger.Error(r.Context(), "invalid JSON", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "ошибка создания предложения")
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		o.logger.Error(r.Context(), "no userID")
		response.HandleError(w, nil, http.StatusUnauthorized, "требуется авторизация")
		return
	}

	userId, err := strconv.Atoi(userID)
	if err != nil {
		o.logger.Error(r.Context(), "invalid user id", zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка пользовательского характера")
	}
	offer := domain.Offer{
		UserID:      userId,
		Title:       req.Title,
		Description: req.Description,
		Image:       req.Image,
		Price:       req.Price,
		Area:        req.Area,
		Rooms:       req.Rooms,
		Address:     req.Address,
		OfferType:   req.OfferType,
	}

	if err := o.offerUsecase.Create(r.Context(), &offer); err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "failed to create offer", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(offer)
}

func (o *offerHandler) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	req, err := parseIntQueryParam(r, "id", 0)
	if err != nil {
		o.logger.Error(r.Context(), "invalid or no id", zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения параметра")
		return
	}
	if err := o.offerUsecase.Delete(r.Context(), strconv.Itoa(req)); err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка удаления предложения")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (o *offerHandler) UpdateOffer(w http.ResponseWriter, r *http.Request) {
	var req UpdateOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		o.logger.Error(r.Context(), "invalid JSON", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "ошибка обработки входных данных")
		return
	}

	offer := domain.Offer{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Image:       req.Image,
		Price:       req.Price,
		Area:        req.Area,
		Rooms:       req.Rooms,
		Address:     req.Address,
		OfferType:   req.OfferType,
	}
	if err := o.offerUsecase.Update(r.Context(), &offer); err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка обновления предложения")
		return
	}
	w.WriteHeader(http.StatusOK)
}