package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"yandex-go-advanced/internal/models"
)

type FileStorage struct {
	Storage
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func (f *FileStorage) WriteRecord(record *models.ShortenRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}
	_, err = f.file.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write record to file: %w", err)
	}
	f.save(record.ShortURL, record.OriginalURL)

	return nil
}

func (f *FileStorage) Close() error {
	if err := f.file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

func InitFileStorage(path string) (*FileStorage, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	storage := newStorage()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var record models.ShortenRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file record: %w", err)
		}
		storage.save(record.ShortURL, record.OriginalURL)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner encountered an error: %w", err)
	}

	return &FileStorage{
		Storage: *storage,
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: scanner,
	}, nil
}
