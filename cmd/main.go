package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"callboard/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := database.InitDataBase(ctx)

	if err != nil {
		fmt.Printf("ERROR: connect with db does not correct: %v\n", err)
		cancel()
		os.Exit(3)
	}

	defer pool.Close()

	// testListing := listing.Listing {
	// 	Title: "2крутая книга",
	// 	Description: "2очень интересная книга",
	// 	Photos: []string{
	// 		"https://i.pinimg.com/736x/d9/fb/1b/d9fb1bf9a64cc80675801ba9a61840e7.jpg",
	// 		"https://i.pinimg.com/736x/d9/fb/1b/d9fb1bf9a64cc80675801ba9a61840e7.jpg",
	// 	},
	// 	SellerID: 2,
	// 	Price: 1000,
	// 	Specifications: specifications.Specifications{Category: "2книги"},
	// 	Place: "2Казань",
	// }

	// if err := listing.ValidationListing(testListing); err != nil {
	// 	fmt.Printf("ERROR: incorrect validation for listing: %v\n", err)
	// 	os.Exit(3)
	// }

	// newListingId, err := createListingInDB(ctx, testListing, pool)

	// if err != nil {
	// 	fmt.Printf("ERROR: incorrect save listing in db: %v", err)
	// } else {
	// 	fmt.Printf("ID of created listing: %d", newListingId)
	// }

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run(":8081")
}
