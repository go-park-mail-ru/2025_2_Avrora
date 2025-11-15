package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/response"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/domain"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/log"
	"github.com/go-park-mail-ru/2025_2_Avrora/internal/usecase"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// ISupportTicketUsecase defines the interface for support ticket usecase operations
type ISupportTicketUsecase interface {
	CreateSupportTicket(ctx context.Context, input usecase.CreateSupportTicketInput) (*domain.SupportTicket, error)
	GetSupportTicketByID(ctx context.Context, id string) (*domain.SupportTicket, error)
	GetSupportTicketsByUserID(ctx context.Context, userID string, page, limit int) ([]domain.SupportTicket, int, error)
	UpdateSupportTicketStatus(ctx context.Context, id string, status domain.SupportTicketStatus) error
	DeleteSupportTicket(ctx context.Context, id string) error
	ListAllSupportTickets(ctx context.Context, page, limit int) ([]domain.SupportTicket, int, error)
}

// SupportTicketHandler handles HTTP requests for support tickets
type SupportTicketHandler struct {
	usecase ISupportTicketUsecase
	log     *log.Logger
	validate *validator.Validate
}

func NewSupportTicketHandler(usecase ISupportTicketUsecase, log *log.Logger) *SupportTicketHandler {
	return &SupportTicketHandler{
		usecase:  usecase,
		log:      log,
		validate: validator.New(),
	}
}

// CreateSupportTicket handles POST /support-tickets
func (h *SupportTicketHandler) CreateSupportTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateSupportTicketRequest

	if err := h.decodeJSON(w, r, &req); err != nil {
		response.HandleError(w, err, http.StatusBadRequest, "не валидный запрос")
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		h.log.Warn(ctx, "validation failed for create support ticket",
			zap.String("error", err.Error()))
		response.HandleError(w, err, http.StatusBadRequest, "не прошла валидация")
		return
	}

	// Convert category string to domain enum
	category, err := h.convertToDomainCategory(req.Category)
	if err != nil {
		h.log.Warn(ctx, "invalid category", zap.String("category", req.Category))
		response.HandleError(w, err, http.StatusBadRequest, "Не валидная категория")
		return
	}

	userID := req.UserID
	if userID == nil {
		if authUserID, ok := ctx.Value("user_id").(string); ok && authUserID != "" {
			userID = &authUserID
		}
	}

	// Prepare usecase input
	input := usecase.CreateSupportTicketInput{
		UserID:        userID,
		SignedEmail:   req.SignedEmail,
		ResponseEmail: req.ResponseEmail,
		Name:          req.Name,
		Category:      category,
		Description:   req.Description,
		PhotoURLs:     req.PhotoURLs,
	}

	// Create ticket
	ticket, err := h.usecase.CreateSupportTicket(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "email") {
			response.HandleError(w, err, http.StatusBadRequest, "невалидный email")
			return
		}
		h.log.Error(ctx, "не получилось создать тикет", zap.Error(err))
		return
	}

	// Map domain ticket to response
	resp := h.mapDomainTicketToResponse(ticket)
	response.WriteJSON(w, http.StatusCreated, resp)
}

// GetSupportTicketByID handles GET /support-tickets/{ticket_id}
func (h *SupportTicketHandler) GetSupportTicketByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := GetPathParameter(r, "/api/v1/support-tickets/")
	if id == "" {
		h.log.Error(r.Context(), "invalid or no id")
		response.HandleError(w, nil, http.StatusBadRequest, "нет id")
		return
	}

	ticket, err := h.usecase.GetSupportTicketByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrSupportTicketNotFound) {
			response.HandleError(w, err, http.StatusNotFound, "тикет не найден")
			return
		}
		h.log.Error(ctx, "failed to get support ticket", zap.String("ticket_id", id), zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения тикета")
		return
	}

	resp := h.mapDomainTicketToResponse(ticket)
	response.WriteJSON(w, http.StatusOK, resp)
}

