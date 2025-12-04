package showtimes

import (
	"context"
	"time"
)

type ServiceInterface interface {
	Create(ctx context.Context, showtime InputDTO) (CreateResponse, error)
	List(ctx context.Context, day time.Time) ([]OutputDTO, error)
}

type Service struct {
	repo RepoInterface
}

type CreateResponse struct {
	ShowTimeID int64 `json:"show_time_id"`
}

func NewService(showtimeRepo RepoInterface) Service {
	return Service{
		repo: showtimeRepo,
	}
}

func (s *Service) Create(ctx context.Context, showtime InputDTO) (CreateResponse, error) {
	err := IsValidTimeRange(showtime.StartTime, showtime.EndTime)
	if err != nil {
		return CreateResponse{}, err
	}
	id, err := s.repo.Create(ctx, showtime)
	if err != nil {
		return CreateResponse{}, err
	}
	return CreateResponse{
		ShowTimeID: id,
	}, nil
}

func (s *Service) List(ctx context.Context, day time.Time) ([]OutputDTO, error) {
	outputList, err := s.repo.List(ctx, day)
	if err != nil {
		return nil, err
	}
	return outputList, nil
}
