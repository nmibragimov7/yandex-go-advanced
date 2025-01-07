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
	db     *sqlx.DB
	Config *config.Config
	Sugar  *zap.SugaredLogger
}

func (p *DatabaseProvider) Init() error {
	db, err := sqlx.Open("postgres", *p.Config.DataBase)
	if err != nil {
		p.Sugar.Errorw(
			"Failed to open database connection",
			"error", err.Error(),
		)
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database connection: %w", err)
	}

	p.db = db
	p.Sugar.Infow("Database connection initialized successfully")

	return nil
}
func (p *DatabaseProvider) Get() *sqlx.DB {
	return p.db
}
