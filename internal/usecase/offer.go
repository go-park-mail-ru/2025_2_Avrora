package usecase

import "github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"

func (uc *offerUsecase) Create(offer *domain.Offer) error {
	if offer.Title == "" {
		return ErrInvalidInput
	}

	return uc.offerRepo.Create(offer)
}

func (uc *offerUsecase) GetByID(id string) (*domain.Offer, error) {
	if id == "" {
		return &domain.Offer{}, ErrInvalidInput
	}
	return uc.offerRepo.GetByID(id)
}

func (uc *offerUsecase) List(page, limit int) ([]*domain.Offer, error) {
	if page < 1 {
		return nil, ErrInvalidInput
	}
	if limit < 1 || limit >= 100 {
		return nil, ErrInvalidInput
	}
	return uc.offerRepo.List(page, limit)
}

func (uc *offerUsecase) Update(offer *domain.Offer) error {
	if offer.Title == "" {
		return ErrInvalidInput
	}
	return uc.offerRepo.Update(offer)
}

func (uc *offerUsecase) ListByUserID(userID string) ([]*domain.Offer, error) {
	return uc.offerRepo.ListByUserID(userID)
}


func (uc *offerUsecase) CountAll() (int, error) {
	return uc.offerRepo.CountAll()
}

func (uc *offerUsecase) Delete(id string) error {
	return uc.offerRepo.Delete(id)
}