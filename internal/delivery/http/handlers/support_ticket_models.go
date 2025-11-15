package handlers

import "time"

type SupportTicket struct {
	ID            string                `json:"id"`
	UserID        *string               `json:"user_id"`
	SignedEmail   string                `json:"signed_email"`
	ResponseEmail string                `json:"response_email"`
	Name          string                `json:"name"`
	Category      SupportTicketCategory `json:"category"`
	Description   string                `json:"description"`
	Status        SupportTicketStatus   `json:"status"`
	PhotoURLs     []string              `json:"photo_urls,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

type SupportTicketCategory string

const (
	GeneralCategory   SupportTicketCategory = "general"
	BillingCategory   SupportTicketCategory = "billing"
	TechnicalCategory SupportTicketCategory = "technical"
)

// SupportTicketStatus represents the status of the ticket.
type SupportTicketStatus string

const (
	OpenStatus       SupportTicketStatus = "open"
	InProgressStatus SupportTicketStatus = "in_progress"
	ResolvedStatus   SupportTicketStatus = "resolved"
	ClosedStatus     SupportTicketStatus = "closed"
)


type CreateSupportTicketRequest struct {
	UserID        *string               `json:"user_id,omitempty"`
	SignedEmail   string                `json:"signed_email" validate:"required,email,max=255"`
	ResponseEmail string                `json:"response_email" validate:"required,email,max=255"`
	Name          string                `json:"name" validate:"required,min=1,max=255"`
	Category      string                `json:"category" validate:"required,oneof=bug general billing feature"`
	Description   string                `json:"description" validate:"required,min=1,max=5000"`
	PhotoURLs     []string              `json:"photo_urls,omitempty"`
}

// CreateSupportTicketResponse represents the response for creating a support ticket
type CreateSupportTicketResponse struct {
	ID            string    `json:"id"`
	UserID        *string   `json:"user_id"`
	SignedEmail   string    `json:"signed_email"`
	ResponseEmail string    `json:"response_email"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	PhotoURLs     []string  `json:"photo_urls"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}