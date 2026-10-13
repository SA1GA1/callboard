package listing

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

	newId, err := CreateListing(c.Request.Context(), listing, h.pool)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newId)

}

func (h *ListingHandler) GetListingHandlerByID(c *gin.Context) {
	idParam, err:= strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	if idParam <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error":"id<=0"})
		return
	}

	listing, err := GetListingById(c.Request.Context(), idParam, h.pool)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, listing)

}

func (h *ListingHandler) GetListingHandler(c *gin.Context) {

}
