package db

import (
	"errors"
	"fmt"
	"log"
	"yandex-go-advanced/internal/models"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

func (s *Storage) Get(key string) (interface{}, error) {
	var record models.ShortenRecord
	query := "SELECT short_url, original_url FROM shortener WHERE short_url = $1"
	err := s.DB.QueryRow(query, key).Scan(&record.ShortURL, &record.OriginalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get record from database: %w", err)
	}

	return &record, nil
}

func (s *Storage) Set(record interface{}) error {
	rec, ok := record.(*models.ShortenRecord)
	if !ok {
		return errors.New("failed to parse record interface")
	}

	query := "INSERT INTO shortener (short_url, original_url) VALUES ($1, $2)"
	_, err := s.DB.Exec(query, rec.ShortURL, rec.OriginalURL)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			var shortURL string
			query = "SELECT short_url FROM shortener WHERE original_url = $1"
			errs := s.DB.QueryRow(query, rec.OriginalURL).Scan(&shortURL)
			if errs != nil {
				return fmt.Errorf("failed to get record from database: %w", err)
			}

			return NewConflictError(
				shortURL,
				pgerrcode.UniqueViolation,
				err,
			)
		}

		return fmt.Errorf("failed to insert record into database: %w", err)
	}
	return nil
}

func (s *Storage) SetByTransaction(records []interface{}) error {
	rcs := make([]*models.ShortenRecord, 0, len(records))
	for _, record := range records {
		rec, ok := record.(*models.ShortenRecord)
		if !ok {
			return errors.New("failed to parse record interface")
		}

		rcs = append(rcs, rec)
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil {
			log.Printf("failed to rollback transaction: %s", err.Error())
		}
	}()

	for _, value := range rcs {
		query := "INSERT INTO shortener (short_url, original_url) VALUES ($1, $2)"
		_, err := tx.Exec(query, value.ShortURL, value.OriginalURL)
		if err != nil {
			return fmt.Errorf("failed to insert record into database: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	err := s.DB.Close()
	if err != nil {
		return fmt.Errorf("failed to close db storage: %w", err)
	}

	return nil
}