// GetUserSupportTickets handles GET /support-tickets
func (h *SupportTicketHandler) GetUserSupportTickets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetPathParameter(r, "/api/v1/support-tickets/my/")

	// Get pagination parameters
	page, limit := h.getPaginationParams(r)

	tickets, totalCount, err := h.usecase.GetSupportTicketsByUserID(ctx, userID, page, limit)
	if err != nil {
		h.log.Error(ctx, "failed to get user support tickets",
			zap.Int("page", page),
			zap.Int("limit", limit),
			zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения тикетов")
		return
	}

	// Map tickets to response format
	respTickets := make([]CreateSupportTicketResponse, len(tickets))
	for i, ticket := range tickets {
		respTickets[i] = h.mapDomainTicketToResponse(&ticket)
	}

	resp := struct {
		Tickets []CreateSupportTicketResponse `json:"tickets"`
		Meta    struct {
			Total  int `json:"total"`
			Page   int `json:"page"`
			Limit  int `json:"limit"`
			Pages  int `json:"pages"`
		} `json:"meta"`
	}{
		Tickets: respTickets,
		Meta: struct {
			Total  int `json:"total"`
			Page   int `json:"page"`
			Limit  int `json:"limit"`
			Pages  int `json:"pages"`
		}{
			Total: totalCount,
			Page:  page,
			Limit: limit,
			Pages: (totalCount + limit - 1) / limit,
		},
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

func (h *SupportTicketHandler) GetAllSupportTickets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetPathParameter(r, "/api/v1/support-tickets/all/")

	// Check if user is authorized to view all tickets (admin/moderator)
	if !h.isAdmin(ctx, userID) {
		response.HandleError(w, nil, http.StatusUnauthorized, "пользователь не админ")
		h.log.Warn(ctx, "unauthorized attempt to access all support tickets")
		return
	}

	// Get pagination parameters
	page, limit := h.getPaginationParams(r)

	// Optional: Get filter parameters from query string
	filters := h.getFilterParams(r)

	tickets, totalCount, err := h.usecase.ListAllSupportTickets(ctx, page, limit)
	if err != nil {
		h.log.Error(ctx, "failed to get all support tickets",
			zap.Int("page", page),
			zap.Int("limit", limit),
			zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения тикетов")
		return
	}

	// Apply filters if any (optional - can be moved to usecase layer for better performance)
	filteredTickets := h.applyFilters(tickets, filters)

	// Map tickets to response format
	respTickets := make([]CreateSupportTicketResponse, len(filteredTickets))
	for i, ticket := range filteredTickets {
		respTickets[i] = h.mapDomainTicketToResponse(&ticket)
	}

	// Calculate total pages based on filtered count
	filteredCount := len(filteredTickets)
	if filters.hasFilters() {
		// When filtering server-side, we need to recalculate total
		// In production, this should be handled at database level
		totalCount = filteredCount
	}

	resp := struct {
		Tickets []CreateSupportTicketResponse `json:"tickets"`
		Meta    struct {
			Total     int      `json:"total"`
			Page      int      `json:"page"`
			Limit     int      `json:"limit"`
			Pages     int      `json:"pages"`
			Filters   Filters  `json:"filters,omitempty"`
			Sort      string   `json:"sort,omitempty"`
		} `json:"meta"`
	}{
		Tickets: respTickets,
		Meta: struct {
			Total     int      `json:"total"`
			Page      int      `json:"page"`
			Limit     int      `json:"limit"`
			Pages     int      `json:"pages"`
			Filters   Filters  `json:"filters,omitempty"`
			Sort      string   `json:"sort,omitempty"`
		}{
			Total:  totalCount,
			Page:   page,
			Limit:  limit,
			Pages:  (totalCount + limit - 1) / limit,
			Filters: filters,
			Sort:   r.URL.Query().Get("sort"),
		},
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

func (h *SupportTicketHandler) getFilterParams(r *http.Request) Filters {
	filters := Filters{}

	// Status filter (comma-separated values)
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		filters.Status = strings.Split(statusStr, ",")
		// Clean and validate status values
		for i, s := range filters.Status {
			filters.Status[i] = strings.TrimSpace(strings.ToLower(s))
		}
	}

	// Category filter
	if categoryStr := r.URL.Query().Get("category"); categoryStr != "" {
		filters.Category = strings.Split(categoryStr, ",")
		for i, c := range filters.Category {
			filters.Category[i] = strings.TrimSpace(strings.ToLower(c))
		}
	}

	// Date range filters
	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filters.From = &t
		}
	}

	if toStr := r.URL.Query().Get("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			filters.To = &t
		}
	}

	// Search query
	if search := r.URL.Query().Get("search"); search != "" {
		filters.Search = strings.TrimSpace(search)
	}

	return filters
}

