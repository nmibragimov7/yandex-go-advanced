package shortener

import "fmt"

type DuplicateError struct {
	Err      error
	ShortURL string
	Code     string
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
