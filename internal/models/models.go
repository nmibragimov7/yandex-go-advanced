package models

type ShortenRequestBody struct {
	URL string `json:"url,omitempty"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type ShortenResponseError struct {
	Message string `json:"message"`
}

type ShortenRecord struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
