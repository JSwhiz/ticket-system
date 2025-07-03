package models

import (
    "time"
    "github.com/google/uuid"
)

type Comment struct {
    ID        uuid.UUID  `db:"comment_id" json:"id"`
    TicketID  uuid.UUID  `db:"ticket_id" json:"ticket_id"`
    AuthorID  uuid.UUID  `db:"author_id" json:"author_id"`
    Content   string     `db:"content" json:"content"`
    CreatedAt time.Time  `db:"created_at" json:"created_at"`
}