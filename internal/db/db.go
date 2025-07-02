package db

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func Connect(dsn string) (*sqlx.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	dbx := sqlx.NewDb(db, "postgres")
	// Тестовый запрос для проверки подключения
	var count int
	err = dbx.Get(&count, "SELECT COUNT(*) FROM users")
	if err != nil {
		log.Printf("Ошибка при проверке таблицы users: %v", err)
	} else {
		log.Printf("Успешное подключение! Количество записей в users: %d", count)
	}
	return dbx, nil
}
