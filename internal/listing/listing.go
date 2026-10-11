package listing

import (
	"callboard/internal/specifications"
	"errors"
	"strings"
	"time"
	"unicode/utf8"
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

func ValidationListing(listing Listing) error {

	titleLettersCount := utf8.RuneCountInString(strings.TrimSpace(listing.Title))
	descriptionLettersCount := utf8.RuneCountInString(strings.TrimSpace(listing.Description))
	price := listing.Price
	photosCount := len(listing.Photos)
	placeLettersCount := utf8.RuneCountInString(strings.TrimSpace(listing.Place))

	if titleLettersCount < 7 || titleLettersCount > 20 {
		return errors.New("ERROR: incorrect len of listing title")
	}

	if descriptionLettersCount < 10 || descriptionLettersCount > 1000 {
		return errors.New("ERROR: incorrect len of listing description")
	}

	if photosCount < 1 {
		return errors.New("ERROR: listing does not contain photos")
	}

	for _, photo := range listing.Photos {
		if len(strings.TrimSpace(photo)) == 0 {
			return errors.New("ERROR: incorrect URL of photo")
		}
	}

	if price < 99 || price > 9_999_999 {
		return errors.New("ERROR: incorrect price of listing")
	}

	if placeLettersCount == 0 || placeLettersCount > 20 {
		return errors.New("ERROR: incorrect place of listing")
	}

	return nil
}
