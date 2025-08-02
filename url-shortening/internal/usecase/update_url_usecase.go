package usecase

import (
	"context"

	"github.com/freitasmatheusrn/url-shortening-app/internal/collections"
	"github.com/freitasmatheusrn/url-shortening-app/internal/dtos"
	"github.com/freitasmatheusrn/url-shortening-app/internal/repo"
)

type UpdateUrlUseCase struct {
	UrlRepository repo.UrlRepositoryInterface
}

func NewUpdateUrlUseCase(UrlRepository repo.UrlRepositoryInterface) *UpdateUrlUseCase {
	return &UpdateUrlUseCase{
		UrlRepository: UrlRepository,
	}
}

func (uc *UpdateUrlUseCase) Execute(ctx context.Context, input dtos.UpdateUrlInputDTO) error {
	u := collections.Url{
		ID: input.ID,
		OriginalURL: input.OriginalURL,
	}
	err := uc.UrlRepository.Update(ctx, u)
	if err != nil {
		return err
	}
	return nil
}
