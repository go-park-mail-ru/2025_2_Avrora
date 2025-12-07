package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// === FEED METHODS ===

// ListOffersInFeed returns paginated offers for the main feed
func (uc *offerUsecase) ListOffersInFeed(ctx context.Context, page, limit int) (*domain.OffersInFeed, error) {
	if page < 1 {
		uc.log.Warn(ctx, "invalid page in feed", zap.Int("page", page))
		return nil, domain.ErrInvalidInput
	}
	if limit < 1 || limit > 100 {
		uc.log.Warn(ctx, "invalid limit in feed", zap.Int("limit", limit))
		return nil, domain.ErrInvalidInput
	}

	offers, err := uc.offerRepo.List(ctx, page, limit)
	if err != nil {
		uc.log.Error(ctx, "failed to list offers for feed", zap.Error(err))
		return nil, err
	}

	return offers, nil
}

// ListOffersInFeedByUserID returns paginated offers for a specific user (e.g., "my offers")
func (uc *offerUsecase) ListOffersInFeedByUserID(ctx context.Context, userID string, page, limit int) (*domain.OffersInFeed, error) {
	if userID == "" {
		uc.log.Warn(ctx, "empty user ID in ListOffersInFeedByUserID")
		return nil, domain.ErrInvalidInput
	}
	if page < 1 {
		uc.log.Warn(ctx, "invalid page", zap.Int("page", page))
		return nil, domain.ErrInvalidInput
	}
	if limit < 1 || limit > 100 {
		uc.log.Warn(ctx, "invalid limit", zap.Int("limit", limit))
		return nil, domain.ErrInvalidInput
	}

	offers, err := uc.offerRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		uc.log.Error(ctx, "failed to list user offers for feed", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	return offers, nil
}

// func (uc *offerUsecase) buildOffersInFeed(
// 	ctx context.Context,
// 	offers []*domain.Offer,
// 	total int,
// 	page, limit int,
// ) (*domain.OffersInFeed, error) {
// 	if len(offers) == 0 {
// 		return &domain.OffersInFeed{
// 			Meta: struct {
// 				Total  int
// 				Offset int
// 			}{Total: total, Offset: (page - 1) * limit},
// 			Offers: []domain.OfferInFeed{},
// 		}, nil
// 	}

// 	// Extract IDs
// 	offerIDs := make([]string, len(offers))
// 	for i, o := range offers {
// 		offerIDs[i] = o.ID
// 	}

// 	// Get first photo for each
// 	firstPhotos, err := uc.offerRepo.GetFirstPhotoForOffers(ctx, offerIDs)
// 	if err != nil {
// 		uc.log.Warn(ctx, "partial photo load in feed", zap.Error(err))
// 		// Proceed with empty photos
// 		if firstPhotos == nil {
// 			firstPhotos = make(map[string]string)
// 		}
// 	}

// 	// Map to OfferInFeed
// 	feedOffers := make([]domain.OfferInFeed, len(offers))
// 	for i, o := range offers {
// 		feedOffers[i] = domain.OfferInFeed{
// 			ID:           o.ID,
// 			UserID:       o.UserID,
// 			OfferType:    o.OfferType,
// 			PropertyType: o.PropertyType,
// 			Price:        o.Price,
// 			Area:         o.Area,
// 			Rooms:        o.Rooms,
// 			Floor:        safeIntDeref(o.Floor),
// 			TotalFloors:  safeIntDeref(o.TotalFloors),
// 			Address:      o.Address,
// 			Metro:        "", // enrich later if needed (e.g., via location join)
// 			ImageURL:     firstPhotos[o.ID],
// 			CreatedAt:    o.CreatedAt,
// 			UpdatedAt:    o.UpdatedAt,
// 		}
// 	}

// 	return &domain.OffersInFeed{
// 		Meta: struct {
// 			Total  int
// 			Offset int
// 		}{
// 			Total:  total,
// 			Offset: (page - 1) * limit,
// 		},
// 		Offers: feedOffers,
// 	}, nil
// }

// === DETAIL VIEW ===

func (uc *offerUsecase) Get(ctx context.Context, id string) (*domain.Offer, error) {
	if id == "" {
		uc.log.Warn(ctx, "empty offer ID")
		return nil, domain.ErrInvalidInput
	}

	offer, err := uc.offerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return offer, nil
}

// === CORE CRUD (unchanged) ===

func (uc *offerUsecase) Create(ctx context.Context, offer *domain.Offer) error {
	if offer == nil || offer.Title == "" {
		uc.log.Warn(ctx, "empty offer title")
		return domain.ErrInvalidInput
	}
	fmt.Println(offer.UserID)
	fmt.Println(offer.Price)
	fmt.Println(offer.Area)
	if offer.UserID == "" || offer.Price <= 0 || offer.Area <= 0 {
		uc.log.Warn(ctx, "invalid offer fields")
		return domain.ErrInvalidInput
	}
	offer.ID = uuid.NewString()
	return uc.offerRepo.Create(ctx, offer)
}

func (uc *offerUsecase) Update(ctx context.Context, offer *domain.Offer) error {
	if offer == nil || offer.ID == "" || offer.Title == "" {
		return domain.ErrInvalidInput
	}
	return uc.offerRepo.Update(ctx, offer)
}

func (uc *offerUsecase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidInput
	}
	return uc.offerRepo.Delete(ctx, id)
}

