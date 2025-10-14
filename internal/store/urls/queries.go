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
	key    keys.KeyValue `bson:"key_value"`
	url    longUrl       `bson:"url"`
	expiry time.Time     `bson:"expiry"`
}

func (s *UrlStore) CreateShortUrl(ctx context.Context, url *urlData) error {
	db := s.conn.Collection(s.collection)

	_, err := db.InsertOne(ctx, url)

	if err != nil {
		return err
	}

	return nil
}

func (s *UrlStore) DeleteUrls(ctx context.Context) ([]keys.KeyValue, error) {
	db := s.conn.Collection(s.collection)

	filter := bson.M{
		"expiry": bson.M{
			"$lte": time.Now(),
		},
	}

	var deletedKeys []keys.KeyValue

	for {
		var deleted urlData

		err := db.FindOneAndDelete(ctx, filter).Decode(&deleted)
		if err == mongo.ErrNoDocuments {
			break
		} else if err != nil {
			return nil, err
		}
		deletedKeys = append(deletedKeys, deleted.key)
	}

	return deletedKeys, nil
}

func (s *UrlStore) GetLongUrl(ctx context.Context, key keys.KeyValue) (longUrl, error) {
	var res urlData
	db := s.conn.Collection(s.collection)

	err := db.FindOne(ctx, bson.M{"key_value": key}).Decode(&res)
	if err != nil {
		return "", err
	}

	return res.url, nil
}
