package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(routes *gin.Engine) {
	// Routes that require authentication
	applyAuthentication := false

	if applyAuthentication {
		routes.Use(AuthenticationMiddleware())
	}

	routes.GET("/users/:user_id", controllers.GetUser())
	routes.POST("/users/become-host/:user_id", controllers.BecomeHost())

	// Routes that don't require authentication
	routes.GET("/users/get", func(c *gin.Context) {
		controllers.GetUsers(c)
	})
}

// AuthenticationMiddleware is a placeholder for your authentication logic
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement your authentication logic here if needed
		// This is just a placeholder
		c.Next()
	}
}
