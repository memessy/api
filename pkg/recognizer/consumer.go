package recognizer

import (
	"context"
	"github.com/rs/zerolog/log"
	"memessy-api/pkg/bus"
	"memessy-api/pkg/fileserver"
	"memessy-api/pkg/storage"
	"path"
)

type Consumer struct {
	Recognizer Recognizer
	FileServer fileserver.FileServer
	Storage    storage.Storage
}

func (c *Consumer) Consume(event bus.EventCreated) {
	entry := log.Info().Str("id", event.Meme.Id)
	entry.Msg("started recognition")
	file, err := c.FileServer.Download(path.Base(event.Meme.FileUrl.Path))
	if err != nil {
		log.Error().Err(err).Msg("caught error while getting meme file from file server")
		return
	}
	entry.Msg("got picture from server")
	text, err := c.Recognizer.Recognize(file.Data)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	entry.Str("text", text).Msg("recognized")
	meme, err := c.Storage.FindOne(context.Background(), event.Meme.Id)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	entry.Msg("found meme in db")
	meme.ParsedText = text
	_, err = c.Storage.UpdateOne(context.Background(), event.Meme.Id, *meme)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	entry.Msg("updated meme with parsed text")
}
