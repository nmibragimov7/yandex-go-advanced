package statistics

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
	return nil, nil
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

// GetStat - func for return stats
func (s *Storage) GetStat() (interface{}, error) {
	var record models.StatResponse

	query := "SELECT (SELECT COUNT(*) FROM shortener), (SELECT COUNT(*) FROM users)"
	err := s.DB.QueryRow(query).Scan(&record.Urls, &record.Users)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats: %w", err)
	}

	return &record, nil
}
