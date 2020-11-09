package rest

import (
	"context"
	"memessy-api/pkg"
)

type Service interface {
	List(ctx context.Context, searchQuery string) ([]pkg.Meme, error)
	Create(ctx context.Context, name string, data []byte) (*pkg.Meme, error)
	Retrieve(ctx context.Context, id string) (*pkg.Meme, error)
	Update(ctx context.Context, id string, meme pkg.Meme) (*pkg.Meme, error)
	Delete(ctx context.Context, id string) error
}
