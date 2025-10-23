package urls

import (
	"context"
	"time"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UrlData struct {
	Key    string     `bson:"key_value"`
	Url    string     `bson:"url"`
	Expiry *time.Time `bson:"expiry,omitempty"`
}

func (s *UrlStore) Insert(ctx context.Context, urlData *UrlData) error {
	db := s.conn.Collection(s.collection)
	_, err := db.InsertOne(ctx, urlData)
	if err != nil {
		return err
	}

	return nil
}

func (s *UrlStore) Delete(ctx context.Context) ([]string, error) {
	db := s.conn.Collection(s.collection)

	filter := bson.M{
		"expiry": bson.M{
			"$lte": time.Now(),
		},
	}

	var deletedKeys []string

	for {
		var deleted UrlData

		err := db.FindOneAndDelete(ctx, filter).Decode(&deleted)
		if err == mongo.ErrNoDocuments {
			break
		} else if err != nil {
			return nil, err
		}
		deletedKeys = append(deletedKeys, deleted.Key)
	}

	return deletedKeys, nil
}

func (s *UrlStore) Get(ctx context.Context, key string) (string, error) {
	url, err := s.cache.Get(key)
	if err == nil {
		return url, nil
	}

	var res UrlData
	db := s.conn.Collection(s.collection)
	err = db.FindOne(ctx, bson.M{"key_value": key}).Decode(&res)
	if err != nil {
		return "", err
	}
	err = s.cache.Add(key, res.Url)
	if err != nil {
		log.Error("Failed to add url to cache", err)
	}

	return res.Url, nil
}
