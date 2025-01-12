package db

import "fmt"

type ConflictError struct {
	ShortURL string
	Code     string
	Err      error
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflicting %s with %s: %v", e.Code, e.ShortURL, e.Err)
}

func NewConflictError(shorten string, code string, err error) error {
	return &ConflictError{
		ShortURL: shorten,
		Code:     code,
		Err:      err,
	}
}
