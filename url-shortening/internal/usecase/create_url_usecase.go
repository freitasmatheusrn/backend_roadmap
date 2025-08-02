package usecase

import (
	"context"

	"github.com/freitasmatheusrn/url-shortening-app/internal/collections"
	"github.com/freitasmatheusrn/url-shortening-app/internal/dtos"
	"github.com/freitasmatheusrn/url-shortening-app/internal/repo"
	"github.com/freitasmatheusrn/url-shortening-app/pkg"
)

type CreateUrlUseCase struct {
	UrlRepository repo.UrlRepositoryInterface
}

func NewCreateUrlUseCase(UrlRepository repo.UrlRepositoryInterface) *CreateUrlUseCase {
	return &CreateUrlUseCase{
		UrlRepository: UrlRepository,
	}
}

func (uc *CreateUrlUseCase) Execute(ctx context.Context, input dtos.UrlInputDTO) error {
	lenght := 8
	id, err := pkg.GenerateShortID(lenght)
	if err != nil {
		return err
	}
	u := collections.Url{
		ID:          id,
		OriginalURL: input.OriginalURL,
		AccessCount: 0,
	}

	err = uc.UrlRepository.Create(ctx, u)
	if err != nil {
		return err
	}
	return nil
}
