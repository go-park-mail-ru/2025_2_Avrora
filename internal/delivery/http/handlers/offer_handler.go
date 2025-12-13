package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type IOfferUsecase interface {
	ListOffersInFeed(ctx context.Context, page, limit int) (*domain.OffersInFeed, error)
	Get(ctx context.Context, id string) (*domain.Offer, error)
	Update(ctx context.Context, offer *domain.Offer) error
	Create(ctx context.Context, offer *domain.Offer) error
	Delete(ctx context.Context, id string) error
	ListOffersInFeedByUserID(ctx context.Context, userID string, page, limit int) (*domain.OffersInFeed, error)
	FilterOffers(ctx context.Context, f *domain.OfferFilter, limit, offset int) ([]domain.OfferInFeed, error)
	GetOfferPriceHistory(ctx context.Context, id string) ([]domain.PricePoint, error)
	ViewOffer(ctx context.Context, offerID string) error
	ToggleOfferLike(ctx context.Context, offerID, userID string) error
	GetOfferViewCount(ctx context.Context, offerID string) (int, error)
	GetOfferLikeCount(ctx context.Context, offerID string) (int, error)
	IsOfferLiked(ctx context.Context, offerID, userID string) (bool, error)
	ListPaidOffers(ctx context.Context, page, limit int) (*domain.OffersInFeed, error)
	InsertPaidAdvertisement(ctx context.Context, offerID string, expiresAt time.Time) error
	ListLikedOffersByUserID(ctx context.Context, userID string, page, limit int) (*domain.OffersInFeed, error)
}

type offerHandler struct {
	offerUsecase IOfferUsecase
	logger       *log.Logger
}

func NewOfferHandler(uc IOfferUsecase, logger *log.Logger) *offerHandler {
	return &offerHandler{offerUsecase: uc, logger: logger}
}
func (h *offerHandler) FilterOffers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	f := &domain.OfferFilter{}

	if v := q.Get("offer_type"); v != "" {
		f.OfferType = &v
	}
	if v := q.Get("property_type"); v != "" {
		f.PropertyType = &v
	}
	if v := q.Get("status"); v != "" {
		f.Status = &v
	}
	if v := q.Get("rooms"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			f.Rooms = &i
		}
	}
	if v := q.Get("price_min"); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.PriceMin = &i
		}
	}
	if v := q.Get("price_max"); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.PriceMax = &i
		}
	}
	if v := q.Get("area_min"); v != "" {
		if f64, err := strconv.ParseFloat(v, 64); err == nil {
			f.AreaMin = &f64
		}
	}
	if v := q.Get("area_max"); v != "" {
		if f64, err := strconv.ParseFloat(v, 64); err == nil {
			f.AreaMax = &f64
		}
	}
	if v := q.Get("address"); v != "" {
		f.Address = &v
	}

	limit := 20
	offset := 0
	if v := q.Get("limit"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			limit = i
		}
	}
	if v := q.Get("offset"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			offset = i
		}
	}

	offers, err := h.offerUsecase.FilterOffers(ctx, f, limit, offset)
	if err != nil {
		h.logger.Error(ctx, "failed to filter offers", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(offers); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// GetLikedOffers returns offers liked by the current authenticated user ("Избранное")
func (h *offerHandler) GetLikedOffers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем userID из контекста (мидлвара аутентификации уже положил его туда)
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		h.logger.Warn(ctx, "unauthenticated request to /liked")
		response.HandleError(w, nil, http.StatusUnauthorized, "требуется аутентификация")
		return
	}

	// Пагинация
	page, err := parseIntQueryParam(r, "page", 1)
	if err != nil {
		page = 1
	}
	limit, err := parseIntQueryParam(r, "limit", 10)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Вызываем usecase
	offers, err := h.offerUsecase.ListLikedOffersByUserID(ctx, userID, page, limit)
	if err != nil {
		h.logger.Error(ctx, "failed to get liked offers",
			zap.String("user_id", userID), zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения избранного")
		return
	}

	// ⚠️ Опционально: проставить IsLiked = true для всех — логично для "избранного"
	for i := range offers.Offers {
		offers.Offers[i].IsLiked = true
	}

	response.WriteJSON(w, http.StatusOK, offers)
}
