package handlers

import "github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"


type IOfferUsecase interface {
	List(page, limit int) ([]*domain.Offer, error)
	GetByID(id string) (*domain.Offer, error)
	Update(offer *domain.Offer) error
	Create(offer *domain.Offer) error
	Delete(id string) error
	ListByUserID(userID string) ([]*domain.Offer, error)
}

type offerHandler struct {
	offerUsecase IOfferUsecase
}

func NewOfferHandler(uc IOfferUsecase) *offerHandler {
	return &offerHandler{offerUsecase: uc}
}