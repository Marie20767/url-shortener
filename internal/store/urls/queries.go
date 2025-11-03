package urls

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Marie20767/url-shortener/internal/store/urls/model"
)

var ErrNotFound = errors.New("url not found")

func (s *UrlStore) Insert(ctx context.Context, urlData *model.UrlData) (any, error) {
	if urlData.Expiry != nil {
		utcTime := urlData.Expiry.UTC()
		urlData.Expiry = &utcTime
	}

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
			"$lte": time.Now().UTC(),
		},
	}

	var deletedKeys []string
	for {
		var deleted *model.UrlData
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

func (s *UrlStore) Get(ctx context.Context, key string, currentTimeStamp time.Time) (string, error) {
	url, ok := s.cache.Get(ctx, key)
	if ok {
		return url, nil
	}

	filters := bson.M{
		"key_value": key,
		"$or": []bson.M{
			{"expiry": bson.M{"$gte": time.Now().UTC()}},
			{"expiry": bson.M{"$exists": false}},
		},
	}

	var res *model.UrlData
	db := s.conn.Collection(s.collection)
	err := db.FindOne(ctx, filters).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", ErrNotFound
		}

		return "", fmt.Errorf("failed to fetch url from db: %w", err)
	}
	s.cache.Add(ctx, res, currentTimeStamp)

	return res.Url, nil
}
