package repository

import (
	"fmt"
	"log"
	"ticket-system/internal/models"
	"time"

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
	// Список допустимых полей для обновления
	allowedFields := map[string]bool{
		"title":       true,
		"description": true,
		"status_id":   true,
		"priority_id": true,
		"assignee_id": true,
	}

	// Проверка допустимых полей
	for key := range updates {
		if !allowedFields[key] {
			return fmt.Errorf("field %s is not allowed for update", key)
		}
	}

	// Автоматическое добавление updated_at
	updates["updated_at"] = time.Now()

	// Построение SQL-запроса
	query := "UPDATE tickets SET "
	var args []interface{}

	for key, value := range updates {
		query += fmt.Sprintf("%s = $%d, ", key, len(args)+1)
		args = append(args, value)
	}

	query = query[:len(query)-2] // Удаление последней запятой
	query += fmt.Sprintf(" WHERE ticket_id = $%d AND deleted_at IS NULL", len(args)+1)
	args = append(args, ticketID)

	// Выполнение запроса
	result, err := r.db.Exec(query, args...)
	if err != nil {
		log.Printf("Failed to update ticket %s: %v", ticketID, err) // Логирование ошибки
		return fmt.Errorf("database error: %v", err)
	}

	// Проверка, был ли тикет обновлён
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ticket not found or already deleted")
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


func (r *TicketRepository) GetAttachmentsByTicketID(ticketID string) ([]*models.TicketAttachment, error) {
	parsedID, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, fmt.Errorf("invalid ticket ID: %v", err)
	}
	attachments := []*models.TicketAttachment{}
	query := `
		SELECT attachment_id, ticket_id, filename, file_path, uploaded_by, uploaded_at
		FROM ticket_attachments
		WHERE ticket_id = $1
	`
	err = r.db.Select(&attachments, query, parsedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachments: %v", err)
	}
	return attachments, nil
}

func (r *TicketRepository) CreateAttachment(attachment *models.TicketAttachment) error {
	query := `
		INSERT INTO ticket_attachments (attachment_id, ticket_id, filename, file_path, uploaded_by, uploaded_at)
		VALUES (:attachment_id, :ticket_id, :filename, :file_path, :uploaded_by, :uploaded_at)
	`
	_, err := r.db.NamedExec(query, attachment)
	if err != nil {
		return fmt.Errorf("failed to create attachment: %v", err)
	}
	return nil
}

func (r *TicketRepository) GetAttachmentByID(attachmentID string) (*models.TicketAttachment, error) {
	parsedID, err := uuid.Parse(attachmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid attachment ID: %v", err)
	}
	attachment := &models.TicketAttachment{}
	query := `
		SELECT attachment_id, ticket_id, filename, file_path, uploaded_by, uploaded_at
		FROM ticket_attachments
		WHERE attachment_id = $1
	`
	err = r.db.Get(attachment, query, parsedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment: %v", err)
	}
	return attachment, nil
}

func (r *TicketRepository) DeleteAttachment(attachmentID string) error {
	parsedID, err := uuid.Parse(attachmentID)
	if err != nil {
		return fmt.Errorf("invalid attachment ID: %v", err)
	}
	query := "DELETE FROM ticket_attachments WHERE attachment_id = $1"
	result, err := r.db.Exec(query, parsedID)
	if err != nil {
		return fmt.Errorf("failed to delete attachment: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("attachment not found")
	}
	return nil
}

