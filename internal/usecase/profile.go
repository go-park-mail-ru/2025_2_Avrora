package usecase

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"go.uber.org/zap"
)

func (uc *profileUsecase) GetProfileByID(ctx context.Context, userID string) (*domain.Profile, error) {
	if userID == "" {
		uc.log.Warn(ctx, "empty user ID in GetProfileByID")
		return nil, domain.ErrInvalidInput
	}

	profile, email, err := uc.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrProfileNotFound) {
			return &domain.Profile{
				UserID:    userID,
				Email:     email,
			}, nil
		}
		return nil, err
	}

	return &domain.Profile{
		ID:        profile.ID,
		UserID:    profile.UserID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Role:      profile.Role,
		Email:     email,
		Phone:     profile.Phone,
		AvatarURL: profile.AvatarURL,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}, nil
}

// UpdateProfileByID updates basic profile info (name, phone, avatar).
// Does NOT update email (that's a separate secure flow).
func (uc *profileUsecase) UpdateProfile(ctx context.Context, userID string, upd *domain.ProfileUpdate) error {
	if userID == "" {
		uc.log.Warn(ctx, "empty user ID in UpdateProfileByID")
		return domain.ErrInvalidInput
	}
	if upd == nil {
		uc.log.Warn(ctx, "nil update payload", zap.String("user_id", userID))
		return domain.ErrInvalidInput
	}

	// Optional: validate phone or avatar URL format here
	// Example: if upd.Phone != "" && !isValidPhone(upd.Phone) { ... }

	return uc.profileRepo.Update(ctx, userID, upd)
}

func (uc *profileUsecase) UpdateProfileSecurityByID(ctx context.Context, userID string, oldPassword, newPassword string) error {
	if userID == "" || oldPassword == "" || newPassword == "" {
		uc.log.Warn(ctx, "empty input in UpdateProfileSecurityByID", zap.String("user_id", userID))
		return domain.ErrInvalidInput
	}

	user, err := uc.profileRepo.GetUserByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	println(oldPassword, newPassword)
	if !uc.passwordHasher.Compare(user.PasswordHash, oldPassword) {
		uc.log.Error(ctx, "invalid credentials", zap.Error(err))
		return ErrInvalidCredentials
	}

	newPasswordHash, err := uc.passwordHasher.Hash(newPassword)
	if err != nil {
		uc.log.Error(ctx, "failed to hash password", zap.Error(err))
		return err
	}

	return uc.profileRepo.UpdateSecurity(ctx, userID, newPasswordHash)
}

func (uc *profileUsecase) UpdateEmail(ctx context.Context, userID string, email string) error {
	if userID == "" || email == "" {
		uc.log.Warn(ctx, "empty input in UpdateEmail", zap.String("user_id", userID))
		return domain.ErrInvalidInput
	}
	return uc.profileRepo.UpdateEmail(ctx, userID, email)
}