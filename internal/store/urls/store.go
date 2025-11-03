package urls

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Marie20767/url-shortener/internal/cache/urls"
	"github.com/Marie20767/url-shortener/internal/utils/config"
)

type UrlStore struct {
	conn       *mongo.Database
	collection string
	cache      *cache.Cache
}

func connectDb(cfg *config.Url) (*mongo.Database, error) {
	timeOut := time.Duration(cfg.DbTimeout) * time.Second
	clientOpts := options.Client().ApplyURI(cfg.DbUrl).SetConnectTimeout(timeOut)
	mongoClient, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to url db: %w", err)
	}

	return mongoClient.Database(cfg.DbName), nil
}

func New(cfg *config.Url) (*UrlStore, error) {
	dbConn, err := connectDb(cfg)
	if err != nil {
		return nil, err
	}

	newCache, err := cache.New(cfg.CacheUrl)
	if err != nil {
		return nil, err
	}

	return &UrlStore{
		conn:       dbConn,
		collection: "urls",
		cache:      newCache,
	}, nil
}

func (s *UrlStore) Close(ctx context.Context) error {
	return s.conn.Client().Disconnect(ctx)
}
