package file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"yandex-go-advanced/internal/models"
)

type Storage struct {
	file   *os.File
	writer *bufio.Writer
}

func (s *Storage) Get(key string) (interface{}, error) {
	if _, err := s.file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		var record models.ShortenRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file record: %w", err)
		}

		if record.ShortURL == key {
			return &record, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return nil, errors.New("failed to find record")
}

func (s *Storage) GetAll(_ interface{}) ([]interface{}, error) {
	return nil, nil
}

func (s *Storage) Set(record interface{}) (interface{}, error) {
	rec, ok := record.(*models.ShortenRecord)
	if !ok {
		return nil, errors.New("failed to parse record interface")
	}

	data, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record: %w", err)
	}
	_, err = s.file.Write(append(data, '\n'))
	if err != nil {
		return nil, fmt.Errorf("failed to write record to file: %w", err)
	}

	return nil, nil
}

func (s *Storage) SetAll(records []interface{}) error {
	rcs := make([]*models.ShortenRecord, 0, len(records))
	for _, record := range records {
		rec, ok := record.(*models.ShortenRecord)
		if !ok {
			return errors.New("failed to parse record interface")
		}

		rcs = append(rcs, rec)
	}

	for _, value := range rcs {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal record: %w", err)
		}
		_, err = s.file.Write(append(data, '\n'))
		if err != nil {
			return fmt.Errorf("failed to write record to file: %w", err)
		}
	}

	return nil
}

func (s *Storage) Close() error {
	if err := s.file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

func (s *Storage) Ping(_ context.Context) error { return nil }

func Init(path string) (*Storage, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var record models.ShortenRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file record: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner encountered an error: %w", err)
	}

	return &Storage{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}
