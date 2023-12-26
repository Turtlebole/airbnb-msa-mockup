package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(routes *gin.Engine) {

	applyAuthentication := true

	if applyAuthentication {
		// Routes that require authentication
		routes.Use(middleware.Authentication())
		routes.GET("/users/:user_id", controllers.GetUser())
		routes.POST("/users/become-host/:user_id", controllers.BecomeHost())
		routes.POST("/users/edit-profile", controllers.EditProfile())
		routes.POST("/users/change-password", controllers.ChangePassword())
	}
	routes.DELETE("/users/delete/:user_id", controllers.DeleteUser())
	// Routes that don't require authentication
	routes.GET("/users/get", func(c *gin.Context) {
		controllers.GetUsers(c)
	})
}
