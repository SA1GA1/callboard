package listing

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup) {
	listings := rg.Group("/listings")
	{
		listings.GET("/")
		listings.POST("/")
		listings.PUT("/")
		listings.DELETE("/")
	}
}
