package repository

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all pending migrations from the internal/migrations/ directory.
// Call from project root so the path resolves correctly.
func RunMigrations() error {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = DefaultDatabaseURL
	}
	m, err := migrate.New(
		"file://internal/migrations",
		url,
	)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

// RollbackMigrations runs the last migration down.
func RollbackMigrations() error {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = DefaultDatabaseURL
	}
	m, err := migrate.New(
		"file://internal/migrations",
		url,
	)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}
	defer m.Close()
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate steps -1: %w", err)
	}
	return nil
}
