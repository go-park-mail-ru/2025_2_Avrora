package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type IComplexRepository interface {
	GetByID(ctx context.Context, id string) (*domain.HousingComplex, error)
	List(ctx context.Context, page, limit int) ([]*domain.HousingComplex, error)
	Create(ctx context.Context, complex *domain.HousingComplex) error
	Update(ctx context.Context, complex *domain.HousingComplex) error
	Delete(ctx context.Context, id string) error
}

type housingComplexUsecase struct {
	complexRepo IComplexRepository
	log         *log.Logger
}

func NewHousingComplexUsecase(
	complexRepo IComplexRepository,
	log *log.Logger,
) *housingComplexUsecase {
	return &housingComplexUsecase{
		complexRepo: complexRepo,
		log:         log,
	}
}

func (u *housingComplexUsecase) GetByID(ctx context.Context, id string) (*domain.HousingComplex, error) {
	return u.complexRepo.GetByID(ctx, id)
}

func (u *housingComplexUsecase) List(ctx context.Context, page, limit int) ([]*domain.HousingComplex, error) {
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