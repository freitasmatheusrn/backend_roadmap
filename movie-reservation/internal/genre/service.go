package genre

import "context"

type ServiceInterface interface {
	Create(ctx context.Context, genre InputDTO) (OutputDTO, error)
	List(ctx context.Context) ([]OutputDTO, error)
}

type Service struct {
	repo RepoInterface
}


func NewService(genreRepo RepoInterface) Service {
	return Service{
		repo: genreRepo,
	}
}

func (s *Service) Create(ctx context.Context, genre InputDTO) (OutputDTO, error) {
	err := NameValid(genre.Name)
	if err != nil {
		return OutputDTO{}, err
	}
	output, err := s.repo.Create(ctx, genre)
	if err != nil {
		return OutputDTO{}, err
	}
	return output, nil
}

func (s *Service) List(ctx context.Context) ([]OutputDTO, error) {
	outputList, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return outputList, nil
}
