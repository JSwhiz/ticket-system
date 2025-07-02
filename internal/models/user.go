package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID  `db:"user_id"`
	Username     string     `db:"username"`
	PasswordHash string     `db:"password_hash"`
	Email        string     `db:"email"`
	RoleID       uuid.UUID  `db:"role_id"`
	DepartmentID int        `db:"department_id"`
	CreatedAt    time.Time  `db:"created_at"`
	LastLogin    *time.Time `db:"last_login"`
	DeletedAt    *time.Time `db:"deleted_at"`
}
