package impl

import (
	"context"
	"memessy-api/pkg"
	"memessy-api/pkg/bus"
	"memessy-api/pkg/fileserver"
	"memessy-api/pkg/storage"
	"time"
)

type Service struct {
	storage    storage.Storage
	fileServer fileserver.FileServer
	eventBus   bus.EventBus
}

func NewService(
	storage storage.Storage,
	fileServer fileserver.FileServer,
	eventBus bus.EventBus,
) *Service {
	return &Service{
		storage:    storage,
		fileServer: fileServer,
		eventBus:   eventBus,
	}
}

func (s *Service) List(ctx context.Context, searchQuery string) ([]pkg.Meme, error) {
	memes, err := s.storage.FindMany(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	return memes, nil
}

func (s *Service) Create(ctx context.Context, name string, data []byte) (*pkg.Meme, error) {
	url, err := s.fileServer.Upload(fileserver.File{
		Name: name,
		Data: data,
	})
	if err != nil {
		return nil, err
	}
	meme, err := s.storage.InsertOne(ctx, pkg.Meme{
		FileUrl:   *url,
		CreatedAt: time.Now(),
	})
	if err != nil {
		// TODO remove file
		return nil, err
	}
	s.eventBus.Created(*meme)
	return meme, nil
}

func (s *Service) Retrieve(ctx context.Context, id string) (*pkg.Meme, error) {
	meme, err := s.storage.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return meme, nil
}

func (s *Service) Update(ctx context.Context, id string, meme pkg.Meme) (*pkg.Meme, error) {
	updated, err := s.storage.UpdateOne(ctx, id, meme)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.storage.Delete(ctx, id)
	// TODO get result of delete, remove file
	return err
}
