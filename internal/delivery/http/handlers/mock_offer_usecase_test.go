package handlers

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/stretchr/testify/mock"
)

// mockOfferUsecase — ручной мок для интерфейса IOfferUsecase
type mockOfferUsecase struct {
	mock.Mock
}

func (m *mockOfferUsecase) ListOffersInFeed(ctx context.Context, page, limit int) (*domain.OffersInFeed, error) {
	args := m.Called(ctx, page, limit)
	if res, ok := args.Get(0).(*domain.OffersInFeed); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockOfferUsecase) Get(ctx context.Context, id string) (*domain.Offer, error) {
	args := m.Called(ctx, id)
	if res, ok := args.Get(0).(*domain.Offer); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockOfferUsecase) Update(ctx context.Context, offer *domain.Offer) error {
	args := m.Called(ctx, offer)
	return args.Error(0)
}

func (m *mockOfferUsecase) Create(ctx context.Context, offer *domain.Offer) error {
	args := m.Called(ctx, offer)
	return args.Error(0)
}

func (m *mockOfferUsecase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockOfferUsecase) ListOffersInFeedByUserID(ctx context.Context, userID string, page, limit int) (*domain.OffersInFeed, error) {
	args := m.Called(ctx, userID, page, limit)
	if res, ok := args.Get(0).(*domain.OffersInFeed); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockOfferUsecase) FilterOffers(ctx context.Context, f *domain.OfferFilter, limit, offset int) ([]domain.OfferInFeed, error) {
	args := m.Called(ctx, f, limit, offset)
	if res, ok := args.Get(0).([]domain.OfferInFeed); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

// Проверка, что мок реализует интерфейс
var _ IOfferUsecase = (*mockOfferUsecase)(nil)
