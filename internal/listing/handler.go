package listing

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListingHandler struct {
	pool *pgxpool.Pool
}

func NewListingHandler(pool *pgxpool.Pool) *ListingHandler {
	return &ListingHandler{pool: pool}
}

func (h *ListingHandler) CreateListingHandler(c *gin.Context) {
	var inputListing CreateListingRequest

	if err := c.ShouldBindJSON(&inputListing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ValidationListing(inputListing)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	listing := Listing{
		Title:          inputListing.Title,
		Description:    inputListing.Description,
		Photos:         inputListing.Photos,
		Price:          inputListing.Price,
		Specifications: inputListing.Specifications,
		Place:          inputListing.Place,

		SellerID: 123, // тестовая заглушка
	}

	newId, err := CreateListingInDB(c.Request.Context(), listing, h.pool)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newId)

}
