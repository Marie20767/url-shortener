package store

import (
	"database/sql"
	"fmt"
)

type Store struct {
	conn    *sql.DB
}

func connectDB(dbURL string) (*sql.DB, error) {
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	err = dbConn.Ping()
	if err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return dbConn, nil
}

func NewStore(dbURL string) (*Store, error) {
	dbConn, err := connectDB(dbURL)

	if err != nil {
		return nil, err
	}

	return &Store{
		conn:    dbConn,
	}, nil
}

func (s *Store) Close() error {
	return s.conn.Close()
}
