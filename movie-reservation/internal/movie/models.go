package movie

import (
	"time"

	"github.com/freitasmatheusrn/movie-reservation/internal/showtimes"
)

type Movie struct {
	ID              int64         `json:"id"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	PosterURL       string        `json:"poster_url"`
	DurationMinutes time.Duration `json:"duration_minutes"`
	ReleaseDate     time.Time     `json:"release_date"`
	Rating          string        `json:"rating"`
	IsActive        bool          `json:"is_active"`
}

type InputDTO struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	PosterURL       string    `json:"poster_url"`
	DurationMinutes int       `json:"duration_minutes"`
	ReleaseDate     time.Time `json:"release_date"`
	Rating          string    `json:"rating"`
	IsActive        bool      `json:"is_active"`
	GenreIDs        []int64   `json:"genre_ids"`
}
type UpdateInputDTO struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	PosterURL       string    `json:"poster_url"`
	DurationMinutes int       `json:"duration_minutes"`
	ReleaseDate     time.Time `json:"release_date"`
	IsActive        bool      `json:"is_active"`
	GenreIDs        []int64   `json:"genre_ids"`
}
type OutputDTO struct {
	ID              int64                 `json:"id"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	PosterURL       string                `json:"poster_url"`
	DurationMinutes time.Duration         `json:"duration_minutes"`
	ReleaseDate     time.Time             `json:"release_date"`
	Rating          []string              `json:"rating"`
	IsActive        bool                  `json:"is_active"`
	Genres          []string              `json:"genres"`
	ShowTimes       []showtimes.OutputDTO `json:"show_times,omitempty"`
}
