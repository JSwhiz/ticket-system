package tickets

import "time"

type Ticket struct {
    ID           string     `db:"ticket_id" json:"ticket_id"`
    Title        string     `db:"title" json:"title"`
    Description  *string    `db:"description" json:"description,omitempty"`
    StatusID     int        `db:"status_id" json:"status_id"`
    PriorityID   int        `db:"priority_id" json:"priority_id"`
    CreatorID    string     `db:"creator_id" json:"creator_id"`
    AssigneeID   *string    `db:"assignee_id" json:"assignee_id,omitempty"`
    DepartmentID *int       `db:"department_id" json:"department_id,omitempty"`
    CreatedAt    time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	SearchVector *string    `db:"search_vector"  json:"-"`
}

type NewTicket struct {
    Title        string  `json:"title" binding:"required"`
    Description  *string `json:"description"`
    StatusID     int     `json:"status_id" binding:"required"`
    PriorityID   int     `json:"priority_id" binding:"required"`
    AssigneeID   *string `json:"assignee_id"`
    DepartmentID *int    `json:"department_id"`
}

type UpdateTicket struct {
    Title        *string `json:"title"`
    Description  *string `json:"description"`
    StatusID     *int    `json:"status_id"`
    PriorityID   *int    `json:"priority_id"`
    AssigneeID   *string `json:"assignee_id"`
    DepartmentID *int    `json:"department_id"`
}
