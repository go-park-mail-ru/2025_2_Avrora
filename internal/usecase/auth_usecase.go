package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
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
	log            *log.Logger
}

func NewAuthUsecase(
	userRepo IUserRepository,
	hasher IPasswordHasher,
	jwt IJWTGenerator,
	log *log.Logger,
) *authUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		passwordHasher: hasher,
		jwtService:     jwt,
		log:            log,
	}
}
