package listing

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateListingInDB(ctx context.Context, item Listing, pool *pgxpool.Pool) (int64, error) {

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
