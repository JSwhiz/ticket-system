package repository

import (
    "ticket-system/internal/models"
    "github.com/jmoiron/sqlx"
    "github.com/google/uuid"
)

type CommentRepository struct {
    db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
    return &CommentRepository{db}
}

func (r *CommentRepository) GetCommentsByTicketID(ticketID uuid.UUID) ([]*models.Comment, error) {
    comments := []*models.Comment{}
    query := "SELECT comment_id, ticket_id, author_id, content, created_at FROM ticket_comments WHERE ticket_id = $1 ORDER BY created_at DESC"
    err := r.db.Select(&comments, query, ticketID)
    if err != nil {
        return nil, err
    }
    return comments, nil
}


func (r *CommentRepository) CreateComment(comment *models.Comment) error {
    query := `
        INSERT INTO ticket_comments (comment_id, ticket_id, author_id, content, created_at)
        VALUES (:comment_id, :ticket_id, :author_id, :content, :created_at)
    `
    _, err := r.db.NamedExec(query, comment)
    return err
}