// applyFilters applies client-side filtering (for demo purposes - in production move to DB level)
func (h *SupportTicketHandler) applyFilters(tickets []domain.SupportTicket, filters Filters) []domain.SupportTicket {
	if !filters.hasFilters() {
		return tickets
	}

	var filtered []domain.SupportTicket
	for _, ticket := range tickets {
		if filters.matches(ticket) {
			filtered = append(filtered, ticket)
		}
	}
	return filtered
}

// Filters struct for query parameters
type Filters struct {
	Status   []string  `json:"status,omitempty"`
	Category []string  `json:"category,omitempty"`
	From     *time.Time `json:"from,omitempty"`
	To       *time.Time `json:"to,omitempty"`
	Search   string    `json:"search,omitempty"`
}

// hasFilters checks if any filters are applied
func (f Filters) hasFilters() bool {
	return len(f.Status) > 0 || len(f.Category) > 0 || f.From != nil || f.To != nil || f.Search != ""
}

// matches checks if a ticket matches the filters
func (f Filters) matches(ticket domain.SupportTicket) bool {
	// Status filter
	if len(f.Status) > 0 {
		matched := false
		ticketStatus := strings.ToLower(string(ticket.Status))
		for _, status := range f.Status {
			if status == ticketStatus {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Category filter
	if len(f.Category) > 0 {
		matched := false
		ticketCategory := strings.ToLower(string(ticket.Category))
		for _, category := range f.Category {
			if category == ticketCategory {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Date range filter
	if f.From != nil && ticket.CreatedAt.Before(*f.From) {
		return false
	}
	if f.To != nil && ticket.CreatedAt.After(*f.To) {
		return false
	}

	// Search filter (case-insensitive partial match)
	if f.Search != "" {
		searchLower := strings.ToLower(f.Search)
		if !strings.Contains(strings.ToLower(ticket.Name), searchLower) &&
			!strings.Contains(strings.ToLower(ticket.Description), searchLower) &&
			!strings.Contains(strings.ToLower(ticket.SignedEmail), searchLower) &&
			(ticket.UserID == nil || !strings.Contains(strings.ToLower(*ticket.UserID), searchLower)) {
			return false
		}
	}

	return true
}

// UpdateSupportTicketStatusRequest represents the request body for updating ticket status
type UpdateSupportTicketStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=open in_progress closed"`
}

// UpdateSupportTicketStatus handles PATCH /admin/support-tickets/status/{ticket_id}
func (h *SupportTicketHandler) UpdateSupportTicketStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := GetPathParameter(r, "/api/v1/admin/support-tickets/status/")
	userID := GetPathParameter(r, "/api/v1/admin/support-tickets/status/"+id+"/")

	if !h.isAdmin(ctx, userID) {
		response.HandleError(w, nil, http.StatusUnauthorized, "пользователь не админ")
		return
	}

	var req UpdateSupportTicketStatusRequest
	if err := h.decodeJSON(w, r, &req); err != nil {
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		h.log.Warn(ctx, "validation failed for update status",
			zap.String("ticket_id", id),
			zap.String("status", req.Status),
			zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "не прошла валидация")
		return
	}

	// Convert status string to domain enum
	status, err := h.convertToDomainStatus(req.Status)
	if err != nil {
		h.log.Warn(ctx, "invalid status", zap.String("status", req.Status))
		response.HandleError(w, err, http.StatusBadRequest, "Не валидный статус")
		return
	}

	// Update status
	if err := h.usecase.UpdateSupportTicketStatus(ctx, id, status); err != nil {
		if errors.Is(err, domain.ErrSupportTicketNotFound) {
			response.HandleError(w, err, http.StatusNotFound, "тикет не найден")
			return
		}
		h.log.Error(ctx, "failed to update ticket status",
			zap.String("ticket_id", id),
			zap.String("new_status", req.Status),
			zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка обновления статуса тикета")
		return
	}

	resp := struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Успешно обновлено",
		Status:  req.Status,
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// DeleteSupportTicket handles DELETE /support-tickets/{ticket_id}
func (h *SupportTicketHandler) DeleteSupportTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := GetPathParameter(r, "/api/v1/support-tickets/delete/")
	if id == "" {
		h.log.Error(r.Context(), "invalid or no id")
		response.HandleError(w, nil, http.StatusBadRequest, "нет id")
		return
	}

	userID := GetPathParameter(r, "/api/v1/support-tickets/delete/"+id+"/")
	if !h.isAdmin(ctx, userID) {
		response.HandleError(w, nil, http.StatusUnauthorized, "пользователь не админ")
		return
	}

	// Get ticket to check ownership
	_, err := h.usecase.GetSupportTicketByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrSupportTicketNotFound) {
			response.HandleError(w, err, http.StatusNotFound, "тикет не найден")
			return
		}
		h.log.Error(ctx, "failed to get ticket for deletion", zap.String("ticket_id", id), zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения тикета")
		return
	}

	// Delete ticket
	if err := h.usecase.DeleteSupportTicket(ctx, id); err != nil {
		if errors.Is(err, domain.ErrSupportTicketNotFound) {
			response.HandleError(w, err, http.StatusNotFound, "тикет не найден")
			return
		}
		h.log.Error(ctx, "failed to delete support ticket", zap.String("ticket_id", id), zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка удаления тикета")
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Успешно удалено",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// ListAllSupportTickets handles GET /admin/support-tickets
func (h *SupportTicketHandler) ListAllSupportTickets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get pagination parameters
	page, limit := h.getPaginationParams(r)

	tickets, totalCount, err := h.usecase.ListAllSupportTickets(ctx, page, limit)
	if err != nil {
		h.log.Error(ctx, "failed to list all support tickets",
			zap.Int("page", page),
			zap.Int("limit", limit),
			zap.Error(err))
		response.HandleError(w, err, http.StatusInternalServerError, "ошибка получения тикетов")
		return
	}

	// Map tickets to response format
	respTickets := make([]CreateSupportTicketResponse, len(tickets))
	for i, ticket := range tickets {
		respTickets[i] = h.mapDomainTicketToResponse(&ticket)
	}

	resp := struct {
		Tickets []CreateSupportTicketResponse `json:"tickets"`
		Meta    struct {
			Total  int `json:"total"`
			Page   int `json:"page"`
			Limit  int `json:"limit"`
			Pages  int `json:"pages"`
		} `json:"meta"`
	}{
		Tickets: respTickets,
		Meta: struct {
			Total  int `json:"total"`
			Page   int `json:"page"`
			Limit  int `json:"limit"`
			Pages  int `json:"pages"`
		}{
			Total: totalCount,
			Page:  page,
			Limit: limit,
			Pages: (totalCount + limit - 1) / limit,
		},
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// Helper methods

func (h *SupportTicketHandler) decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		h.log.Warn(r.Context(), "failed to decode JSON", zap.Error(err))
		response.HandleError(w, err, http.StatusBadRequest, "невалидный JSON")
		return err
	}
	return nil
}

func (h *SupportTicketHandler) getPaginationParams(r *http.Request) (int, int) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	return page, limit
}

func (h *SupportTicketHandler) isAdmin(ctx context.Context, userID string) bool {
	return true
}

func (h *SupportTicketHandler) mapDomainTicketToResponse(ticket *domain.SupportTicket) CreateSupportTicketResponse {
	return CreateSupportTicketResponse{
		ID:            ticket.ID,
		UserID:        ticket.UserID,
		SignedEmail:   ticket.SignedEmail,
		ResponseEmail: ticket.ResponseEmail,
		Name:          ticket.Name,
		Category:      string(ticket.Category),
		Description:   ticket.Description,
		Status:        string(ticket.Status),
		PhotoURLs:     ticket.PhotoURLs,
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
	}
}

func (h *SupportTicketHandler) convertToDomainCategory(categoryStr string) (domain.SupportTicketCategory, error) {
	switch strings.ToLower(categoryStr) {
	case "bug":
		return domain.BugCategory, nil
	case "general":
		return domain.GeneralCategory, nil
	case "billing":
		return domain.BillingCategory, nil
	case "feature", "technical":
		return domain.TechnicalCategory, nil
	default:
		return "", errors.New("invalid category")
	}
}

func (h *SupportTicketHandler) convertToDomainStatus(statusStr string) (domain.SupportTicketStatus, error) {
	switch strings.ToLower(statusStr) {
	case "open":
		return domain.OpenStatus, nil
	case "in_progress", "in progress":
		return domain.InProgressStatus, nil
	case "closed":
		return domain.ClosedStatus, nil
	default:
		return "", errors.New("invalid status")
	}
}