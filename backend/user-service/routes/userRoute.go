package routes

import (
	"12factorapp/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(routes *gin.Engine) {
	routes.GET("/users", controllers.GetUser())
	routes.GET("/users/:user_id", controllers.GetUser())
}
