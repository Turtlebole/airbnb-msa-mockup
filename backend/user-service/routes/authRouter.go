package routes

import (
	controller "backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(routes *gin.Engine) {
	routes.POST("/users/register", controller.Register())
	routes.POST("/users/login", controller.Login())
	routes.POST("/users/logout", controller.Logout())
	routes.POST("/request-password-reset", controller.RequestPasswordReset())
	routes.POST("/reset-password", controller.ResetPassword())
}
