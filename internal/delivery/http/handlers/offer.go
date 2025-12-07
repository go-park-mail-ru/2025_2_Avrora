package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware"
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
		response.HandleError(w, err, http.StatusBadRequest, "нет параметра page")
		return
	}
	q := r.URL.Query()

	// Если в запросе есть хоть один фильтрующий параметр то мы вызываем FilterOffers
	filterKeys := []string{
		"offer_type", "property_type", "status",
		"price_min", "price_max",
		"area_min", "area_max",
		"address",
	}

	hasFilter := false
	for _, key := range filterKeys {
		if q.Get(key) != "" {
			hasFilter = true
			break
		}
	}

	if hasFilter {
		o.FilterOffers(w, r)
		return
	}
	limit, err := parseIntQueryParam(r, "limit", 10)
	if err != nil {
		o.logger.Error(r.Context(), "invalid or no limit", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "нет параметра limit")
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
		Title:            req.Title,
		Description:      req.Description,
		ImageURLs:        req.ImageURLs,
		HousingComplexID: &utils.AddressToComplex(req.Address).ID,
		LocationID:       utils.AddressToLocation(req.Address).ID,
		Price:            int64(req.Price),
		Area:             req.Area,
		Rooms:            req.Rooms,
		Address:          req.Address,
		OfferType:        domain.OfferType(req.OfferType),
		PropertyType:     domain.PropertyType(req.PropertyType),
		Floor:            &req.Floor,
		TotalFloors:      &req.TotalFloors,
		KitchenArea:      &req.KitchenArea,
		LivingArea:       &req.LivingArea,
		Deposit:          &req.Deposit,
		Commission:       &req.Commission,
		RentalPeriod:     &req.RentalPeriod,
		UserID:           req.UserID,
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
		response.HandleError(w, nil, http.StatusBadRequest, "нет id")
		return
	}
	if err := o.offerUsecase.Delete(r.Context(), id); err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка удаления предложения")
		return
	}
	response.WriteJSON(w, http.StatusOK, "success")
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
		response.HandleError(w, nil, http.StatusBadRequest, "нет id")
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
		response.HandleError(w, nil, http.StatusBadRequest, "нет id")
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
	userID := GetPathParameter(r, "/api/v1/profile/myoffers")
	if userID == "" {
		o.logger.Error(r.Context(), "invalid or no userID")
		response.HandleError(w, nil, http.StatusBadRequest, "нет userID")
		return
	}
	offers, err := o.offerUsecase.ListOffersInFeedByUserID(r.Context(), userID, 1, 10)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	response.WriteJSON(w, http.StatusOK, offers)
}

func (o *offerHandler) GetOfferPriceHistory(w http.ResponseWriter, r *http.Request) {
	id := GetPathParameter(r, "/api/v1/offers/pricehistory/")
	if id == "" {
		o.logger.Error(r.Context(), "invalid or no id")
		response.HandleError(w, nil, http.StatusBadRequest, "нет id")
		return
	}
	points, err := o.offerUsecase.GetOfferPriceHistory(r.Context(), id)
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения предложений")
		return
	}
	response.WriteJSON(w, http.StatusOK, points)
}

// ViewOffer logs a view event for an offer (can be anonymous)
func (o *offerHandler) ViewOffer(w http.ResponseWriter, r *http.Request) {
	offerID := GetPathParameter(r, "/api/v1/offers/view/")
	if offerID == "" {
		o.logger.Error(r.Context(), "invalid or no offerID for view")
		response.HandleError(w, nil, http.StatusBadRequest, "нет offerID")
		return
	}
	
	if err := o.offerUsecase.ViewOffer(r.Context(), offerID); err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			response.HandleError(w, err, http.StatusBadRequest, "неверные данные запроса")
			return
		}
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка регистрации просмотра")
		return
	}
	
	response.WriteJSON(w, http.StatusNoContent, nil)
}

