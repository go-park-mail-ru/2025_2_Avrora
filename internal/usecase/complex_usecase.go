package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type IComplexRepository interface {
	GetByID(ctx context.Context, id string) (*domain.HousingComplex, error)
	List(ctx context.Context, page, limit int) (*domain.ComplexesInFeed, error)
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
