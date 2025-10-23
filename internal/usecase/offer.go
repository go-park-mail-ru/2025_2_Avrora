package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
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
		return domain.ErrInvalidInput
	}
	if offer.UserID == "" || offer.LocationID == "" || offer.Price <= 0 || offer.Area <= 0 {
		return domain.ErrInvalidInput
	}
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