// ToggleLike toggles like status for an offer
func (o *offerHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	offerID := GetPathParameter(r, "/api/v1/offers/like/")
	if offerID == "" {
		o.logger.Error(r.Context(), "invalid or no offerID for like")
		response.HandleError(w, nil, http.StatusBadRequest, "нет offerID")
		return
	}
	
	// Must have authenticated user
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		o.logger.Warn(r.Context(), "unauthenticated like attempt")
		response.HandleError(w, nil, http.StatusUnauthorized, "требуется аутентификация")
		return
	}
	
	if err := o.offerUsecase.ToggleOfferLike(r.Context(), offerID, userID); err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			response.HandleError(w, err, http.StatusBadRequest, "неверные данные запроса")
			return
		}
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка обработки лайка")
		return
	}
	
	response.WriteJSON(w, http.StatusNoContent, nil)
}

// GetViewCount returns total views for an offer
func (o *offerHandler) GetViewCount(w http.ResponseWriter, r *http.Request) {
	offerID := GetPathParameter(r, "/api/v1/offers/viewcount/")
	if offerID == "" {
		o.logger.Error(r.Context(), "invalid or no offerID for view count")
		response.HandleError(w, nil, http.StatusBadRequest, "нет offerID")
		return
	}
	
	count, err := o.offerUsecase.GetOfferViewCount(r.Context(), offerID)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			response.HandleError(w, err, http.StatusBadRequest, "неверные данные запроса")
			return
		}
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения просмотров")
		return
	}
	
	response.WriteJSON(w, http.StatusOK, map[string]int{"count": count})
}

// GetLikeCount returns total likes for an offer
func (o *offerHandler) GetLikeCount(w http.ResponseWriter, r *http.Request) {
	offerID := GetPathParameter(r, "/api/v1/offers/likecount/")
	if offerID == "" {
		o.logger.Error(r.Context(), "invalid or no offerID for like count")
		response.HandleError(w, nil, http.StatusBadRequest, "нет offerID")
		return
	}
	
	count, err := o.offerUsecase.GetOfferLikeCount(r.Context(), offerID)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			response.HandleError(w, err, http.StatusBadRequest, "неверные данные запроса")
			return
		}
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения лайков")
		return
	}
	
	response.WriteJSON(w, http.StatusOK, map[string]int{"count": count})
}

// IsOfferLiked checks if current user liked an offer
func (o *offerHandler) IsOfferLiked(w http.ResponseWriter, r *http.Request) {
	offerID := GetPathParameter(r, "/api/v1/offers/isliked/")
	if offerID == "" {
		o.logger.Error(r.Context(), "invalid or no offerID for like check")
		response.HandleError(w, nil, http.StatusBadRequest, "нет offerID")
		return
	}
	
	// Must have authenticated user
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		o.logger.Warn(r.Context(), "unauthenticated like check")
		response.HandleError(w, nil, http.StatusUnauthorized, "требуется аутентификация")
		return
	}
	
	isLiked, err := o.offerUsecase.IsOfferLiked(r.Context(), offerID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			response.HandleError(w, err, http.StatusBadRequest, "неверные данные запроса")
			return
		}
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка проверки лайка")
		return
	}
	
	response.WriteJSON(w, http.StatusOK, map[string]bool{"is_liked": isLiked})
}

func (o *offerHandler) WebHook(w http.ResponseWriter, r *http.Request) {
	o.logger.Info(r.Context(), "webhook request received")
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var body WebhookRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Ошибка при чтении данных", http.StatusBadRequest)
		return
	}

	offer_id := body.Object.Metadata["offer_id"]
	_, err = o.offerUsecase.Get(r.Context(), offer_id)
	if err != nil {
		o.logger.Error(r.Context(), "ошибка получения предложения")
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения предложения")
		return
	}

	err = o.offerUsecase.InsertPaidAdvertisement(r.Context(), offer_id, time.Now().Add(7 * time.Hour * 24))
	if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка регистрации оплаты")
		return
	}

	response.WriteJSON(w, http.StatusOK, "успешно регистрация оплаты")
}

func (o *offerHandler) GetPaidOffers(w http.ResponseWriter, r *http.Request) {
	offers, err := o.offerUsecase.ListPaidOffers(r.Context(), 1, 10); if err != nil {
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения оплаченных предложений")
		return
	}

	response.WriteJSON(w, http.StatusOK, offers)
}