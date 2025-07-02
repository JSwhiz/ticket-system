package repository

import (
	"fmt"
	"ticket-system/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) *TicketRepository {
	return &TicketRepository{db}
}

func (r *TicketRepository) CreateTicket(ticket *models.Ticket) error {
	query := `
		INSERT INTO tickets (title, description, status_id, priority_id, creator_id, assignee_id, department_id)
		VALUES (:title, :description, :status_id, :priority_id, :creator_id, :assignee_id, :department_id)
	`
	_, err := r.db.NamedExec(query, ticket)
	return err
}
func (r *TicketRepository) GetAllTickets() ([]*models.Ticket, error) {
	tickets := []*models.Ticket{}
	err := r.db.Select(&tickets, "SELECT ticket_id, title, description, status_id, priority_id, creator_id, assignee_id, department_id, created_at, updated_at, deleted_at, search_vector FROM tickets WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) GetTicketByID(ticketID string) (*models.Ticket, error) {
	parsedID, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, err
	}
	ticket := &models.Ticket{}
	err = r.db.Get(ticket, "SELECT ticket_id, title, description, status_id, priority_id, creator_id, assignee_id, department_id, created_at, updated_at, deleted_at, search_vector FROM tickets WHERE ticket_id = $1 AND deleted_at IS NULL", parsedID)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *TicketRepository) UpdateTicket(ticketID string, updates map[string]interface{}) error {
	query := "UPDATE tickets SET "
	var args []interface{}

	for key, value := range updates {
		query += fmt.Sprintf("%s = $%d, ", key, len(args)+1)
		args = append(args, value)
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE ticket_id = $%d AND deleted_at IS NULL", len(args)+1)
	args = append(args, ticketID)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *TicketRepository) DeleteTicket(ticketID string) error {
	query := "UPDATE tickets SET deleted_at = NOW() WHERE ticket_id = $1 AND deleted_at IS NULL"
	_, err := r.db.Exec(query, ticketID)
	if err != nil {
		return err
	}
	return nil
}