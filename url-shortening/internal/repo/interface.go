package repo

import (
	"context"

	"github.com/freitasmatheusrn/url-shortening-app/internal/collections"
)

type UrlRepositoryInterface interface {
	Create(ctx context.Context, url collections.Url) error
	ListAll(ctx context.Context) ([]collections.Url, error)
	IncrementAccess(ctx context.Context, id string) error
	Update(ctx context.Context, url collections.Url) error
	Delete(ctx context.Context, id string) error
	GetOne(ctx context.Context, id string) (string, error)
}
