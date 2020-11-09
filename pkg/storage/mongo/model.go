package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"memessy-api/pkg"
	"net/url"
)

type dbMeme struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	FileURL     primitive.Binary   `bson:"file_url,omitempty"`
	Description string             `bson:"description,omitempty"`
	ParsedText  string             `bson:"parsed_text,omitempty"`
	Categories  []string           `bson:"categories,omitempty"`
	ProcessedAt primitive.DateTime `bson:"processed_at,omitempty"`
	CreatedAt   primitive.DateTime `bson:"created_at,omitempty"`
}

func domainToDb(m pkg.Meme) dbMeme {
	id, _ := primitive.ObjectIDFromHex(m.Id)
	binaryUrl, _ := m.FileUrl.MarshalBinary()
	return dbMeme{
		Id: id,
		FileURL: primitive.Binary{
			Data: binaryUrl,
		},
		Description: m.Description,
		ParsedText:  m.ParsedText,
		Categories:  m.Categories,
		ProcessedAt: primitive.NewDateTimeFromTime(m.ProcessedAt),
		CreatedAt:   primitive.NewDateTimeFromTime(m.CreatedAt),
	}
}

func (m dbMeme) toDomain() pkg.Meme {
	unmarshalledUrl := &url.URL{}
	_ = unmarshalledUrl.UnmarshalBinary(m.FileURL.Data)
	return pkg.Meme{
		Id:          m.Id.Hex(),
		FileUrl:     *unmarshalledUrl,
		Description: m.Description,
		ParsedText:  m.ParsedText,
		Categories:  m.Categories,
		ProcessedAt: m.ProcessedAt.Time(),
		CreatedAt:   m.CreatedAt.Time(),
	}
}
