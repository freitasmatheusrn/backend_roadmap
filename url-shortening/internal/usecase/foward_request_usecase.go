package usecase

import (
	"context"

	"github.com/freitasmatheusrn/url-shortening-app/internal/repo"
)

type FowardRequestUseCase struct {
	UrlRepository repo.UrlRepositoryInterface
}

func NewFowardRequestUseCase(UrlRepository repo.UrlRepositoryInterface) *FowardRequestUseCase {
	return &FowardRequestUseCase{
		UrlRepository: UrlRepository,
	}
}

func (uc *FowardRequestUseCase) Execute(ctx context.Context, id string) (string, error) {
	originalURL, err := uc.UrlRepository.GetOne(ctx, id)
	if err != nil{
		return "Error Retrieving record", err
	}
	return originalURL, nil
}
