package usecase

import (
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/infrastructure/db"
)

type OfferUsecase interface {
	Create(offer domain.Offer) error
	Update(offer domain.Offer) error
	Delete(id string) error
	GetByID(id string) (domain.Offer, error)
	List(page, limit int) ([]domain.Offer, error)
	ListByUserID(userID string) ([]domain.Offer, error)
	CountAll() (int, error)
}

var _ OfferUsecase = (*offerUsecase)(nil)

type offerUsecase struct {
	offerRepo *db.OfferRepository
}

func NewOfferUsecase(repo *db.OfferRepository) *offerUsecase {
	return &offerUsecase{offerRepo: repo}
}