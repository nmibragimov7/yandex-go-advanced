package models

// Response - struct for response
type Response struct {
	Message string `json:"message"`
}

// ShortenRequestBody - struct for shorten request body
type ShortenRequestBody struct {
	URL string `json:"url,omitempty"`
}

// ShortenResponseSuccess - struct for shorten response
type ShortenResponseSuccess struct {
	Result string `json:"result"`
}

// ShortenRecord - struct for shorten entity
type ShortenRecord struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      int64  `json:"user_id"`
	DeletedFlag bool   `json:"is_deleted"`
}

// UserRecord - struct for user entity
type UserRecord struct{}

// ShortenBatchRequest - struct for shorten batch request body
type ShortenBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ShortenBatchResponse - struct for shorten batch response
type ShortenBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ShortenBatchUpdateRequest - struct for shorten batch entity
type ShortenBatchUpdateRequest struct {
	ShortURL string `json:"short_url"`
	UserID   int64  `json:"user_id"`
}
