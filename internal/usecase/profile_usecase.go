package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type IProfileRepository interface {
	GetByUserID(ctx context.Context, userID string) (*domain.Profile, string, error)
	Update(ctx context.Context, userID string, upd *domain.ProfileUpdate) error
	UpdateSecurity(ctx context.Context, userID string, passwordHash string) error
	UpdateEmail(ctx context.Context, userID string, email string) error
	GetUserByUserID(ctx context.Context, userID string) (*domain.User, error)
}

type profileUsecase struct {
	profileRepo IProfileRepository
	passwordHasher IPasswordHasher
	log  *log.Logger
}

func NewProfileUsecase(profileRepo IProfileRepository, ph IPasswordHasher, log *log.Logger) *profileUsecase {
	return &profileUsecase{
		profileRepo: profileRepo, 
		passwordHasher: ph, 
		log: log,
	}
}