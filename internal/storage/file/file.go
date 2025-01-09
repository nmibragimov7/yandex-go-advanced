package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"yandex-go-advanced/internal/models"
)

type Storage struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func (s *Storage) Get(key string) (string, error) {
	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		var record models.ShortenRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return "", fmt.Errorf("failed to unmarshal file record: %w", err)
		}

		if record.ShortURL == key {
			return record.OriginalURL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to scan: %w", err)
	}

	return "", errors.New("failed to find record")
}

func (s *Storage) Set(record *models.ShortenRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}
	_, err = s.file.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write record to file: %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	if err := s.file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

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
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: scanner,
	}, nil
}
