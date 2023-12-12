// routes.go

package routes

import (
	"accommodation-service/controllers"
	"github.com/gin-gonic/gin"
)

func AccommodationRoutes(router *gin.Engine) {
	router.POST("/accommodations/create", controllers.CreateAccommodation())
	router.GET("/accommodations/:accommodation_id", controllers.GetAccommodation())
	router.GET("/accommodations", controllers.GetAllAccommodations())

	// Add more routes as needed for updating, deleting accommodations, etc.
} // ok? aj pa sta Bog da
