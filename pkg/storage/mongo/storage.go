package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"memessy-api/pkg"
)

type MemeStorage struct {
	Collection *mongo.Collection
}

func (storage *MemeStorage) InsertOne(ctx context.Context, m pkg.Meme) (*pkg.Meme, error) {
	dbm := domainToDb(m)
	result, err := storage.Collection.InsertOne(ctx, dbm)
	if err != nil {
		return nil, err
	}
	objectId, _ := result.InsertedID.(primitive.ObjectID)
	m.Id = objectId.Hex()
	return &m, nil
}

func (storage *MemeStorage) FindMany(ctx context.Context, search string) ([]pkg.Meme, error) {
	var filter interface{}
	ops := options.Find()
	if search == "" {
		filter = primitive.D{}
	} else {
		filter = primitive.D{{
			"$text",
			primitive.D{{"$search", search}}},
		}
		ops.
			SetSort(primitive.M{"score": -1})
	}
	cursor, err := storage.Collection.Find(ctx, filter, ops)
	if err != nil {
		return nil, err
	}
	var dbMemes []dbMeme
	err = cursor.All(ctx, &dbMemes)
	if err != nil {
		return nil, err
	}
	memes := make([]pkg.Meme, 0, len(dbMemes))
	for _, m := range dbMemes {
		memes = append(memes, m.toDomain())
	}
	return memes, nil
}

func (storage *MemeStorage) FindOne(ctx context.Context, id string) (*pkg.Meme, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	result := storage.Collection.FindOne(
		ctx,
		primitive.D{{"_id", objectId}},
	)
	meme := dbMeme{}
	err := result.Decode(&meme)
	if err != nil {
		return nil, err
	}
	d := meme.toDomain()
	return &d, nil
}

func (storage *MemeStorage) UpdateOne(ctx context.Context, id string, meme pkg.Meme) (*pkg.Meme, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	dbm := domainToDb(meme)
	result := storage.Collection.FindOneAndUpdate(
		ctx,
		primitive.D{{"_id", objectId}},
		primitive.D{{"$set",
			dbm,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	m := dbMeme{}
	err := result.Decode(&m)
	if err != nil {
		return nil, err
	}
	d := m.toDomain()
	return &d, nil
}

func (storage *MemeStorage) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := storage.Collection.DeleteOne(
		ctx,
		primitive.D{{"_id", objectId}},
	)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("not found")
	}
	return nil
}

func (storage *MemeStorage) Init() error {
	storage.Collection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys: primitive.D{{"parsed_text", "text"}},
			Options: options.Index().
				SetName("searchIndex").
				SetDefaultLanguage("russian"),
		},
	)
	return nil
}
