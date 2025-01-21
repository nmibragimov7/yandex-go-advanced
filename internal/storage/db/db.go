package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository interface {
	Get(key string) (interface{}, error)
	Set(record interface{}) error
	SetAll(records []interface{}) error
	Close() error
	Ping(ctx context.Context) error
}

func Init(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database connection: %w", err)
	}

	err = bootstrap(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create table queries: %w", err)
	}

	return db, nil
}

func bootstrap(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil {
			log.Printf("failed to rollback transaction: %s", err.Error())
		}
	}()

	tables := []string{
		`CREATE TABLE IF NOT EXISTS shortener (
			id SERIAL PRIMARY KEY,
			short_url VARCHAR(10) NOT NULL,
			original_url VARCHAR(100) UNIQUE NOT NULL
		)`,
	}

	for _, query := range tables {
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("failed to create table query: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
