package dtos

type UrlInputDTO struct {
	OriginalURL string `json:"original_url" bson:"original_url"`
}

type UrlOutputDTO struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	AccessCount int    `json:"access_count"`
}

type UpdateUrlInputDTO struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
}
