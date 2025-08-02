package collections

type Url struct {
	ID          string `bson:"_id"`
	OriginalURL string `bson:"original_url"`
	AccessCount int    `bson:"access_count"`
}