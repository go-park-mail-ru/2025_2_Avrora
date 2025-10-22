package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type IOfferRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Offer, error)
	List(ctx context.Context, page, limit int) ([]*domain.Offer, error)
	Create(ctx context.Context, offer *domain.Offer) error
	Update(ctx context.Context, offer *domain.Offer) error
	Delete(ctx context.Context, id string) error
	CountAll(ctx context.Context, ) (int, error)
	ListByUserID(ctx context.Context, userID string) ([]*domain.Offer, error)
}

type offerUsecase struct {
	offerRepo IOfferRepository
	log *log.Logger
}

func NewOfferUsecase(repo IOfferRepository, log *log.Logger) *offerUsecase {
	return &offerUsecase{offerRepo: repo, log: log}
}