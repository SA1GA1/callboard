package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"callboard/internal/listing"
	"callboard/internal/specifications"

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
	}

	defer pool.Close()

	testListing := listing.Listing {
		Title: "2крутая книга",
		Description: "2очень интересная книга",
		Photos: []string{
			"https://i.pinimg.com/736x/d9/fb/1b/d9fb1bf9a64cc80675801ba9a61840e7.jpg",
			"https://i.pinimg.com/736x/d9/fb/1b/d9fb1bf9a64cc80675801ba9a61840e7.jpg",
		},
		SellerID: 2,
		Price: 1000,
		Specifications: specifications.Specifications{Category: "2книги"},
		Place: "2Казань",
	}

	if err := listing.ValidationListing(testListing); err != nil {
		fmt.Printf("ERROR: incorrect validation for listing: %v\n", err)
		os.Exit(3)
	} 
	
	newListingId, err := createListingInDB(ctx, testListing, pool)

	if err != nil {
		fmt.Printf("ERROR: incorrect save listing in db: %v", err)
	} else {
		fmt.Printf("ID of created listing: %d", newListingId)
	}


}

func initDataBase(ctx context.Context) (*pgxpool.Pool, error) {

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

func createListingInDB(ctx context.Context, item listing.Listing, pool *pgxpool.Pool) (int64, error) {

	tx, err := pool.Begin(ctx)

	if err != nil {
		return 0, fmt.Errorf("ERROR: incorrect create transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	var newListingId int64

	err = tx.QueryRow(ctx,
		"INSERT INTO listings (title, description, seller_id, price, place) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		item.Title, item.Description, item.SellerID, item.Price, item.Place,
	).Scan(&newListingId)

	if err != nil {
		return 0, fmt.Errorf("ERROR: %w", err)
	}

	if item.Specifications.Category != "" || item.Specifications.Brand != "" || item.Specifications.Model != "" || item.Specifications.Condition != "" || item.Specifications.SellingReason != "" {

		_, err = tx.Exec(ctx,
			"INSERT INTO specifications (listing_id, category, brand, model, condition, selling_reason) VALUES ($1, $2, $3, $4, $5, $6)",
			newListingId, item.Specifications.Category, item.Specifications.Brand, item.Specifications.Model, item.Specifications.Condition, item.Specifications.SellingReason,
		)

		if err != nil {
			return 0, fmt.Errorf("ERROR: %w", err)
		}
	}

	for i := range item.Photos {
		_,err = tx.Exec(ctx,
			"INSERT INTO photos (listing_id, url) VALUES ($1, $2)",
			newListingId, item.Photos[i],
		)
		if err != nil {
			return 0, fmt.Errorf("ERROR: %w", err)
		}
	}

	if err:=tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("ERROR: %w", err)
	}

	return newListingId, nil
}
