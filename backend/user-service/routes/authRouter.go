package routes

import (
	controller "backend/controllers"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
<<<<<<< HEAD
func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/register", controller.SignUp())
	incomingRoutes.POST("/users/login", controller.Login())
=======
func AuthRoutes(routes *gin.Engine) {
	routes.POST("/users/register", controller.Register())
	routes.POST("/users/login", controller.Login())
	routes.POST("/users/logout", controller.Logout())
>>>>>>> main
}
