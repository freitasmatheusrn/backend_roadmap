package genre

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type RepoInterface interface {
	Create(ctx context.Context, genre InputDTO) (OutputDTO, error)
	List(ctx context.Context) ([]OutputDTO, error)
}

type Repository struct {
	DB *pgx.Conn
}

func NewRepo(db *pgx.Conn) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) Create(ctx context.Context, genre InputDTO) (OutputDTO, error) {
	var g OutputDTO
	err := r.DB.QueryRow(
		ctx,
		`INSERT INTO genres (name)
		 VALUES ($1)
		 RETURNING id, name`,
		genre.Name,
	).Scan(&g.ID, &g.Name)
	if err != nil {
		return OutputDTO{}, fmt.Errorf("erro ao inserir gênero: %w", err)
	}

	return g, nil
}

func (r *Repository) List(ctx context.Context) ([]OutputDTO, error) {
	var genres []OutputDTO

	rows, err := r.DB.Query(ctx, `SELECT id, name FROM genres`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g OutputDTO
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}