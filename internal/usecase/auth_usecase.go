package usecase

import (
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
)

type IUserRepository interface {
	Create(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
}

type IPasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) bool
}

type IJWTGenerator interface {
	GenerateJWT(id string) (string, error)
	GenerateExpiredJWT() (string, error)
}

type authUsecase struct {
	userRepo       IUserRepository
	passwordHasher IPasswordHasher
	jwtService     IJWTGenerator
}

func NewAuthUsecase(
	userRepo IUserRepository,
	hasher IPasswordHasher,
	jwt IJWTGenerator,
) *authUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		passwordHasher: hasher,
		jwtService:     jwt,
	}
}