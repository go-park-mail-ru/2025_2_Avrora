package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"go.uber.org/zap"
)

var (
	ErrUserAlreadyExists   = errors.New("пользователь с таким email уже существует")
	ErrInvalidCredentials  = errors.New("неправильные email или пароль")
	ErrInvalidInput        = errors.New("невалидные данные")
	ErrServerSideError     = errors.New("серверная ошибка")
)

func (uc *authUsecase) Register(ctx context.Context, email, password string) error {
	_, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		uc.log.Error(ctx, "user already exists", zap.Error(err))
		return ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		uc.log.Error(ctx, "server side error", zap.Error(err))
		return err
	}

	hashed, err := uc.passwordHasher.Hash(password)
	if err != nil {
		uc.log.Error(ctx, "server side error", zap.Error(err))
		return err
	}

	user := domain.User{
		Email:     email,
		Password:  hashed,
	}

	uc.log.Info(ctx, "successfull sign up", zap.String("email", email))
	return uc.userRepo.Create(ctx, &user)
}

func (uc *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			uc.log.Error(ctx, "user not found", zap.Error(err))
			return "", ErrInvalidCredentials
		}
		uc.log.Error(ctx, "server side error", zap.Error(err))
		return "", err
	}

	if !uc.passwordHasher.Compare(user.Password, password) {
		uc.log.Error(ctx, "invalid credentials", zap.Error(err))
		return "", ErrInvalidCredentials
	}

	uc.log.Info(ctx, "successfull login", zap.String("email", email))
	return uc.jwtService.GenerateJWT(strconv.Itoa(user.ID))
}

func (uc *authUsecase) Logout(ctx context.Context) (string, error) {
	expiredToken, err := uc.jwtService.GenerateExpiredJWT()
	if err != nil {
		uc.log.Error(ctx, "server side error", zap.Error(err))
		return "", err
	}

	uc.log.Info(ctx, "successfull logout")
	return expiredToken, nil
}