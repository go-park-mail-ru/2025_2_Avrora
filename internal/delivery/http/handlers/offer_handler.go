package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
)

type OfferHandler interface {
	GetOffersHandler(w http.ResponseWriter, r *http.Request)
	CreateOfferHandler(w http.ResponseWriter, r *http.Request)
	UpdateOfferHandler(w http.ResponseWriter, r *http.Request)
	DeleteOfferHandler(w http.ResponseWriter, r *http.Request)
}

type offerHandler struct {
	offerUsecase usecase.OfferUsecase
}

var _ OfferHandler = (*offerHandler)(nil)

func NewOfferHandler(uc usecase.OfferUsecase) *offerHandler {
	return &offerHandler{offerUsecase: uc}
}