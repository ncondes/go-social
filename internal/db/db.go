package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq" // Register the pq driver
)

func New(
	addr string,
	maxOpenConns int,
	maxIdleConns int,
	maxIdleTime time.Duration,
) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
