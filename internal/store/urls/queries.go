package urls

import (
	"context"
	"time"

	"github.com/Marie20767/url-shortener/internal/store/keys"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LongURL string

type URL struct {
	Key    keys.KeyValue `bson:"key"`
	URL    LongURL        `bson:"url"`
	Expiry time.Time     `bson:"expiry"`
}

const collection = "urls"

func (s *UrlStore) CreateShortURL(c context.Context, u *URL) error {
	collection := s.conn.Collection(collection)

	_, err := collection.InsertOne(c, u)

	if err != nil {
		return err
	}

	return nil
}

func (s *UrlStore) DeleteURLs(c context.Context) ([]keys.KeyValue, error) {
	collection := s.conn.Collection(collection)

	filter := bson.M{
		"expiry": bson.M{
			"$lte": time.Now(),
		},
	}

	var deletedKeys []keys.KeyValue

	for {
		var deleted URL

		err := collection.FindOneAndDelete(c, filter).Decode(&deleted)
		if err == mongo.ErrNoDocuments {
			break
		} else if err != nil {
			return nil, err
		}
		deletedKeys = append(deletedKeys, deleted.Key)
	}

	return deletedKeys, nil
}

func (s *UrlStore) GetLongURL(c context.Context, k keys.KeyValue) (LongURL, error) {
	var res URL
	collection := s.conn.Collection(collection)

	err := collection.FindOne(c, bson.M{"key": k}).Decode(&res)
	if err != nil {
		return "", err
	}

	return res.URL, nil
}
