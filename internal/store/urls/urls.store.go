package urls

import (
	"context"
	"time"

	"github.com/Marie20767/url-shortener/internal/utils/cache/lru"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	cacheCapacity = 2
)

type UrlStore struct {
	conn       *mongo.Database
	collection string
	cache      *cache.LRUCache
}

func connectDb(dbUrl, dbName string) (*mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(dbUrl).SetConnectTimeout(5 * time.Second)
	mongoClient, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	return mongoClient.Database(dbName), nil
}

func New(dbUrl, dbName string) (*UrlStore, error) {
	dbConn, err := connectDb(dbUrl, dbName)
	if err != nil {
		return nil, err
	}

	lru, err := cache.New(cacheCapacity)
	if err != nil {
		return nil, err
	}

	return &UrlStore{
		conn:       dbConn,
		collection: "urls",
		cache:      lru,
	}, nil
}

func (s *UrlStore) Close(ctx context.Context) error {
	return s.conn.Client().Disconnect(ctx)
}
