package postgres

import (
	"errors"
	"fmt"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/storage/db"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	DB *sqlx.DB
}

func (s *Storage) Get(key string) (interface{}, error) {
	var record models.ShortenRecord
	query := "SELECT short_url, original_url FROM shortener WHERE short_url = $1"
	err := s.DB.QueryRow(query, key).Scan(&record.ShortURL, &record.OriginalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get record from database: %w", err)
	}

	return &record, nil
}

func (s *db.Storage) Set(record interface{}) error {
	rec, ok := record.(*models.ShortenRecord)
	if !ok {
		return errors.New("failed to parse record interface")
	}

	query := "INSERT INTO shortener (short_url, original_url) VALUES ($1, $2)"
	_, err := s.DB.Exec(query, rec.ShortURL, rec.OriginalURL)
	if err != nil {
		return fmt.Errorf("failed to insert record into database: %w", err)
	}
	return nil
}

func (s *db.Storage) Close() error {
	err := s.DB.Close()
	if err != nil {
		return fmt.Errorf("failed to close db storage: %w", err)
	}

	return nil
}
