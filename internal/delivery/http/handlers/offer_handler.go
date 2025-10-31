package handlers

import (
	"context"

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
}

type offerHandler struct {
	offerUsecase IOfferUsecase
	logger *log.Logger
}

func NewOfferHandler(uc IOfferUsecase, logger *log.Logger) *offerHandler {
	return &offerHandler{offerUsecase: uc, logger: logger}
}