package users

import (
	"context"
	"errors"
	"fmt"
	"yandex-go-advanced/internal/models"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	DB *sqlx.DB
}

func (s *Storage) Get(_ string) (interface{}, error) {
	var record models.UserRecord
	return &record, nil
}

func (s *Storage) GetAll(_ interface{}) ([]interface{}, error) {
	return nil, nil
}

func (s *Storage) Set(record interface{}) (interface{}, error) {
	_, ok := record.(*models.UserRecord)
	if !ok {
		return nil, errors.New("failed to parse record interface")
	}

	var id int64
	query := "INSERT INTO users DEFAULT VALUES RETURNING id"
	err := s.DB.QueryRowx(query).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert record into database: %w", err)
	}
	return id, nil
}

func (s *Storage) SetAll(_ []interface{}) error {
	return nil
}

func (s *Storage) Ping(ctx context.Context) error {
	if err := s.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close db storage: %w", err)
	}

	return nil
}

func (s *Storage) UpdateAll(_ chan struct{}, _ ...chan interface{}) { return }
