package usecase

import (
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/infrastructure/db"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/infrastructure/utils"
)

type AuthUsecase interface {
	Register(email, password string) error
	Login(email, password string) (string, error)
	Logout() (string, error)
}

type authUsecase struct {
	userRepo       *db.UserRepository
	passwordHasher *utils.PasswordHasher
	jwtService     *utils.JwtGenerator
}

var _ AuthUsecase = (*authUsecase)(nil)

func NewAuthUsecase(
	userRepo *db.UserRepository,
	hasher *utils.PasswordHasher,
	jwt *utils.JwtGenerator,
) AuthUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		passwordHasher: hasher,
		jwtService:     jwt,
	}
}