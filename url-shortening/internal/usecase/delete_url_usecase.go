package usecase

import (
	"context"

	"github.com/freitasmatheusrn/url-shortening-app/internal/repo"
)

type DeleteUrlUseCase struct {
	UrlRepository repo.UrlRepositoryInterface
}

func NewDeleteUrlUseCase(UrlRepository repo.UrlRepositoryInterface) *DeleteUrlUseCase {
	return &DeleteUrlUseCase{
		UrlRepository: UrlRepository,
	}
}

func (uc *DeleteUrlUseCase) Execute(ctx context.Context, id string) error {
	err := uc.UrlRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
