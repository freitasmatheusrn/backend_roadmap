package showtimes

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type RepoInterface interface {
	Create(ctx context.Context, showtime InputDTO) (int64, error)
	List(ctx context.Context, day time.Time) ([]OutputDTO, error)
}

type Repository struct {
	DB *pgx.Conn
}

func NewRepo(db *pgx.Conn) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) Create(ctx context.Context, showtime InputDTO) (int64, error) {
	var showtimeID int64
	err := r.DB.QueryRow(
		ctx,
		`INSERT INTO showtimes (movie_id, theater_id, start_time, end_time, base_price)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		showtime.MovieID,
		showtime.TheaterID,
		showtime.StartTime,
		showtime.EndTime,
		showtime.BasePrice,
	).Scan(&showtimeID)
	if err != nil {
		return 0, fmt.Errorf("erro ao inserir sessão: %w", err)
	}
	return showtimeID, nil
}

func (r *Repository) List(ctx context.Context, day time.Time) ([]OutputDTO, error) {
	var showtimeList []OutputDTO
	start := day.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	log.Println(start)
	log.Println(end)
	rows, err := r.DB.Query(ctx,
		`
		SELECT id, movie_id, theater_id, start_time, end_time, base_price, is_active
		FROM showtimes
		WHERE created_at >= $1
		AND created_at < $2
		ORDER BY movie_id
		`, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var o OutputDTO
		if err := rows.Scan(
			&o.ID, &o.MovieID, &o.TheaterID,
			&o.StartTime, &o.EndTime, &o.BasePrice,
			&o.IsActive,
		); err != nil {
			return nil, err
		}
		showtimeList = append(showtimeList, o)
	}
	return showtimeList, nil
}
