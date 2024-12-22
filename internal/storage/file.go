package storage

import (
	"bufio"
	"encoding/json"
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
		return err
	}
	_, err = f.file.Write(append(data, '\n'))
	return err
}

func (f *FileStorage) ReadRecord() (*models.ShortenRecord, error) {
	_, err := f.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f.file)
	var records []*models.ShortenRecord

	for scanner.Scan() {
		var record models.ShortenRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	return records[len(records)-1], nil
}

func (f *FileStorage) Close() error {
	return f.file.Close()
}

func NewFileStorage(path string) (*FileStorage, error) {
	file, err := os.OpenFile(path+"/"+"storage.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	storage := NewStorage()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var record models.ShortenRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, err
		}
		storage.Save(record.ShortURL, record.OriginalURL)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &FileStorage{
		Storage: *storage,
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: scanner,
	}, nil
}
