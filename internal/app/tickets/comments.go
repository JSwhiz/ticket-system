package tickets

import (
    "time"

    "github.com/jmoiron/sqlx"
)

type Comment struct {
    CommentID string    `db:"comment_id" json:"comment_id"`
    TicketID  string    `db:"ticket_id"  json:"ticket_id"`
    AuthorID  string    `db:"author_id"  json:"author_id"`
    Content   string    `db:"content"    json:"content"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type NewComment struct {
    Content string `json:"content" binding:"required"`
}

type CommentRepo struct {
    db *sqlx.DB
}

func NewCommentRepo(db *sqlx.DB) *CommentRepo {
    return &CommentRepo{db: db}
}

func (r *CommentRepo) ListByTicket(ticketID string) ([]Comment, error) {
    const query = `
        SELECT comment_id, ticket_id, author_id, content, created_at
          FROM ticket_comments
         WHERE ticket_id = $1
         ORDER BY created_at ASC
    `
    comments := make([]Comment, 0)
    if err := r.db.Select(&comments, query, ticketID); err != nil {
        return nil, err
    }
    return comments, nil
}

func (r *CommentRepo) Create(ticketID, authorID, content string) (Comment, error) {
    const query = `
        INSERT INTO ticket_comments (ticket_id, author_id, content)
        VALUES ($1, $2, $3)
        RETURNING comment_id, ticket_id, author_id, content, created_at
    `
    var c Comment
    if err := r.db.Get(&c, query, ticketID, authorID, content); err != nil {
        return Comment{}, err
    }
    return c, nil
}
