package file

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/util"

	"github.com/stretchr/testify/assert"
)

func createTempFileWithContent(t *testing.T, content string) *os.File {
	t.Helper()

	f, err := os.CreateTemp("", "set-test")
	assert.NoError(t, err)

	_, err = f.WriteString(content)
	assert.NoError(t, err)

	_, err = f.Seek(0, 0)
	assert.NoError(t, err)

	t.Cleanup(func() {
		err = os.Remove(f.Name())
		assert.NoError(t, err)
	})
	return f
}

type brokenWriteFile struct {
	*os.File
}

func (b *brokenWriteFile) Write(_ []byte) (n int, err error) {
	return 0, errors.New("failed to write record to file")
}

func (b *brokenWriteFile) Close() error {
	return b.File.Close()
}

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

	t.Run("File set with already exists error", func(t *testing.T) {
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

	t.Run("File set with parse record error", func(t *testing.T) {
		f := createTempFileWithContent(t, "")
		s := &Storage{file: f}

		_, err := s.Set("not a struct")
		assert.EqualError(t, err, "failed to parse record interface")
	})

	t.Run("File set with seek error", func(t *testing.T) {
		f := createTempFileWithContent(t, "")
		err := f.Close()
		assert.NoError(t, err)

		s := &Storage{file: f}

		rec := &models.ShortenRecord{
			ShortURL:    "abc",
			OriginalURL: "https://test.com",
		}

		_, err = s.Set(rec)
		assert.ErrorContains(t, err, "seek")
	})

	t.Run("File set with unmarshal error", func(t *testing.T) {
		f := createTempFileWithContent(t, "bad-json\n")

		s := &Storage{file: f}

		rec := &models.ShortenRecord{
			ShortURL:    "abc",
			OriginalURL: "https://test.com",
		}

		_, err := s.Set(rec)
		assert.ErrorContains(t, err, "unmarshal")
	})

	t.Run("File set with unmarshal error", func(t *testing.T) {
		f := createTempFileWithContent(t, "bad-json\n")

		s := &Storage{file: f}

		rec := &models.ShortenRecord{
			ShortURL:    "abc",
			OriginalURL: "https://test.com",
		}

		_, err := s.Set(rec)
		assert.ErrorContains(t, err, "unmarshal")
	})

	t.Run("File set with write error", func(t *testing.T) {
		f, err := os.CreateTemp("", "set-test")
		assert.NoError(t, err)

		s := &Storage{
			file: &brokenWriteFile{File: f},
		}

		rec := &models.ShortenRecord{
			ShortURL:    "abc",
			OriginalURL: "https://writeerror.com",
		}

		_, err = s.Set(rec)
		assert.ErrorContains(t, err, "failed to write record to file")
	})
}

type brokenSeekFile struct {
	*os.File
}

func (b *brokenSeekFile) Seek(_ int64, _ int) (int64, error) {
	return 0, errors.New("seek error")
}

func (b *brokenSeekFile) Close() error {
	return b.File.Close()
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

	t.Run("File get with no find error", func(t *testing.T) {
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

	t.Run("File get with seek error", func(t *testing.T) {
		f, err := os.CreateTemp("", "broken")
		assert.NoError(t, err)
		defer func(name string) {
			err = os.Remove(name)
			assert.NoError(t, err)
		}(f.Name())

		st := &Storage{
			file: &brokenSeekFile{f},
		}

		_, err = st.Get("key")
		assert.ErrorContains(t, err, "seek")
	})

	t.Run("File get with unmarshal error", func(t *testing.T) {
		f, err := os.CreateTemp("", "bad_json")
		assert.NoError(t, err)
		defer func(name string) {
			err = os.Remove(name)
			assert.NoError(t, err)
		}(f.Name())

		_, err = f.WriteString("not-a-json\n")
		assert.NoError(t, err)
		_, err = f.Seek(0, 0)
		assert.NoError(t, err)

		st := &Storage{
			file: f,
		}

		_, err = st.Get("key")
		assert.ErrorContains(t, err, "unmarshal")
	})

	t.Run("File get with scanner error", func(t *testing.T) {
		f, err := os.CreateTemp("", "long_line")
		assert.NoError(t, err)
		defer func(name string) {
			err = os.Remove(name)
			assert.NoError(t, err)
		}(f.Name())

		_, err = f.WriteString(strings.Repeat("a", 1024*1024))
		assert.NoError(t, err)
		_, err = f.Seek(0, 0)
		assert.NoError(t, err)

		st := &Storage{
			file: f,
		}

		_, err = st.Get("key")
		assert.ErrorContains(t, err, "scan")
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

type brokenCloseFile struct {
	*os.File
}

func (b *brokenCloseFile) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (b *brokenCloseFile) Close() error {
	return errors.New("close error")
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

	t.Run("File close with error", func(t *testing.T) {
		f, err := os.CreateTemp("", "broken")
		assert.NoError(t, err)
		defer func(name string) {
			err = os.Remove(name)
			assert.NoError(t, err)
		}(f.Name())

		st := &Storage{
			file: &brokenCloseFile{f},
		}

		err = st.Close()
		assert.ErrorContains(t, err, "close")
	})
}

func TestPing(t *testing.T) {
	t.Run("File ping", func(t *testing.T) {
		_ = os.Remove("./storage.txt")

		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)

		ctx := context.Context(context.Background())
		err = str.Ping(ctx)
		assert.Empty(t, err)
	})
}

func TestInit(t *testing.T) {
	t.Run("File initialized", func(t *testing.T) {
		str, err := Init("./storage.txt")
		assert.NoError(t, err)
		assert.NotNil(t, str)
	})

	t.Run("File initialized with open file error", func(t *testing.T) {
		osOpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
			return nil, errors.New("open failed")
		}

		_, err := Init("./storage.txt")
		assert.ErrorContains(t, err, "open failed")
	})

	t.Run("File initialized with unmarshal error", func(t *testing.T) {
		f, err := os.CreateTemp("", "badjson")
		assert.NoError(t, err)
		defer func(name string) {
			err = os.Remove(name)
			assert.NoError(t, err)
		}(f.Name())

		_, err = f.WriteString("this-is-not-json\n")
		assert.NoError(t, err)
		_, err = f.Seek(0, 0)
		assert.NoError(t, err)

		osOpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
			return os.Open(f.Name())
		}

		_, err = Init("./storage.txt")
		assert.ErrorContains(t, err, "unmarshal")
	})

	t.Run("File initialized with scanner error", func(t *testing.T) {
		f, err := os.CreateTemp("", "long_line")
		assert.NoError(t, err)
		defer func(name string) {
			err = os.Remove(name)
			assert.NoError(t, err)
		}(f.Name())

		_, err = f.WriteString(strings.Repeat("a", 1024*1024)) // > bufio.Scanner limit
		assert.NoError(t, err)
		_, err = f.Seek(0, 0)
		assert.NoError(t, err)

		osOpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
			return os.Open(f.Name())
		}

		_, err = Init("./storage.txt")
		assert.ErrorContains(t, err, "scanner")
	})
}
