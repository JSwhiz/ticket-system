package tickets

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
    db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) List(filters map[string]interface{}) ([]Ticket, error) {
    where := []string{"deleted_at IS NULL"}
    args := []interface{}{}
    i := 1

    if v, ok := filters["status_id"]; ok {
        where = append(where, fmt.Sprintf("status_id = $%d", i))
        args = append(args, v)
        i++
    }
    if v, ok := filters["priority_id"]; ok {
        where = append(where, fmt.Sprintf("priority_id = $%d", i))
        args = append(args, v)
        i++
    }
    if v, ok := filters["assignee_id"]; ok {
        where = append(where, fmt.Sprintf("assignee_id = $%d", i))
        args = append(args, v)
        i++
    }
    if v, ok := filters["department_id"]; ok {
        where = append(where, fmt.Sprintf("department_id = $%d", i))
        args = append(args, v)
        i++
    }
    if v, ok := filters["created_from"]; ok {
        if t, ok2 := v.(time.Time); ok2 {
            where = append(where, fmt.Sprintf("created_at >= $%d", i))
            args = append(args, t)
            i++
        }
    }
    if v, ok := filters["created_to"]; ok {
        if t, ok2 := v.(time.Time); ok2 {
            where = append(where, fmt.Sprintf("created_at <= $%d", i))
            args = append(args, t)
            i++
        }
    }
    if v, ok := filters["search"]; ok {
        where = append(where, fmt.Sprintf("search_vector @@ plainto_tsquery($%d)", i))
        args = append(args, v)
        i++
    }

    query := `SELECT * FROM tickets WHERE ` + strings.Join(where, " AND ")
    var out []Ticket
    if err := r.db.Select(&out, query, args...); err != nil {
        return nil, err
    }
    return out, nil
}

func (r *Repository) Create(n NewTicket, creatorID string) (Ticket, error) {
    var t Ticket
    err := r.db.Get(&t, `
        INSERT INTO tickets (title, description, status_id, priority_id, creator_id, assignee_id, department_id)
        VALUES ($1,$2,$3,$4,$5,$6,$7)
        RETURNING *
    `, n.Title, n.Description, n.StatusID, n.PriorityID, creatorID, n.AssigneeID, n.DepartmentID)
    return t, err
}

func (r *Repository) GetByID(id string) (Ticket, error) {
    var t Ticket
    err := r.db.Get(&t, `SELECT * FROM tickets WHERE ticket_id=$1 AND deleted_at IS NULL`, id)
    return t, err
}

func (r *Repository) Update(id string, u UpdateTicket) (Ticket, error) {
    sets := []string{}
    args := []interface{}{}
    i := 1

    if u.Title != nil {
        sets = append(sets, fmt.Sprintf("title = $%d", i))
        args = append(args, *u.Title)
        i++
    }
    if u.Description != nil {
        sets = append(sets, fmt.Sprintf("description = $%d", i))
        args = append(args, *u.Description)
        i++
    }
    if u.StatusID != nil {
        sets = append(sets, fmt.Sprintf("status_id = $%d", i))
        args = append(args, *u.StatusID)
        i++
    }
    if u.PriorityID != nil {
        sets = append(sets, fmt.Sprintf("priority_id = $%d", i))
        args = append(args, *u.PriorityID)
        i++
    }
    if u.AssigneeID != nil {
        sets = append(sets, fmt.Sprintf("assignee_id = $%d", i))
        args = append(args, *u.AssigneeID)
        i++
    }
    if u.DepartmentID != nil {
        sets = append(sets, fmt.Sprintf("department_id = $%d", i))
        args = append(args, *u.DepartmentID)
        i++
    }

    if len(sets) == 0 {
        return r.GetByID(id)
    }

    args = append(args, id)
    query := fmt.Sprintf(
        "UPDATE tickets SET %s WHERE ticket_id = $%d AND deleted_at IS NULL RETURNING *",
        strings.Join(sets, ", "),
        i,
    )

    var t Ticket
    err := r.db.Get(&t, query, args...)
    return t, err
}

func (r *Repository) Delete(id string) error {
    _, err := r.db.Exec(
        "UPDATE tickets SET deleted_at = now() WHERE ticket_id = $1 AND deleted_at IS NULL",
        id,
    )
    return err
}

