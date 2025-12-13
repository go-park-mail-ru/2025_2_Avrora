package usecase

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type IOfferRepository interface {
	List(ctx context.Context, page, limit int) (*domain.OffersInFeed, error)
	Create(ctx context.Context, offer *domain.Offer) error
	Update(ctx context.Context, offer *domain.Offer) error
	Delete(ctx context.Context, id string) error
	CountAll(ctx context.Context) (int, error)
	ListByUserID(ctx context.Context, userID string, page, limit int) (*domain.OffersInFeed, error)
	GetByID(ctx context.Context, id string) (*domain.Offer, error)
	FilterOffers(ctx context.Context, f *domain.OfferFilter, limit, offset int) ([]domain.OfferInFeed, error)
	GetOfferPriceHistory(ctx context.Context, id string) ([]domain.PricePoint, error)
	LogView(ctx context.Context, offerID string) error
	ToggleLike(ctx context.Context, offerID, userID string) error
	GetOfferViewCount(ctx context.Context, id string) (int, error)
	GetOfferLikeCount(ctx context.Context, id string) (int, error)
	IsOfferLiked(ctx context.Context, offerID, userID string) (bool, error)
	ListPaidOffers(ctx context.Context, page, limit int) (*domain.OffersInFeed, error)
	InsertPaidAdvertisement(ctx context.Context, offerID string, expiresAt time.Time) error
}

type offerUsecase struct {
	offerRepo IOfferRepository
	log       *log.Logger
}

func NewOfferUsecase(repo IOfferRepository, log *log.Logger) *offerUsecase {
	return &offerUsecase{offerRepo: repo, log: log}
}

func (uc *offerUsecase) FilterOffers(ctx context.Context, f *domain.OfferFilter, limit, offset int) ([]domain.OfferInFeed, error) {
	if f == nil {
		return nil, domain.ErrInvalidInput
	}
	offers, err := uc.offerRepo.FilterOffers(ctx, f, limit, offset)
	if err != nil {
		uc.log.Error(ctx, "failed to filter offers", zap.Error(err))
		return nil, err
	}
	return offers, nil
}

func (uc *offerUsecase) GetOfferPriceHistory(ctx context.Context, id string) ([]domain.PricePoint, error) {
	return uc.offerRepo.GetOfferPriceHistory(ctx, id)
}