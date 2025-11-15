package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateSupportTicketInput defines the input parameters for creating a support ticket
type CreateSupportTicketInput struct {
	UserID        *string                      `json:"user_id,omitempty"`
	SignedEmail   string                       `json:"signed_email" validate:"required,email,max=255"`
	ResponseEmail string                       `json:"response_email" validate:"required,email,max=255"`
	Name          string                       `json:"name" validate:"required,min=1,max=255"`
	Category      domain.SupportTicketCategory `json:"category" validate:"required,oneof=bug general billing feature"`
	Description   string                       `json:"description" validate:"required,min=1,max=5000"`
	PhotoURLs     []string                     `json:"photo_urls,omitempty"`
}

// CreateSupportTicket creates a new support ticket
func (uc *supportTicketUsecase) CreateSupportTicket(ctx context.Context, input CreateSupportTicketInput) (*domain.SupportTicket, error) {
	// Generate new UUID for the ticket
	ticketID, err := uuid.NewRandom()
	if err != nil {
		uc.log.Error(ctx, "failed to generate UUID for support ticket", zap.Error(err))
		return nil, errors.New("failed to generate ticket ID")
	}

	// Create the ticket domain object
	ticket := &domain.SupportTicket{
		ID:            ticketID.String(),
		UserID:        input.UserID,
		SignedEmail:   input.SignedEmail,
		ResponseEmail: input.ResponseEmail,
		Name:          input.Name,
		Category:      input.Category,
		Description:   input.Description,
		Status:        domain.OpenStatus, // Default status
		PhotoURLs:     input.PhotoURLs,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	// Validate the ticket data
	if err := uc.validateTicket(ticket); err != nil {
		uc.log.Warn(ctx, "invalid support ticket data", zap.Error(err))
		return nil, err
	}

	// Create the ticket in repository
	if err := uc.supportTicketRepo.Create(ctx, ticket); err != nil {
		uc.log.Error(ctx, "failed to create support ticket in repository",
			zap.String("ticket_id", ticket.ID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create support ticket: %w", err)
	}

	uc.log.Info(ctx, "successfully created support ticket",
		zap.String("ticket_id", ticket.ID),
		zap.String("user_id", fmt.Sprintf("%v", ticket.UserID)))

	return ticket, nil
}

// GetSupportTicketByID retrieves a support ticket by its ID
func (uc *supportTicketUsecase) GetSupportTicketByID(ctx context.Context, ticketID string) (*domain.SupportTicket, error) {
	ticket, err := uc.supportTicketRepo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, domain.ErrSupportTicketNotFound) {
			uc.log.Warn(ctx, "support ticket not found",
				zap.String("ticket_id", ticketID))
			return nil, domain.ErrSupportTicketNotFound
		}
		uc.log.Error(ctx, "failed to get support ticket by ID",
			zap.String("ticket_id", ticketID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get support ticket: %w", err)
	}

	uc.log.Info(ctx, "successfully retrieved support ticket",
		zap.String("ticket_id", ticket.ID))
	return ticket, nil
}

// GetSupportTicketsByUserID retrieves paginated support tickets for a user
func (uc *supportTicketUsecase) GetSupportTicketsByUserID(ctx context.Context, userID string, page, limit int) ([]domain.SupportTicket, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Maximum limit for performance
	}

	// Get tickets
	tickets, err := uc.supportTicketRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		uc.log.Error(ctx, "failed to get support tickets by user ID",
			zap.String("user_id", userID),
			zap.Int("page", page),
			zap.Int("limit", limit),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get support tickets: %w", err)
	}

	// Get total count
	totalCount, err := uc.supportTicketRepo.CountByUserID(ctx, userID)
	if err != nil {
		uc.log.Warn(ctx, "failed to count support tickets by user ID",
			zap.String("user_id", userID),
			zap.Error(err))
		totalCount = len(tickets)
	}
	return tickets, totalCount, nil
}

// UpdateSupportTicketStatus updates the status of a support ticket
func (uc *supportTicketUsecase) UpdateSupportTicketStatus(ctx context.Context, ticketID string, status domain.SupportTicketStatus) error {
	// Validate status transition
	if err := uc.validateStatusTransition(ctx, ticketID, status); err != nil {
		return err
	}

	if err := uc.supportTicketRepo.UpdateStatus(ctx, ticketID, status); err != nil {
		if errors.Is(err, domain.ErrSupportTicketNotFound) {
			uc.log.Warn(ctx, "support ticket not found for status update",
				zap.String("ticket_id", ticketID),
				zap.String("new_status", string(status)))
			return domain.ErrSupportTicketNotFound
		}
		uc.log.Error(ctx, "failed to update support ticket status",
			zap.String("ticket_id", ticketID),
			zap.String("new_status", string(status)),
			zap.Error(err))
		return fmt.Errorf("failed to update support ticket status: %w", err)
	}
	// TODO: Consider sending notification to user about status change
	return nil
}

// DeleteSupportTicket deletes a support ticket
func (uc *supportTicketUsecase) DeleteSupportTicket(ctx context.Context, ticketID string) error {
	if err := uc.supportTicketRepo.Delete(ctx, ticketID); err != nil {
		if errors.Is(err, domain.ErrSupportTicketNotFound) {
			uc.log.Warn(ctx, "support ticket not found for deletion",
				zap.String("ticket_id", ticketID))
			return domain.ErrSupportTicketNotFound
		}
		uc.log.Error(ctx, "failed to delete support ticket",
			zap.String("ticket_id", ticketID),
			zap.Error(err))
		return fmt.Errorf("failed to delete support ticket: %w", err)
	}
	// TODO: Consider cleanup of associated photo files from storage
	return nil
}

// ListAllSupportTickets retrieves all support tickets for admin view
func (uc *supportTicketUsecase) ListAllSupportTickets(ctx context.Context, page, limit int) ([]domain.SupportTicket, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 50 { // Lower limit for admin view to prevent performance issues
		limit = 50
	}

	tickets, err := uc.supportTicketRepo.ListAll(ctx, page, limit)
	if err != nil {
		uc.log.Error(ctx, "failed to list all support tickets",
			zap.Int("page", page),
			zap.Int("limit", limit),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list support tickets: %w", err)
	}

	totalCount, err := uc.supportTicketRepo.CountAll(ctx)
	if err != nil {
		uc.log.Warn(ctx, "failed to count all support tickets", zap.Error(err))
		totalCount = len(tickets)
	}
	return tickets, totalCount, nil
}

// validateTicket validates the support ticket data before creation
func (uc *supportTicketUsecase) validateTicket(ticket *domain.SupportTicket) error {
	// Basic validations (repository has constraints but we validate early)
	if ticket.SignedEmail == "" || ticket.ResponseEmail == "" {
		return errors.New("email fields cannot be empty")
	}

	if len(ticket.SignedEmail) > 255 || len(ticket.ResponseEmail) > 255 {
		return errors.New("email fields exceed maximum length")
	}

	if !isValidEmail(ticket.SignedEmail) || !isValidEmail(ticket.ResponseEmail) {
		return errors.New("invalid email format")
	}

	if ticket.Name == "" || len(ticket.Name) > 255 {
		return errors.New("invalid name field")
	}

	if ticket.Description == "" || len(ticket.Description) > 5000 {
		return errors.New("invalid description field")
	}

	// Validate category
	validCategories := map[domain.SupportTicketCategory]bool{
		domain.BugCategory:       true,
		domain.GeneralCategory:   true,
		domain.BillingCategory:   true,
		domain.TechnicalCategory: true,
	}
	if !validCategories[ticket.Category] {
		return errors.New("invalid ticket category")
	}

	// Validate status
	validStatuses := map[domain.SupportTicketStatus]bool{
		domain.OpenStatus:       true,
		domain.InProgressStatus: true,
		domain.ClosedStatus:     true,
	}
	if !validStatuses[ticket.Status] {
		return errors.New("invalid ticket status")
	}

	return nil
}

// validateStatusTransition validates if a status transition is allowed
func (uc *supportTicketUsecase) validateStatusTransition(ctx context.Context, ticketID string, newStatus domain.SupportTicketStatus) error {
	// Get current ticket to check status transition rules
	currentTicket, err := uc.supportTicketRepo.GetByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, domain.ErrSupportTicketNotFound) {
			return domain.ErrSupportTicketNotFound
		}
		return fmt.Errorf("failed to get current ticket status: %w", err)
	}

	// Define allowed transitions
	allowedTransitions := map[domain.SupportTicketStatus][]domain.SupportTicketStatus{
		domain.OpenStatus: {
			domain.InProgressStatus,
			domain.ClosedStatus,
		},
		domain.InProgressStatus: {
			domain.ClosedStatus,
			domain.OpenStatus, // Can revert to open
		},
		domain.ClosedStatus: {
			// Typically no transitions from closed, but could allow reopening
			domain.OpenStatus,
		},
	}

	allowed := false
	for _, allowedStatus := range allowedTransitions[currentTicket.Status] {
		if allowedStatus == newStatus {
			allowed = true
			break
		}
	}

	if !allowed {
		uc.log.Warn(ctx, "invalid status transition attempt",
			zap.String("ticket_id", ticketID),
			zap.String("current_status", string(currentTicket.Status)),
			zap.String("new_status", string(newStatus)))
		return fmt.Errorf("invalid status transition from %s to %s", currentTicket.Status, newStatus)
	}

	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	// Simple validation - could be enhanced with regex or proper email validation library
	if len(email) < 3 || len(email) > 255 {
		return false
	}

	atIndex := -1
	for i, char := range email {
		if char == '@' {
			atIndex = i
			break
		}
	}

	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return false
	}

	dotIndex := -1
	for i := atIndex + 1; i < len(email); i++ {
		if email[i] == '.' {
			dotIndex = i
			break
		}
	}

	return dotIndex != -1 && dotIndex < len(email)-1 && dotIndex > atIndex+1
}

// SupportTicketService provides methods for support ticket operations
type SupportTicketService interface {
	CreateSupportTicket(ctx context.Context, input CreateSupportTicketInput) (*domain.SupportTicket, error)
	GetSupportTicketByID(ctx context.Context, ticketID string) (*domain.SupportTicket, error)
	GetSupportTicketsByUserID(ctx context.Context, userID string, page, limit int) ([]domain.SupportTicket, int, error)
	UpdateSupportTicketStatus(ctx context.Context, ticketID string, status domain.SupportTicketStatus) error
	DeleteSupportTicket(ctx context.Context, ticketID string) error
	ListAllSupportTickets(ctx context.Context, page, limit int) ([]domain.SupportTicket, int, error)
}

// Ensure supportTicketUsecase implements SupportTicketService interface
var _ SupportTicketService = (*supportTicketUsecase)(nil)
