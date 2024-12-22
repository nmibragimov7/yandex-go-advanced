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
	f.Save(record.ShortURL, record.OriginalURL)

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}
	_, err = f.file.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write record to file: %w", err)
	}

	return nil
}

func (f *FileStorage) ReadRecord() (*models.ShortenRecord, error) {
	_, err := f.file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek file to the beginning: %w", err)
	}

	scanner := bufio.NewScanner(f.file)
	var records []*models.ShortenRecord

	for scanner.Scan() {
		var record models.ShortenRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file record: %w", err)
		}
		records = append(records, &record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner encountered an error: %w", err)
	}

	if len(records) == 0 {
		return &models.ShortenRecord{}, nil
	}

	return records[len(records)-1], nil
}

func (f *FileStorage) Close() error {
	if err := f.file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

func NewFileStorage(path string) (*FileStorage, error) {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return nil, fmt.Errorf("%s is a directory, not a file", path)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	storage := NewStorage()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var record models.ShortenRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file record: %w", err)
		}
		storage.Save(record.ShortURL, record.OriginalURL)
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
