package domain

import (
	"errors"
	"time"
)

type SupportTicket struct {
    ID            string                // Unique identifier for the ticket
    UserID        *string               // Nullable user ID due to ON DELETE SET NULL
    SignedEmail   string                // Email address provided by the user
    ResponseEmail string                // Email address where responses should be sent
    Name          string                // Name of the user or subject of the ticket
    Category      SupportTicketCategory // Category of the support ticket
    Description   string                // Detailed description of the issue
    Status        SupportTicketStatus   // Current status of the ticket
	PhotoURLs     []string 
    CreatedAt     time.Time             // Timestamp when the ticket was created
    UpdatedAt     time.Time             // Timestamp when the ticket was last updated
}

type SupportTicketCategory string

const (
	BugCategory      SupportTicketCategory = "bug"
	GeneralCategory   SupportTicketCategory = "general"
	BillingCategory   SupportTicketCategory = "billing"
	TechnicalCategory SupportTicketCategory = "feature"
)

// SupportTicketStatus represents the status of the ticket.
type SupportTicketStatus string

const (
	OpenStatus       SupportTicketStatus = "open"
	InProgressStatus SupportTicketStatus = "in_progress"
	ClosedStatus     SupportTicketStatus = "closed"
)

var (
    ErrSupportTicketNotFound = errors.New("support ticket not found")
    
)