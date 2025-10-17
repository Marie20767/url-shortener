package urls

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UrlData struct {
	Key    string    `bson:"key_value"`
	Url    string    `bson:"url" validate:"required,url"`
	Expiry time.Time `bson:"expiry,omitempty" validate:"expiry"`
}

func (s *UrlStore) InsertUrlData(ctx context.Context, urlData *UrlData) error {
	db := s.conn.Collection(s.collection)
	_, err := db.InsertOne(ctx, urlData)
	if err != nil {
		return err
	}

	return nil
}

func (s *UrlStore) DeleteUrlData(ctx context.Context) ([]string, error) {
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

func (s *UrlStore) GetUrl(ctx context.Context, key string) (string, error) {
	var res UrlData
	db := s.conn.Collection(s.collection)
	err := db.FindOne(ctx, bson.M{"key_value": key}).Decode(&res)
	if err != nil {
		return "", err
	}

	return res.Url, nil
}
