package models

import (
	"time"
	"github.com/google/uuid"
)

type Ticket struct {
	TicketID    uuid.UUID  `db:"ticket_id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"` // NULLable
	StatusID    int16      `db:"status_id"`   // smallint
	PriorityID  int16      `db:"priority_id"` // smallint
	CreatorID   uuid.UUID  `db:"creator_id"`
	AssigneeID  *uuid.UUID `db:"assignee_id"` // NULLable
	DepartmentID *int16    `db:"department_id"` // NULLable smallint
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"` // NULLable
	DeletedAt   *time.Time `db:"deleted_at"` // NULLable
	SearchVector *string   `db:"search_vector"` // NULLable
}