// ViewOffer records a view event for an offer (can be anonymous)
func (uc *offerUsecase) ViewOffer(ctx context.Context, offerID string) error {
	if offerID == "" {
		uc.log.Warn(ctx, "empty offer ID for view")
		return domain.ErrInvalidInput
	}

	return uc.offerRepo.LogView(ctx, offerID)
}

// ToggleOfferLike toggles like status for an authenticated user
func (uc *offerUsecase) ToggleOfferLike(ctx context.Context, offerID, userID string) error {
	if offerID == "" || userID == "" {
		uc.log.Warn(ctx, "missing offer or user ID for like toggle",
			zap.String("offer_id", offerID),
			zap.String("user_id", userID))
		return domain.ErrInvalidInput
	}
	
	return uc.offerRepo.ToggleLike(ctx, offerID, userID)
}

// GetOfferViewCount retrieves total views for an offer
func (uc *offerUsecase) GetOfferViewCount(ctx context.Context, offerID string) (int, error) {
	if offerID == "" {
		uc.log.Warn(ctx, "empty offer ID for view count")
		return 0, domain.ErrInvalidInput
	}
	
	return uc.offerRepo.GetOfferViewCount(ctx, offerID)
}

// GetOfferLikeCount retrieves total likes for an offer
func (uc *offerUsecase) GetOfferLikeCount(ctx context.Context, offerID string) (int, error) {
	if offerID == "" {
		uc.log.Warn(ctx, "empty offer ID for like count")
		return 0, domain.ErrInvalidInput
	}
	
	return uc.offerRepo.GetOfferLikeCount(ctx, offerID)
}

// IsOfferLiked checks if current user has liked an offer
func (uc *offerUsecase) IsOfferLiked(ctx context.Context, offerID, userID string) (bool, error) {
	if offerID == "" || userID == "" {
		uc.log.Warn(ctx, "missing offer or user ID for like check",
			zap.String("offer_id", offerID),
			zap.String("user_id", userID))
		
		// Return false for invalid inputs (safe default)
		return false, nil
	}
	
	return uc.offerRepo.IsOfferLiked(ctx, offerID, userID)
}

// InsertPaidAdvertisement inserts a new paid advertisement.
func (uc *offerUsecase) InsertPaidAdvertisement(ctx context.Context, offerID string, expiresAt time.Time) error {
	// Validate inputs
	if offerID == "" {
		uc.log.Warn(ctx, "empty offer ID for paid advertisement")
		return domain.ErrInvalidInput
	}
	if expiresAt.Before(time.Now()) {
		uc.log.Warn(ctx, "invalid expiration time for paid advertisement",
			zap.Time("expires_at", expiresAt))
		return domain.ErrInvalidInput
	}

	// Delegate to repository
	err := uc.offerRepo.InsertPaidAdvertisement(ctx, offerID, expiresAt)
	if err != nil {
		uc.log.Error(ctx, "failed to insert paid advertisement",
			zap.String("offer_id", offerID),
			zap.Time("expires_at", expiresAt),
			zap.Error(err))
		return fmt.Errorf("insert paid advertisement: %w", err)
	}

	return nil
}

// ListPaidOffers retrieves paginated paid offers in the OffersInFeed format.
func (uc *offerUsecase) ListPaidOffers(ctx context.Context, page, limit int) (*domain.OffersInFeed, error) {
	// Validate inputs
	if page <= 0 || limit <= 0 {
		uc.log.Warn(ctx, "invalid pagination parameters",
			zap.Int("page", page),
			zap.Int("limit", limit))
		return nil, domain.ErrInvalidInput
	}

	// Delegate to repository
	paidOffers, err := uc.offerRepo.ListPaidOffers(ctx, page, limit)
	if err != nil {
		uc.log.Error(ctx, "failed to list paid offers",
			zap.Int("page", page),
			zap.Int("limit", limit),
			zap.Error(err))
		return nil, fmt.Errorf("list paid offers: %w", err)
	}

	return paidOffers, nil
}