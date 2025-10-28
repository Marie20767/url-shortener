package urls

import (
	"context"
	"time"

	"github.com/Marie20767/url-shortener/internal/utils/cache"
	"github.com/Marie20767/url-shortener/internal/utils/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func New(cfg *config.Url) (*UrlStore, error) {
	dbConn, err := connectDb(cfg.DbUrl, cfg.DbName)
	if err != nil {
		return nil, err
	}

	cache, err := cache.New(cfg.CacheCapacity)
	if err != nil {
		return nil, err
	}

	return &UrlStore{
		conn:       dbConn,
		collection: "urls",
		cache:      cache,
	}, nil
}

func (s *UrlStore) Close(ctx context.Context) error {
	return s.conn.Client().Disconnect(ctx)
}
