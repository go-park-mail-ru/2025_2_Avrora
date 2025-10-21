package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"go.uber.org/zap"
)

func (uc *offerUsecase) Create(ctx context.Context, offer *domain.Offer) error {
	if offer.Title == "" {
		return ErrInvalidInput
	}

	return uc.offerRepo.Create(ctx, offer)
}

func (uc *offerUsecase) GetByID(ctx context.Context, id string) (*domain.Offer, error) {
	if id == "" {
		uc.log.Error(ctx, "id is empty", zap.Error(ErrInvalidInput))
		return &domain.Offer{}, ErrInvalidInput
	}
	return uc.offerRepo.GetByID(ctx, id)
}

func (uc *offerUsecase) List(ctx context.Context, page, limit int) ([]*domain.Offer, error) {
	if page < 1 {
		uc.log.Error(ctx, "page is empty", zap.Error(ErrInvalidInput))
		return nil, ErrInvalidInput
	}
	if limit < 1 || limit >= 100 {
		uc.log.Error(ctx, "limit is empty", zap.Error(ErrInvalidInput))
		return nil, ErrInvalidInput
	}
	return uc.offerRepo.List(ctx, page, limit)
}

func (uc *offerUsecase) Update(ctx context.Context, offer *domain.Offer) error {
	if offer.Title == "" {
		uc.log.Error(ctx, "title is empty", zap.Error(ErrInvalidInput))
		return ErrInvalidInput
	}
	return uc.offerRepo.Update(ctx, offer)
}

func (uc *offerUsecase) ListByUserID(ctx context.Context, userID string) ([]*domain.Offer, error) {
	return uc.offerRepo.ListByUserID(ctx, userID)
}


func (uc *offerUsecase) CountAll(ctx context.Context, ) (int, error) {
	return uc.offerRepo.CountAll(ctx)
}

func (uc *offerUsecase) Delete(ctx context.Context, id string) error {
	return uc.offerRepo.Delete(ctx, id)
}