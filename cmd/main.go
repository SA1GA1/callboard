package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"callboard/internal/database"
	"callboard/internal/listing"

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

	listingHandler := listing.NewListingHandler(pool)

	router := gin.Default()

	mainGroup := router.Group("")

	listingHandler.RegisterRoutes(mainGroup)

	if err := router.Run(":8081"); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
}
