package urls

import (
	"context"
	"time"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type longUrl string

type urlData struct {
	Key    keys.KeyValue `bson:"key_value"`
	Url    longUrl       `bson:"url"`
	Expiry time.Time     `bson:"expiry"`
}

const collection = "urls"

func (s *UrlStore) CreateShortUrl(ctx context.Context, url *urlData) error {
	collection := s.conn.Collection(collection)

	_, err := collection.InsertOne(ctx, url)

	if err != nil {
		return err
	}

	return nil
}

func (s *UrlStore) DeleteUrls(ctx context.Context) ([]keys.KeyValue, error) {
	collection := s.conn.Collection(collection)

	filter := bson.M{
		"expiry": bson.M{
			"$lte": time.Now(),
		},
	}

	var deletedKeys []keys.KeyValue

	for {
		var deleted urlData

		err := collection.FindOneAndDelete(ctx, filter).Decode(&deleted)
		if err == mongo.ErrNoDocuments {
			break
		} else if err != nil {
			return nil, err
		}
		deletedKeys = append(deletedKeys, deleted.Key)
	}

	return deletedKeys, nil
}

func (s *UrlStore) GetLongUrl(ctx context.Context, key keys.KeyValue) (longUrl, error) {
	var res urlData
	collection := s.conn.Collection(collection)

	err := collection.FindOne(ctx, bson.M{"key_value": key}).Decode(&res)
	if err != nil {
		return "", err
	}

	return res.Url, nil
}
