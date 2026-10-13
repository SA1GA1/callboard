package listing

import (
	"github.com/gin-gonic/gin"
)

func (h *ListingHandler) RegisterRoutes(rg *gin.RouterGroup) {
	listings := rg.Group("/listings")
	{
		listings.GET("/page/:page", h.GetListingHandler)
		listings.GET("/:id", h.GetListingHandlerByID)
		listings.POST("/", h.CreateListingHandler)
		listings.PUT("/")
		listings.DELETE("/")
	}
}
