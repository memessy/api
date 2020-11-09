package meme

import (
	"memessy-api/pkg"
	"time"
)

func CreateRetrieveSchema(meme *pkg.Meme) RetrieveSchema {
	return RetrieveSchema{
		ID:          meme.Id,
		FileURL:     meme.FileUrl.String(),
		Description: meme.Description,
		ParsedText:  meme.ParsedText,
		Categories:  meme.Categories,
		ProcessedAt: meme.ProcessedAt,
		CreatedAt:   meme.CreatedAt,
	}
}

func CreateListSchema(memes []pkg.Meme) []RetrieveSchema {
	schemas := make([]RetrieveSchema, 0, len(memes))
	for _, m := range memes {
		schemas = append(schemas, CreateRetrieveSchema(&m))
	}
	return schemas
}


type RetrieveSchema struct {
	ID          string    `json:"id"`
	FileURL     string    `json:"file_url"`
	Description string    `json:"description"`
	ParsedText  string    `json:"parsed_text"`
	Categories  []string  `json:"categories"`
	ProcessedAt time.Time `json:"processed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateSchema struct {
	Description string    `json:"description"`
	ParsedText  string    `json:"parsed_text"`
	Categories  []string  `json:"categories"`
}