package showtimes

import "time"

type ShowTime struct {
	ID        int64
	MovieID   int64
	TheaterID int64
	StartTime time.Time
	EndTime   time.Time
	BasePrice float64
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type InputDTO struct {
	MovieID   int64     `json:"movie_id"`
	TheaterID int64     `json:"theater_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	BasePrice float64   `json:"base_price"`
}

type OutputDTO struct {
	ID        int64     `json:"id"`
	MovieID   int64     `json:"movie_id"`
	TheaterID int64     `json:"theater_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	BasePrice float64   `json:"base_price"`
	IsActive  bool      `json:"is_active"`
}
