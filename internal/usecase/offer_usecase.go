package usecase

import (
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
)

type IOfferRepository interface {
	GetByID(id string) (*domain.Offer, error)
	List(page, limit int) ([]*domain.Offer, error)
	Create(offer *domain.Offer) error
	Update(offer *domain.Offer) error
	Delete(id string) error
	CountAll() (int, error)
	ListByUserID(userID string) ([]*domain.Offer, error)
}

type offerUsecase struct {
	offerRepo IOfferRepository
}

func NewOfferUsecase(repo IOfferRepository) *offerUsecase {
	return &offerUsecase{offerRepo: repo}
}