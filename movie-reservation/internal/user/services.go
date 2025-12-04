package user

import (
	"context"
	"fmt"

	"github.com/freitasmatheusrn/movie-reservation/pkg/errsx"
)

type ServiceInterface interface {
	Login(ctx context.Context, email, password string) (*OutputDTO, error)
	Signup(ctx context.Context, name, email, password string) (int64, error)
}

type Service struct {
	repo RepoInterface
}

func NewService(userRepo RepoInterface) *Service {
	return &Service{
		repo: userRepo,
	}
}

type SignupResp struct {
	UserID string `json:"user_id"`
}

func (s *Service) Login(ctx context.Context, credentials LoginDTO) (*OutputDTO, error) {
	output, err := s.repo.Login(ctx, credentials)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (s *Service) Signup(ctx context.Context, input InputDTO) (*SignupResp, error) {
	var errs errsx.Map
	var err error
	err = NameValid(input.Name)
	if err != nil {
		errs.Set("name", err)
	}
	err = EmailValid(input.Email)
	if err != nil {
		errs.Set("email", err)
	}
	err = PasswordValid(input.Password)
	if err != nil {
		errs.Set("password", err)
	}
	if errs != nil {
		return nil, fmt.Errorf("bad Request: %w", errs)
	}
	id, err := s.repo.Signup(ctx, input)
	if err != nil {
		return nil, err
	}
	return &SignupResp{UserID: id}, err
}
