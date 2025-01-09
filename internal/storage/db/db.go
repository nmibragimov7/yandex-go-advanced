package db

import (
	"context"
	"fmt"
	"time"
	"yandex-go-advanced/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DatabaseProvider struct {
	Config *config.Config
	Sugar  *zap.SugaredLogger
}

func (p *DatabaseProvider) Init() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", *p.Config.DataBase)
	if err != nil {
		p.Sugar.Errorw(
			"Failed to open database connection",
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database connection: %w", err)
	}

	p.Sugar.Infow("Database connection initialized successfully")

	return db, nil
}

func (p *DatabaseProvider) CreateTables(db *sqlx.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS shortener (
			id SERIAL PRIMARY KEY,
			short_url VARCHAR(10) NOT NULL,
			original_url VARCHAR(100) NOT NULL
		)`,
	}

	for _, query := range tables {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table query: %w", err)
		}
	}

	p.Sugar.Infow("All tables created successfully")
	return nil
}
