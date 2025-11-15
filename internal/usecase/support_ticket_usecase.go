package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
)

type ISupportTicketRepository interface {
	Create(ctx context.Context, ticket *domain.SupportTicket) error
	GetByID(ctx context.Context, id string) (*domain.SupportTicket, error)
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]domain.SupportTicket, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
	UpdateStatus(ctx context.Context, ticketID string, status domain.SupportTicketStatus) error
	Delete(ctx context.Context, ticketID string) error
	ListAll(ctx context.Context, page, limit int) ([]domain.SupportTicket, error)
	CountAll(ctx context.Context) (int, error)
}

type supportTicketUsecase struct {
	supportTicketRepo ISupportTicketRepository
	log *log.Logger
}

func NewSupportTicketUsecase(supportTicketRepo ISupportTicketRepository, log *log.Logger) *supportTicketUsecase {
	return &supportTicketUsecase{
		supportTicketRepo: supportTicketRepo,
		log: log,
	}
}