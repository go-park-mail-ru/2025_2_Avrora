package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

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
	result, err := o.offerUsecase.ListOffersInFeed(r.Context(), page, limit)
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

	// где то тут надо из адреса сделать location ???
	// то же самое с комлексом
	offer := &domain.Offer{
		Title:        req.Title,
		Description:  req.Description,
		ImageURLs:    req.ImageURLs,
		Price:        int64(req.Price),
		Area:         req.Area,
		Rooms:        req.Rooms,
		Address:      req.Address,
		OfferType:    domain.OfferType(req.OfferType),
		PropertyType: domain.PropertyType(req.PropertyType),
		Floor:        &req.Floor,
		TotalFloors:  &req.TotalFloors,
		UserID:       userID,
	}

	if err := o.offerUsecase.Create(r.Context(), offer); err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "failed to create offer", http.StatusInternalServerError)
		}
		return
	}

	response.WriteJSON(w, http.StatusCreated, offer)
}

func (o *offerHandler) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(string)
	if !ok {
		o.logger.Error(r.Context(), "invalid or no id")
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	if err := o.offerUsecase.Delete(r.Context(), id); err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка удаления предложения")
		return
	}
	response.WriteJSON(w, http.StatusOK, nil)
}

func (o *offerHandler) UpdateOffer(w http.ResponseWriter, r *http.Request) {
	var req UpdateOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		o.logger.Error(r.Context(), "invalid JSON", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "ошибка обработки входных данных")
		return
	}

	offer := domain.Offer{
		Title:        req.Title,
		Description:  req.Description,
		ImageURLs:    req.ImageURLs,
		Price:        int64(req.Price),
		Area:         req.Area,
		Rooms:        req.Rooms,
		Address:      req.Address,
		OfferType:    domain.OfferType(req.OfferType),
		PropertyType: domain.PropertyType(req.PropertyType),
		Floor:        &req.Floor,
		TotalFloors:  &req.TotalFloors,
		Deposit:      &req.Deposit,
		Commission:   &req.Commission,
		RentalPeriod: &req.RentalPeriod,
		Status:       domain.OfferStatus(req.Status),
		LivingArea:   &req.LivingArea,
		KitchenArea:  &req.KitchenArea,
	}
	if err := o.offerUsecase.Update(r.Context(), &offer); err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка обновления предложения")
		return
	}
	response.WriteJSON(w, http.StatusOK, offer)
}

func (o *offerHandler) GetOffer(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(string)
	if !ok {
		o.logger.Error(r.Context(), "invalid or no id")
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	offer, err := o.offerUsecase.Get(r.Context(), id)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	response.WriteJSON(w, http.StatusOK, offer)
}