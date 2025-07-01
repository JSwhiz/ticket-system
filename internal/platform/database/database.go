package database

import (
    "fmt"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func Open(dbURL, migrationsDir string) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        return nil, fmt.Errorf("connect to db: %w", err)
    }

    driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
    if err != nil {
        return nil, fmt.Errorf("create migrate driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://"+migrationsDir,
        "postgres", driver,
    )
    if err != nil {
        return nil, fmt.Errorf("init migrate: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return nil, fmt.Errorf("apply migrations: %w", err)
    }

    return db, nil
}
