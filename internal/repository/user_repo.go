package repository

import (
	"database/sql"
	"fmt"
	"log"
	"ticket-system/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	log.Printf("Попытка найти пользователя с username: %s", username)
	err := r.db.Get(user, "SELECT user_id, username, password_hash, email, role_id, department_id, created_at, last_login, deleted_at FROM users WHERE LOWER(username) = LOWER($1) AND deleted_at IS NULL", username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Пользователь %s не найден в базе данных", username)
			return nil, fmt.Errorf("пользователь не найден")
		}
		log.Printf("Ошибка при запросе к базе данных: %v", err)
		return nil, err
	}
	log.Printf("Найден пользователь: %+v", user)
	return user, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (user_id, username, password_hash, email, role_id, department_id, created_at, last_login, deleted_at)
		VALUES (:user_id, :username, :password_hash, :email, :role_id, :department_id, :created_at, :last_login, :deleted_at)
	`
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		return err
	}
	return nil
}
