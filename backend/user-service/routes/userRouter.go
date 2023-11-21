package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
<<<<<<< HEAD
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
=======
func UserRoutes(routes *gin.Engine) {
	routes.Use(middleware.Authentication())
	routes.GET("/users/:user_id", controller.GetUser())
	routes.GET("/users/get", func(c *gin.Context) {
		controller.GetUsers(c)
	})
>>>>>>> main
}
