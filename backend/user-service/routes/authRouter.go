package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(routes *gin.Engine) {
	routes.Use(func(c *gin.Context) {
		middleware.Rbac(c)
	})
	routes.POST("/users/register", controller.Register())
	routes.POST("/users/login", func(c *gin.Context) {
		controller.Login(c)
	})
	routes.POST("/users/logout", controller.Logout())
}
