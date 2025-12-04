package movie

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/freitasmatheusrn/movie-reservation/internal/showtimes"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RepoInterface interface {
	Create(ctx context.Context, movie InputDTO) (int64, error)
	Update(ctx context.Context, movie UpdateInputDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]OutputDTO, error)
	ListByGenre(ctx context.Context, genres []string) ([]OutputDTO, error)
	ListByName(ctx context.Context, name string) ([]OutputDTO, error)
	Show(ctx context.Context, id int64) (OutputDTO, error)
}

type Repository struct {
	DB *pgx.Conn
}

func NewRepo(db *pgx.Conn) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) Create(ctx context.Context, movie InputDTO) (int64, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	var movieID int64
	err = tx.QueryRow(
		ctx,
		`INSERT INTO movies 
            (title, description, poster_url, duration_minutes, release_date, rating)
         VALUES ($1, $2, $3, $4, $5, $6)
         RETURNING id`,
		movie.Title,
		movie.Description,
		movie.PosterURL,
		movie.DurationMinutes,
		movie.ReleaseDate,
		movie.Rating,
	).Scan(&movieID)
	if err != nil {
		return 0, err
	}
	if len(movie.GenreIDs) > 0 {
		batch := &pgx.Batch{}
		for _, genreID := range movie.GenreIDs {
			batch.Queue(
				`INSERT INTO movie_genres (movie_id, genre_id) VALUES ($1, $2)`,
				movieID, genreID,
			)
		}

		br := tx.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return movieID, nil
}

func (r *Repository) Update(ctx context.Context, input UpdateInputDTO) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
        UPDATE movies
        SET 
            title = $1,
            description = $2,
            poster_url = $3,
            duration_minutes = $4,
            release_date = $5,
            is_active = $6,
            updated_at = NOW()
        WHERE id = $7
    `,
		input.Title,
		input.Description,
		input.PosterURL,
		input.DurationMinutes,
		input.ReleaseDate,
		input.IsActive,
		input.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar filme: %w", err)
	}

	_, err = tx.Exec(ctx, `
        DELETE FROM movie_genres WHERE movie_id = $1
    `, input.ID)
	if err != nil {
		return fmt.Errorf("falha ao limpar gêneros: %w", err)
	}

	for _, genreID := range input.GenreIDs {
		_, err = tx.Exec(ctx, `
            INSERT INTO movie_genres (movie_id, genre_id)
            VALUES ($1, $2)
        `, input.ID, genreID)
		if err != nil {
			return fmt.Errorf("falha ao inserir gêneros %d: %w", genreID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) List(ctx context.Context) ([]OutputDTO, error) {
	var movies []OutputDTO

	rows, err := r.DB.Query(ctx, `
        SELECT 
            m.id, 
            m.title, 
            m.description,
            m.poster_url,
            m.duration_minutes,
            m.release_date,
            m.rating,
            m.is_active,
            COALESCE(array_agg(g.name ORDER BY g.name) 
                FILTER (WHERE g.id IS NOT NULL), '{}') AS genres
        FROM movies m
        LEFT JOIN movie_genres mg ON mg.movie_id = m.id
        LEFT JOIN genres g ON g.id = mg.genre_id
        GROUP BY m.id
        ORDER BY m.id;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m OutputDTO
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Description,
			&m.PosterURL,
			&m.DurationMinutes,
			&m.ReleaseDate,
			&m.Rating,
			&m.IsActive,
			&m.Genres,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *Repository) ListByGenre(ctx context.Context, genres []string) ([]OutputDTO, error) {
	var movies []OutputDTO
	rows, err := r.DB.Query(ctx, `
        SELECT 
            m.id,
            m.title,
            m.description,
            m.poster_url,
            m.duration_minutes,
            m.release_date,
            m.rating,
            m.is_active,
            COALESCE(array_agg(DISTINCT g2.name ORDER BY g2.name)
                     FILTER (WHERE g2.id IS NOT NULL), '{}') AS genres
        FROM movies m
        JOIN movie_genres mg ON mg.movie_id = m.id
        JOIN genres g ON g.id = mg.genre_id
        LEFT JOIN movie_genres mg2 ON mg2.movie_id = m.id
        LEFT JOIN genres g2 ON g2.id = mg2.genre_id
        WHERE g.name ILIKE ANY($1)
        GROUP BY m.id
        ORDER BY m.id;
    `, genres)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m OutputDTO
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Description,
			&m.PosterURL,
			&m.DurationMinutes,
			&m.ReleaseDate,
			&m.Rating,
			&m.IsActive,
			&m.Genres,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *Repository) ListByName(ctx context.Context, name string) ([]OutputDTO, error) {
	var movies []OutputDTO
	rows, err := r.DB.Query(ctx, `
        SELECT 
            m.id, 
            m.title, 
            m.description,
            m.poster_url,
            m.duration_minutes,
            m.release_date,
            m.rating,
            m.is_active,
            COALESCE(array_agg(g.name ORDER BY g.name) 
                FILTER (WHERE g.id IS NOT NULL), '{}') AS genres
        FROM movies m
        LEFT JOIN movie_genres mg ON mg.movie_id = m.id
        LEFT JOIN genres g ON g.id = mg.genre_id
        WHERE m.title ILIKE $1
        GROUP BY m.id
        ORDER BY m.id DESC;
    `, "%"+strings.TrimSpace(name)+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m OutputDTO
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Description,
			&m.PosterURL,
			&m.DurationMinutes,
			&m.ReleaseDate,
			&m.Rating,
			&m.IsActive,
			&m.Genres,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *Repository) Show(ctx context.Context, id int64) (OutputDTO, error) {
	var m OutputDTO
	err := r.DB.QueryRow(ctx,
		`
		SELECT id, title, description, poster_url, duration_minutes,
			release_date, rating, is_active
		FROM movies
			WHERE id = $1
		`, id).Scan(
		&m.ID,
		&m.Title,
		&m.Description,
		&m.PosterURL,
		&m.DurationMinutes,
		&m.ReleaseDate,
		&m.Rating,
		&m.IsActive,
	)
	if err != nil {
		return OutputDTO{}, err
	}

	rows, err := r.DB.Query(ctx,
		`
		SELECT id, start_time, end_time, base_price
		FROM showtimes
		WHERE movie_id = $1
		ORDER BY start_time ASC
		`, id)
	if err != nil {
		return OutputDTO{}, err
	}
	defer rows.Close()

	var showtimesList []showtimes.OutputDTO
	for rows.Next() {
		var s showtimes.OutputDTO
		if err := rows.Scan(&s.ID, &s.StartTime, &s.EndTime, &s.BasePrice); err != nil {
			return OutputDTO{}, err
		}
		showtimesList = append(showtimesList, s)
	}
	m.ShowTimes = showtimesList
	return m, nil
}


func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	cmdTag, err := r.DB.Exec(ctx,
		`DELETE FROM reservations WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("reservation %s not found", id)
	}

	return nil
}