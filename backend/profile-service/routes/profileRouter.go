package routes

import (
	"github.com/gin-gonic/gin"
	"profile-service/controllers"
	"profile-service/database"
)

func ProfileRoutes(router *gin.Engine) {
	router.POST("/profiles/create", controllers.CreateProfileHandler(database.Client))
	router.GET("/profiles/getall", controllers.GetAllProfilesHandler(database.Client))
	router.GET("/profiles/:id", controllers.GetProfileByIDHandler(database.Client))
	router.PUT("/profiles/:id", controllers.UpdateProfileHandler(database.Client))
	router.DELETE("/profiles/:id", controllers.DeleteProfileHandler(database.Client))
	router.GET("/profiles/by-email", controllers.GetProfileByEmailHandler(database.Client))
	// Add more routes as needed for updating, deleting accommodations, etc.
} // ok? aj pa sta Bog da
