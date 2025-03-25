package shortener

import "fmt"

// DuplicateError - struct that contains the necessary settings
type DuplicateError struct {
	Err      error
	ShortURL string
	Code     string
}

// Error - func for return error string
func (e *DuplicateError) Error() string {
	return fmt.Sprintf("db error code %s, exists %s: %v", e.Code, e.ShortURL, e.Err)
}

// NewDuplicateError - func for return error instance
func NewDuplicateError(shorten string, code string, err error) error {
	return &DuplicateError{
		ShortURL: shorten,
		Code:     code,
		Err:      err,
	}
}
