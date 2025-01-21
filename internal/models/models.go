package models

type Response struct {
	Message string `json:"message"`
}

type ShortenRequestBody struct {
	URL string `json:"url,omitempty"`
}

type ShortenResponseSuccess struct {
	Result string `json:"result"`
}

type ShortenRecord struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      int64  `json:"user_id"`
}

type UserRecord struct{}

type ShortenBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
