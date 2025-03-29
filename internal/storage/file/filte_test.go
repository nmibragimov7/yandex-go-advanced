package file

import (
	"os"
	"testing"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	t.Run("File set", func(t *testing.T) {
		_ = os.Remove("./storage.txt")

		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)

		key := util.GetKey()
		record := &models.ShortenRecord{
			ShortURL:    key,
			OriginalURL: "https://1cb.kz/",
			UserID:      int64(1),
		}

		originalURL, err := str.Set(record)
		assert.NoError(t, err)
		assert.NotEmpty(t, originalURL)
	})

	t.Run("File set with error", func(t *testing.T) {
		_ = os.Remove("./storage.txt")

		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)

		key := util.GetKey()
		record := &models.ShortenRecord{
			ShortURL:    key,
			OriginalURL: "https://fcbk.kz/",
			UserID:      int64(1),
		}

		originalURL, err := str.Set(record)
		assert.NoError(t, err)
		assert.NotEmpty(t, originalURL)

		originalURL, err = str.Set(record)
		assert.Error(t, err)
		assert.Empty(t, originalURL)
	})
}

func TestGet(t *testing.T) {
	t.Run("File get", func(t *testing.T) {
		_ = os.Remove("./storage.txt")

		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)

		key := util.GetKey()
		record := &models.ShortenRecord{
			ShortURL:    key,
			OriginalURL: "https://biometry.kz/",
			UserID:      int64(1),
		}

		originalURL, err := str.Set(record)
		assert.NoError(t, err)
		assert.NotEmpty(t, originalURL)

		result, err := str.Get(key)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Memory get with error", func(t *testing.T) {
		_ = os.Remove("./storage.txt")

		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)

		key := util.GetKey()
		record := &models.ShortenRecord{
			ShortURL:    key,
			OriginalURL: "https://dataset.kz/",
			UserID:      int64(1),
		}

		originalURL, err := str.Set(record)
		assert.NoError(t, err)
		assert.NotEmpty(t, originalURL)

		result, err := str.Get("12345678")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("File get all", func(t *testing.T) {
		_ = os.Remove("./storage.txt")

		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)

		records, err := str.GetAll(int64(1))
		assert.Empty(t, err)
		assert.Empty(t, records)
	})
}

func TestSetAll(t *testing.T) {
	t.Run("File set all", func(t *testing.T) {
		_ = os.Remove("./storage.txt")

		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)

		records := []models.ShortenRecord{
			{
				ShortURL:    util.GetKey(),
				OriginalURL: "https://1cb.kz/",
				UserID:      int64(1),
			},
			{
				ShortURL:    util.GetKey(),
				OriginalURL: "https://biometry.kz/",
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

		err = str.SetAll(values)
		assert.NoError(t, err)
	})
}

func TestClose(t *testing.T) {
	t.Run("File close", func(t *testing.T) {
		_ = os.Remove("./storage.txt")

		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)

		err = str.Close()
		assert.Empty(t, err)
	})
}

func TestPing(t *testing.T) {
	t.Run("File ping", func(t *testing.T) {
		_ = os.Remove("./storage.txt")

		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)

		err = str.Ping(nil)
		assert.Empty(t, err)
	})
}

func TestInit(t *testing.T) {
	t.Run("File initialized", func(t *testing.T) {
		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)
	})
}
