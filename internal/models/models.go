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
}
