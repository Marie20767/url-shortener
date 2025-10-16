package urls

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


type urlData struct {
	key    string `bson:"key_value"`
	url    string       `bson:"url"`
	expiry time.Time     `bson:"expiry"`
}

func (s *UrlStore) InsertKey(ctx context.Context, key *urlData) error {
	db := s.conn.Collection(s.collection)
	_, err := db.InsertOne(ctx, key)
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

func (s *UrlStore) GetUrl(ctx context.Context, key string) (string, error) {
	var res urlData
	db := s.conn.Collection(s.collection)
	err := db.FindOne(ctx, bson.M{"key_value": key}).Decode(&res)
	if err != nil {
		return "", err
	}

	return res.url, nil
}
