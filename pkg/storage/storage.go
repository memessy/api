package storage

import (
	"context"
	"memessy-api/pkg"
)

type Storage interface {
	InsertOne(context.Context, pkg.Meme) (*pkg.Meme, error)
	FindMany(context.Context, string) ([]pkg.Meme, error)
	FindOne(context.Context, string) (*pkg.Meme, error)
	UpdateOne(context.Context, string, pkg.Meme) (*pkg.Meme, error)
	Delete(ctx context.Context, id string) error
}
