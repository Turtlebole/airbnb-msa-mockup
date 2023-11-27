package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(routes *gin.Engine) {
	routes.Use(middleware.Authentication())
	routes.GET("/users/:user_id", controller.GetUser())
	routes.GET("/users/get", func(c *gin.Context) {
		controller.GetUsers(c)
	})
	routes.POST("/users/become-host/:user_id", controller.BecomeHost())
}
