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
