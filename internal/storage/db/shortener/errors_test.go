package shortener

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := errors.New("unique constraint failed")
	dupErr := &DuplicateError{
		Err:      err,
		ShortURL: "https://short.url/abc123",
		Code:     "23505",
	}

	expectedErrMsg := "db error code 23505, exists https://short.url/abc123: unique constraint failed"
	assert.Equal(t, expectedErrMsg, dupErr.Error())
}

func TestNewDuplicateError(t *testing.T) {
	err := errors.New("duplicate entry")
	shortURL := "https://short.url/xyz789"
	code := "23505"

	dupErr := NewDuplicateError(shortURL, code, err)
	assert.IsType(t, &DuplicateError{}, dupErr)

	var de *DuplicateError
	ok := errors.As(dupErr, &de)
	assert.True(t, ok)
	assert.Equal(t, shortURL, de.ShortURL)
	assert.Equal(t, code, de.Code)
	assert.Equal(t, err, de.Err)
}
