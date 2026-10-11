package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := checkPostgres(); err != nil {
		log.Fatalf("%v", err)
	} else {
		fmt.Println("SUCCESS: connect with db is correct")
	}
}

func checkPostgres() error {

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return errors.New("ERROR: incorrect env params")
	}

	url := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(dbUser, dbPassword),
		Host: net.JoinHostPort(dbHost,dbPort),
		Path: dbName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	connStr := url.String()

	pool, err := pgxpool.New(ctx, connStr)

	if err != nil {
		return errors.New("ERROR: pool for postgres not created")
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		return errors.New("ERROR: data base unavailable")
	}

	return nil
}
