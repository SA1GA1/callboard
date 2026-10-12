package listing

import (
	"callboard/internal/specifications"
	"time"
)

type Listing struct {
	ID             int64
	Title          string
	Description    string
	Photos         []string
	SellerID       int64
	Price          int64
	Specifications specifications.Specifications
	Place          string
	CreatedAt      time.Time
	IsSold         bool
}

type CreateListingRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Photos         []string `json:"photos"`
	Price          int64 `json:"price"`
	Specifications specifications.Specifications `json:"specifications"`
	Place          string `json:"place"`
}
