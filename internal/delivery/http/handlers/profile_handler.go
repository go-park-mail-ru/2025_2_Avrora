package handlers

import (
	"context"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type IProfileUsecase interface {
	UpdateProfile(ctx context.Context, userID string, profile *domain.ProfileUpdate) error
	GetProfileByID(ctx context.Context, userID string) (*domain.Profile, error)
	UpdateProfileSecurityByID(ctx context.Context, userID string, oldPassword, newPassword string) error
	UpdateEmail(ctx context.Context, userID string, email string) error
}

type profileHandler struct {
	profileUsecase IProfileUsecase
	log     *log.Logger
}

func NewProfileHandler(uc IProfileUsecase, log *log.Logger) *profileHandler {
	return &profileHandler{profileUsecase: uc, log: log}
}