package refs

import (
    "github.com/jmoiron/sqlx"
)

type Repository struct {
    db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
    return &Repository{db: db}
}

type Department struct {
    ID   int    `db:"id" json:"id"`
    Name string `db:"name" json:"name"`
}

type Status struct {
    ID    int    `db:"id" json:"id"`
    Label string `db:"label" json:"label"`
}

type Priority struct {
    ID    int    `db:"id" json:"id"`
    Label string `db:"label" json:"label"`
}

type Role struct {
    ID   string `db:"role_id" json:"id"`
    Name string `db:"name"    json:"name"`
}

type Permission struct {
    ID   string `db:"permission_id" json:"id"`
    Name string `db:"name"          json:"name"`
}

func (r *Repository) Departments() ([]Department, error) {
    var out []Department
    err := r.db.Select(&out, `SELECT id, name FROM mv_departments`)
    return out, err
}

func (r *Repository) Statuses() ([]Status, error) {
    var out []Status
    err := r.db.Select(&out, `SELECT id, label FROM mv_ticket_statuses`)
    return out, err
}

func (r *Repository) Priorities() ([]Priority, error) {
    var out []Priority
    err := r.db.Select(&out, `SELECT id, label FROM mv_ticket_priorities`)
    return out, err
}

func (r *Repository) Roles() ([]Role, error) {
    var out []Role
    err := r.db.Select(&out, `SELECT role_id, name FROM roles`)
    return out, err
}

func (r *Repository) Permissions() ([]Permission, error) {
    var out []Permission
    err := r.db.Select(&out, `SELECT permission_id, name FROM permissions`)
    return out, err
}
