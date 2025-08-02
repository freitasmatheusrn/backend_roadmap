package usecase

import (
	"context"

	"github.com/freitasmatheusrn/url-shortening-app/internal/dtos"
	"github.com/freitasmatheusrn/url-shortening-app/internal/repo"
)

type ListAllUrlUseCase struct {
	UrlRepository repo.UrlRepositoryInterface
}

func NewListAllUrlUseCase(UrlRepository repo.UrlRepositoryInterface) *ListAllUrlUseCase {
	return &ListAllUrlUseCase{
		UrlRepository: UrlRepository,
	}
}

func (uc *ListAllUrlUseCase) Execute(ctx context.Context) ([]dtos.UrlOutputDTO, error) {
	urls, err := uc.UrlRepository.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	dtoList := make([]dtos.UrlOutputDTO, 0)
	for _, url := range urls {
		dto := dtos.UrlOutputDTO{
			ID:          url.ID,
			OriginalURL: url.OriginalURL,
			AccessCount: url.AccessCount,
		}
		dtoList = append(dtoList, dto)
	}
	return dtoList, nil
}
