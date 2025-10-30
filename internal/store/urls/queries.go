package urls

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UrlData struct {
	Key    string     `bson:"key_value"`
	Url    string     `bson:"url"`
	Expiry *time.Time `bson:"expiry,omitempty"`
}

func (s *UrlStore) Insert(ctx context.Context, urlData *UrlData) (any, error) {
	db := s.conn.Collection(s.collection)
	res, err := db.InsertOne(ctx, urlData)
	if err != nil {
		return "", fmt.Errorf("failed to insert new url into db: %w", err)
	}

	return res.InsertedID, nil
}

func (s *UrlStore) DeleteById(ctx context.Context, id any) error {
	db := s.conn.Collection(s.collection)
	_, err := db.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete url from db: %w", err)
	}

	return nil
}

func (s *UrlStore) DeleteExpired(ctx context.Context) ([]string, error) {
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
			return nil, fmt.Errorf("failed to delete expired url from db: %w", err)
		}
		deletedKeys = append(deletedKeys, deleted.Key)
	}

	return deletedKeys, nil
}

func (s *UrlStore) Get(ctx context.Context, key string) (string, error) {
	url, ok := s.cache.Get(ctx, key)
	if ok {
		return url, nil
	}

	var res UrlData
	db := s.conn.Collection(s.collection)
	err := db.FindOne(ctx, bson.M{"key_value": key}).Decode(&res)
	if err != nil {
		return "", fmt.Errorf("failed to fetch url from db: %w", err)
	}
	s.cache.Add(ctx, key, res.Url)

	return res.Url, nil
}
