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
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UUID        int    `json:"uuid"`
}
