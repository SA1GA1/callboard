package database

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDataBase(ctx context.Context) (*pgxpool.Pool, error) {

	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return nil, errors.New("ERROR: incorrect env params")
	}

	url := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(dbUser, dbPassword),
		Host: net.JoinHostPort(dbHost,dbPort),
		Path: dbName,
	}

	connStr := url.String()

	pool, err := pgxpool.New(ctx, connStr)

	if err != nil {
		return nil, fmt.Errorf("ERROR: pool for postgres not created: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		defer pool.Close()
		return nil, fmt.Errorf("ERROR: data base unavailable: %w", err)
	}

	return pool, nil
}