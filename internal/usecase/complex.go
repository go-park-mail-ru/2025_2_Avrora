package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
)

func (u *housingComplexUsecase) GetByID(ctx context.Context, id string) (*domain.HousingComplex, error) {
	return u.complexRepo.GetByID(ctx, id)
}

func (u *housingComplexUsecase) List(ctx context.Context, page, limit int) (*domain.ComplexesInFeed, error) {
	return u.complexRepo.List(ctx, page, limit)
}

func (u *housingComplexUsecase) Create(ctx context.Context, complex *domain.HousingComplex) error {
	return u.complexRepo.Create(ctx, complex)
}

func (u *housingComplexUsecase) Update(ctx context.Context, complex *domain.HousingComplex) error {
	return u.complexRepo.Update(ctx, complex)
}

func (u *housingComplexUsecase) Delete(ctx context.Context, id string) error {
	return u.complexRepo.Delete(ctx, id)
}