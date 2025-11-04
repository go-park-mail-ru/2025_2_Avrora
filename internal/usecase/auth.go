package usecase

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"go.uber.org/zap"
)

var (
	ErrUserAlreadyExists   = errors.New("пользователь с таким email уже существует")
	ErrInvalidCredentials  = errors.New("неправильные email или пароль")
	ErrInvalidInput        = errors.New("невалидные данные")
	ErrServerSideError     = errors.New("серверная ошибка")
	ErrUserNotFound        = errors.New("пользователь не найден")
)

func (uc *authUsecase) Register(ctx context.Context, email, password string) error {
	_, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		uc.log.Error(ctx, "user already exists", zap.String("email", email))
		return ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}

	hashed, err := uc.passwordHasher.Hash(password)
	if err != nil {
		uc.log.Error(ctx, "failed to hash password", zap.Error(err))
		return err
	}

	user := domain.User{
		Email:     email,
		PasswordHash:  hashed,
	}

	return uc.userRepo.Create(ctx, &user)
}

func (uc *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			uc.log.Error(ctx, "user not found", zap.String("email", email))
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if !uc.passwordHasher.Compare(user.PasswordHash, password) {
		uc.log.Error(ctx, "invalid credentials", zap.Error(err))
		return "", ErrInvalidCredentials
	}

	return uc.jwtService.GenerateJWT(user.ID)
}

func (uc *authUsecase) Logout(ctx context.Context) (string, error) {
	expiredToken, err := uc.jwtService.GenerateExpiredJWT()
	if err != nil {
		uc.log.Error(ctx, "failed to generate expired token", zap.Error(err))
		return "", err
	}

	return expiredToken, nil
}