package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type RepoInterface interface {
	Login(ctx context.Context, credentials LoginDTO) (*OutputDTO, error)
	Signup(ctx context.Context, user InputDTO) (string, error)
}

type Repository struct {
	DB *pgx.Conn
}

func NewRepo(db *pgx.Conn) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) Login(ctx context.Context, credentials LoginDTO) (*OutputDTO, error) {
	var user User
	err := r.DB.QueryRow(
		ctx,
		`SELECT id, name, email, password, role
		 FROM users
		 WHERE email = $1`,
		credentials.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("email não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	if err := ValidatePassword([]byte(credentials.Password), []byte(user.Password)); err != nil {
		return nil, fmt.Errorf("senha inválida")
	}
	return &OutputDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (r *Repository) Signup(ctx context.Context, user InputDTO) (string, error) {
	var userID string
	hashPassword, err := HashPassword(user.Password)
	if err != nil{
		return "", err
	}
	err = r.DB.QueryRow(
		ctx,
		`INSERT INTO users (name, email, password)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		user.Name,
		user.Email,
		hashPassword,
	).Scan(&userID)

	if err != nil {
		return "", fmt.Errorf("erro ao inserir usuário: %w", err)
	}

	return userID, nil
}
