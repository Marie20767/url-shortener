package keys

import (
	"database/sql"
	"fmt"
)

type KeyStore struct {
	conn *sql.DB
}

func connectDb(dbUrl string) (*sql.DB, error) {
	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return dbConn, nil
}

func NewStore(dbUrl string) (*KeyStore, error) {
	dbConn, err := connectDb(dbUrl)
	if err != nil {
		return nil, err
	}

	return &KeyStore{
		conn: dbConn,
	}, nil
}

func (s *KeyStore) Close() error {
	return s.conn.Close()
}
