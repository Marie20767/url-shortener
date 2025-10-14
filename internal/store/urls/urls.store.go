package urls

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UrlStore struct {
	conn *mongo.Database
}

func connectDb(dbUrl string) (*mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(dbUrl).SetConnectTimeout(5 * time.Second)
	mongoClient, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	dbName := os.Getenv("URL_DB_NAME")
	return mongoClient.Database(dbName), nil
}

func NewStore(dbUrl string) (*UrlStore, error) {
	dbConn, err := connectDb(dbUrl)

	if err != nil {
		return nil, err
	}

	return &UrlStore{
		conn: dbConn,
	}, nil
}

func (s *UrlStore) Close(ctx context.Context) error {
	return s.conn.Client().Disconnect(ctx)
}
