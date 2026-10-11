package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx,cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := initDataBase(ctx)
	if err != nil {
		fmt.Printf("ERROR: connect with db does not correct: %v\n",err)
		cancel()
		os.Exit(3)
	} else {
		fmt.Println("SUCCESS: connect with db is correct")
	}
	defer pool.Close()
}

func initDataBase(ctx context.Context) (*pgxpool.Pool, error) {

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

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
