package usecase

import (
	"errors"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
)

var (
	ErrUserAlreadyExists   = errors.New("пользователь с таким email уже существует")
	ErrInvalidCredentials  = errors.New("неправильные email или пароль")
	ErrInvalidInput        = errors.New("невалидные данные")
	ErrServerSideError     = errors.New("серверная ошибка")
)

func (uc *authUsecase) Register(email, password string) error {
	_, err := uc.userRepo.GetUserByEmail(email)
	if err == nil {
		return ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}

	hashed, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return err
	}

	user := domain.User{
		Email:     email,
		Password:  hashed,
	}

	return uc.userRepo.Create(&user)
}

func (uc *authUsecase) Login(email, password string) (string, error) {
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if !uc.passwordHasher.Compare(user.Password, password) {
		return "", ErrInvalidCredentials
	}

	return uc.jwtService.GenerateJWT(strconv.Itoa(user.ID))
}

func (uc *authUsecase) Logout() (string, error) {
	expiredToken, err := uc.jwtService.GenerateExpiredJWT()
	if err != nil {
		return "", err
	}
	return expiredToken, nil
}