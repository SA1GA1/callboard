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
