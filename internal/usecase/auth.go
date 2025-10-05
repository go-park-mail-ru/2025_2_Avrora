package usecase

import (
	"errors"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
)

var (
	ErrUserAlreadyExists   = errors.New("user with this email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidInput        = errors.New("invalid input")
)

func (uc *authUsecase) Register(email, password string) error {
	err := validateRegisterRequest(&RegisterRequest{Email: email, Password: password})
	if err != nil {
		return err
	}

	_, err = uc.userRepo.GetByEmail(email)
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

	return uc.userRepo.Create(user)
}

func (uc *authUsecase) Login(email, password string) (string, error) {
	err := validateLoginRequest(&LoginRequest{Email: email, Password: password})
	if err != nil {
		return "", err
	}

	user, err := uc.userRepo.GetByEmail(email)
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