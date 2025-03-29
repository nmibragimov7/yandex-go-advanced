package memory

import (
	"sync"
	"testing"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	t.Run("Memory set", func(t *testing.T) {
		stp := Storage{
			storage: map[string]string{},
			mtx:     &sync.Mutex{},
		}

		key := util.GetKey()
		record := &models.ShortenRecord{
			ShortURL:    key,
			OriginalURL: "https://yandex.ru/",
			UserID:      int64(1),
		}

		originalURL, err := stp.Set(record)
		assert.NoError(t, err)
		assert.NotEmpty(t, originalURL)
	})

	t.Run("Memory set with error", func(t *testing.T) {
		stp := Storage{
			storage: map[string]string{},
			mtx:     &sync.Mutex{},
		}

		key := util.GetKey()
		record := &models.ShortenRecord{
			ShortURL:    key,
			OriginalURL: "https://yandex.ru/",
			UserID:      int64(1),
		}

		originalURL, err := stp.Set(record)
		assert.NoError(t, err)
		assert.NotEmpty(t, originalURL)

		originalURL, err = stp.Set(record)
		assert.Error(t, err)
		assert.Empty(t, originalURL)
	})
}

func TestGet(t *testing.T) {
	t.Run("Memory get", func(t *testing.T) {
		stp := Storage{
			storage: map[string]string{},
			mtx:     &sync.Mutex{},
		}

		key := util.GetKey()
		record := &models.ShortenRecord{
			ShortURL:    key,
			OriginalURL: "https://kyc.kz/",
			UserID:      int64(1),
		}

		originalURL, err := stp.Set(record)
		assert.NoError(t, err)
		assert.NotEmpty(t, originalURL)

		result, err := stp.Get(key)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Memory get with error", func(t *testing.T) {
		stp := Storage{
			storage: map[string]string{},
			mtx:     &sync.Mutex{},
		}

		key := util.GetKey()
		record := &models.ShortenRecord{
			ShortURL:    key,
			OriginalURL: "https://kyc.kz/",
			UserID:      int64(1),
		}

		originalURL, err := stp.Set(record)
		assert.NoError(t, err)
		assert.NotEmpty(t, originalURL)

		result, err := stp.Get("12345678")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("Memory get all", func(t *testing.T) {
		stp := Storage{
			storage: map[string]string{},
			mtx:     &sync.Mutex{},
		}

		records, err := stp.GetAll(int64(1))
		assert.Empty(t, err)
		assert.Empty(t, records)
	})
}

func TestSetAll(t *testing.T) {
	t.Run("Memory set all", func(t *testing.T) {
		stp := Storage{
			storage: map[string]string{},
			mtx:     &sync.Mutex{},
		}

		records := []models.ShortenRecord{
			{
				ShortURL:    util.GetKey(),
				OriginalURL: "https://yandex.ru/",
				UserID:      int64(1),
			},
			{
				ShortURL:    util.GetKey(),
				OriginalURL: "https://kyc.kz/",
				UserID:      int64(1),
			},
		}
		values := make([]interface{}, 0, len(records))
		for _, value := range records {
			values = append(values, &models.ShortenRecord{
				OriginalURL: value.OriginalURL,
				ShortURL:    value.ShortURL,
				UserID:      value.UserID,
			})
		}

		err := stp.SetAll(values)
		assert.NoError(t, err)
	})
}

func TestClose(t *testing.T) {
	t.Run("Memory close", func(t *testing.T) {
		stp := Storage{
			storage: map[string]string{},
			mtx:     &sync.Mutex{},
		}

		err := stp.Close()
		assert.Empty(t, err)
	})
}

func TestPing(t *testing.T) {
	t.Run("Memory ping", func(t *testing.T) {
		stp := Storage{
			storage: map[string]string{},
			mtx:     &sync.Mutex{},
		}

		err := stp.Ping(nil)
		assert.Empty(t, err)
	})
}

func TestInit(t *testing.T) {
	t.Run("Memory initialized", func(t *testing.T) {
		cnf := Init()
		assert.NotNil(t, cnf)
	})
}
