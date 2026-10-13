package listing

import (
	"github.com/gin-gonic/gin"
)

func (h *ListingHandler) RegisterRoutes(rg *gin.RouterGroup) {
	listings := rg.Group("/listings")
	{
		listings.GET("/")
		listings.GET("/:id", h.GetListingHandlerByID)
		listings.POST("/", h.CreateListingHandler)
		listings.PUT("/")
		listings.DELETE("/")
	}
}
