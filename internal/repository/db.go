package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DefaultDatabaseURL is the fallback connection string when DATABASE_URL is not set.
const DefaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/app_db_name?sslmode=disable"

// NewDBPool opens a PostgreSQL connection pool from DATABASE_URL or default.
func NewDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = DefaultDatabaseURL
	}
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db ping: %w", err)
	}
	return pool, nil
}
