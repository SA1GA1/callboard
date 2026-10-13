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
		_, err = tx.Exec(ctx,
			"INSERT INTO photos (listing_id, url) VALUES ($1, $2)",
			newListingId, item.Photos[i],
		)
		if err != nil {
			return 0, fmt.Errorf("ERROR: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("ERROR: %w", err)
	}

	return newListingId, nil
}

func GetListingByIdInDB(ctx context.Context, id int64, pool *pgxpool.Pool) (Listing, error) {

	tx, err := pool.Begin(ctx)

	if err != nil {
		return Listing{}, fmt.Errorf("ERROR: incorrect create transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var l Listing

	query := `
		SELECT 
			l.id, l.title, l.description, l.seller_id, l.price, l.place, l.created_at, l.is_sold,
			COALESCE(array_agg(p.url ORDER BY p.id) FILTER (WHERE p.url IS NOT NULL), '{}'),
			CASE 
				WHEN s.listing_id IS NOT NULL THEN
					jsonb_build_object(
						'category', s.category,
						'brand', s.brand,
						'model', s.model,
						'condition', s.condition,
						'sellingreason', s.selling_reason
					)
				ELSE '{}'::jsonb
			END
		FROM listings l
		LEFT JOIN photos p ON l.id = p.listing_id
		LEFT JOIN specifications s ON l.id = s.listing_id
		WHERE l.id = $1
		GROUP BY l.id, s.listing_id, s.category, s.brand, s.model, s.condition, s.selling_reason;
	`

	err = tx.QueryRow(ctx, query, id).Scan(
		&l.ID,
		&l.Title,
		&l.Description,
		&l.SellerID,
		&l.Price,
		&l.Place,
		&l.CreatedAt,
		&l.IsSold,
		&l.Photos,
		&l.Specifications,
	)

	if err != nil {
		return Listing{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Listing{}, fmt.Errorf("ERROR: fail of commit: %w", err)
	}
	return l, nil
}

func GetListingWithPaginationInDB(ctx context.Context, pool *pgxpool.Pool, page int) ([]Listing, error) {

	if page < 1 {
		page = 1
	}
	const limit = 20

	offset := (page-1) * limit

	tx, err := pool.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("ERROR: incorrect create transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT 
			l.id, l.title, l.description, l.seller_id, l.price, l.place, l.created_at, l.is_sold,
			COALESCE(array_agg(p.url ORDER BY p.id) FILTER (WHERE p.url IS NOT NULL), '{}'),
			CASE 
				WHEN s.listing_id IS NOT NULL THEN
					jsonb_build_object(
						'category', s.category,
						'brand', s.brand,
						'model', s.model,
						'condition', s.condition,
						'sellingreason', s.selling_reason
					)
				ELSE '{}'::jsonb
			END
		FROM (
			SELECT id, title, description, seller_id, price, place, created_at, is_sold
			FROM listings
			WHERE is_sold = false
			ORDER BY created_at DESC, id DESC
			LIMIT $1 OFFSET $2
		) l
		LEFT JOIN photos p ON l.id = p.listing_id
		LEFT JOIN specifications s ON l.id = s.listing_id
		GROUP BY l.id, l.title, l.description, l.seller_id, l.price, l.place, l.created_at, l.is_sold,
		         s.listing_id, s.category, s.brand, s.model, s.condition, s.selling_reason
		ORDER BY l.created_at DESC, l.id DESC;
	`

	rows, err := tx.Query(ctx,query,limit,offset)
	if err != nil {
		return nil, fmt.Errorf("ERROR: fail of pagination query: %w",err)
	}
	defer rows.Close()

	listings := make([]Listing,0)

	for rows.Next() {
		var l Listing

		err = rows.Scan(
			&l.ID,
			&l.Title,
			&l.Description,
			&l.SellerID,
			&l.Price,
			&l.Place,
			&l.CreatedAt,
			&l.IsSold,
			&l.Photos,
			&l.Specifications,
		)

		if err != nil {
			return nil, fmt.Errorf("ERROR: fail to scan listing: %w", err)
		}

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ERROR: rows error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("ERROR: fail of commit: %w", err)
	}

	return listings,nil
}