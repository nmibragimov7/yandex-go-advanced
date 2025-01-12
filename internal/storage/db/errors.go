package db

import "fmt"

type DuplicateError struct {
	ShortURL string
	Code     string
	Err      error
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("db error code %s, exists %s: %v", e.Code, e.ShortURL, e.Err)
}

func NewDuplicateError(shorten string, code string, err error) error {
	return &DuplicateError{
		ShortURL: shorten,
		Code:     code,
		Err:      err,
	}
}
