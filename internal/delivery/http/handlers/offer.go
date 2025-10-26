package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/utils"
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

	// где то тут надо из адреса сделать location ???
	// то же самое с комлексом
	offer := &domain.Offer{
		Title:        req.Title,
		Description:  req.Description,
		ImageURLs:    req.ImageURLs,
		LocationID:   utils.AddressToLocation(req.Address).ID,
		Price:        int64(req.Price),
		Area:         req.Area,
		Rooms:        req.Rooms,
		Address:      req.Address,
		OfferType:    domain.OfferType(req.OfferType),
		PropertyType: domain.PropertyType(req.PropertyType),
		Floor:        &req.Floor,
		TotalFloors:  &req.TotalFloors,
		UserID:       req.UserID,
	}

	if err := o.offerUsecase.Create(r.Context(), offer); err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			response.HandleError(w, err, http.StatusBadRequest, "невалидные данные")
		} else {
			response.HandleError(w, err, http.StatusInternalServerError, "ошибка создания предложения")
		}
		return
	}

	response.WriteJSON(w, http.StatusCreated, offer)
}

func (o *offerHandler) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	id := GetPathParameter(r, "/api/v1/offers/delete/")
	if id == "" {
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

	id := GetPathParameter(r, "/api/v1/offers/update/")
	if id == "" {
		o.logger.Error(r.Context(), "invalid or no id")
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}

	// Надо достать location_id как то через адрес
	// YandexMapsURL
	// Широта долгота
	// Нормализованный адрес

	location := utils.AddressToLocation(req.Address) // Пока так

	offer := domain.Offer{
		ID:           id,
		UserID:       req.UserID,
		LocationID:   location.ID,
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
	id := GetPathParameter(r, "/api/v1/offers")
	if id == "" {
		o.logger.Error(r.Context(), "invalid or no id")
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}

	println(id)

	offer, err := o.offerUsecase.Get(r.Context(), id)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	response.WriteJSON(w, http.StatusOK, offer)
}

func (o *offerHandler) GetMyOffers(w http.ResponseWriter, r *http.Request) {
	userID := GetPathParameter(r, "/api/v1/profile/myoffers/")
	if userID == "" {
		o.logger.Error(r.Context(), "invalid or no userID")
		response.HandleError(w, nil, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	offers, err := o.offerUsecase.ListOffersInFeedByUserID(r.Context(), userID, 1, 10)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	response.WriteJSON(w, http.StatusOK, offers)
}
