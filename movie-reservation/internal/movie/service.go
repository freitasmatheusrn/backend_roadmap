package movie

import (
	"context"
	"fmt"

	"github.com/freitasmatheusrn/movie-reservation/internal/genre"
	"github.com/freitasmatheusrn/movie-reservation/pkg/errsx"
	"github.com/google/uuid"
)

type ServiceInterface interface {
	Create(ctx context.Context, movie InputDTO) (int64, error)
	Update(ctx context.Context, movie UpdateInputDTO) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]OutputDTO, error)
	ListByGenre(ctx context.Context, genres []string) ([]OutputDTO, error)
	ListByName(ctx context.Context, name string) ([]OutputDTO, error)
	Show(ctx context.Context, id int64) (OutputDTO, error)
}

type Service struct {
	repo RepoInterface
}

func NewService(movieRepo RepoInterface) Service {
	return Service{
		repo: movieRepo,
	}
}

func (s *Service) Create(ctx context.Context, input InputDTO) (int64, error) {
	var errs errsx.Map
	var err error
	err = NameValid(input.Title)
	if err != nil {
		errs.Set("title", err)
	}
	err = DescriptionValid(input.Description)
	if err != nil {
		errs.Set("description", err)
	}
	err = IsPosterValid(input.PosterURL)
	if err != nil {
		errs.Set("poster", err)
	}
	err = IsValidDuration(input.DurationMinutes)
	if err != nil {
		errs.Set("duration", err)
	}
	err = IsValidReleaseDate(input.ReleaseDate)
	if err != nil {
		errs.Set("release date", err)
	}
	err = IsValidRating(input.Rating)
	if err != nil {
		errs.Set("rating", err)
	}
	if errs != nil {
		return 0, fmt.Errorf("bad Request: %w", errs)
	}
	id, err := s.repo.Create(ctx, input)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (s *Service) List(ctx context.Context) ([]OutputDTO, error) {
	outputList, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return outputList, nil
}

func (s *Service) ListByGenre(ctx context.Context, genres []string) ([]OutputDTO, error) {
	sanitizedGenres := genre.SanitizeGenres(genres)
	outputList, err := s.repo.ListByGenre(ctx, sanitizedGenres)
	if err != nil {
		return nil, err
	}
	return outputList, err
}
func (s *Service) ListByName(ctx context.Context, name string) ([]OutputDTO, error) {
	outputList, err := s.repo.ListByName(ctx, name)
	if err != nil {
		return []OutputDTO{}, err
	}
	return outputList, nil
}

func (s *Service) Show(ctx context.Context, id int64) (OutputDTO, error) {
	output, err := s.repo.Show(ctx, id)
	if err != nil {
		return OutputDTO{}, nil
	}
	return output, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil{
		return fmt.Errorf("id no formato inválido")
	}
	err = s.repo.Delete(ctx, uid)
	if err != nil{
		return err
	}
	return nil
}

func (s *Service) Update(ctx context.Context, movie UpdateInputDTO) error{
	return nil
}
