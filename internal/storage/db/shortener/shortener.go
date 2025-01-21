package shortener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"yandex-go-advanced/internal/models"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (s *Storage) GetAll(key interface{}) ([]interface{}, error) {
	var records []interface{}

	rows, err := s.DB.Query("SELECT short_url, original_url FROM shortener WHERE user_id = $1", key)
	if err != nil {
		return nil, fmt.Errorf("failed to query records: %w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %s", err.Error())
		}
	}()

	for rows.Next() {
		var record models.ShortenRecord

		err := rows.Scan(&record.ShortURL, &record.OriginalURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan records: %w", err)
	}

	return records, nil
}

func (s *Storage) Set(record interface{}) (interface{}, error) {
	rec, ok := record.(*models.ShortenRecord)
	if !ok {
		return nil, errors.New("failed to parse record interface")
	}

	query := "INSERT INTO shortener (short_url, original_url, user_id) VALUES ($1, $2, $3)"
	result, err := s.DB.Exec(query, rec.ShortURL, rec.OriginalURL, rec.UserID)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			var shortURL string
			query = "SELECT short_url FROM shortener WHERE original_url = $1"
			errs := s.DB.QueryRow(query, rec.OriginalURL).Scan(&shortURL)
			if errs != nil {
				return nil, fmt.Errorf("failed to get record from database: %w", err)
			}

			return nil, fmt.Errorf("shortener already exists: %w", NewDuplicateError(
				shortURL,
				pgerrcode.UniqueViolation,
				err,
			))
		}

		return nil, fmt.Errorf("failed to insert record into database: %w", err)
	}
	return result, nil
}

func (s *Storage) SetAll(records []interface{}) error {
	maxRetries := 5

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := s.SaveBatches(records)
		if err == nil {
			return nil
		}

		if !(strings.Contains(err.Error(), "deadlock") || strings.Contains(err.Error(), "timeout")) {
			return fmt.Errorf("failed to insert records into database: %w", err)
		}

		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return errors.New("failed to attempt retries")
}

func (s *Storage) SaveBatches(records []interface{}) error {
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

	query := `INSERT INTO shortener (short_url, original_url, user_id) VALUES `
	params := make([]interface{}, 0, len(rcs)*3)
	for i, record := range rcs {
		query += fmt.Sprintf("($%d,$%d,$%d),", i*2+1, i*2+2, i*2+3)
		params = append(params, record.ShortURL, record.OriginalURL, record.UserID)
	}

	query = query[:len(query)-1]

	_, err = tx.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("failed to insert records into database: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

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
