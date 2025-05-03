package users

import (
	"context"
	"errors"
	"fmt"

	"yandex-go-advanced/internal/models"

	"github.com/jmoiron/sqlx"
)

// Storage - struct that contains the necessary settings
type Storage struct {
	DB *sqlx.DB
}

// Get - func for return record
func (s *Storage) Get(_ string) (interface{}, error) {
	var record models.UserRecord
	return &record, nil
}

// GetAll - func for return records
func (s *Storage) GetAll(_ interface{}) ([]interface{}, error) {
	return nil, nil
}

// Set - func for saving record in database
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

// SetAll - func for saving records in database
func (s *Storage) SetAll(_ []interface{}) error {
	return nil
}

// Ping - func for ping database
func (s *Storage) Ping(ctx context.Context) error {
	if err := s.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close - func for close database
func (s *Storage) Close() error {
	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close db storage: %w", err)
	}

	return nil
}

// AddToChannel - func for add value in channel
func (s *Storage) AddToChannel(_ chan struct{}, _ ...chan interface{}) {}
