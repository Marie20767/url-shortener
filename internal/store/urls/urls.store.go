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

func connectDb(dbURL string) (*mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(dbURL).SetConnectTimeout(5 * time.Second)
	mongoClient, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	dbName := os.Getenv("URLS_DB_Name")
	return mongoClient.Database(dbName), nil
}

func NewStore(dbURL string) (*UrlStore, error) {
	dbConn, err := connectDb(dbURL)

	if err != nil {
		return nil, err
	}

	return &UrlStore{
		conn: dbConn,
	}, nil
}

func (s *UrlStore) Close(ctx context.Context) {
	s.conn.Client().Disconnect(ctx)
}
