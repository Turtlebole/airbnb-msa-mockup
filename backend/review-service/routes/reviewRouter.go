package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"review-service/controllers"
)

func ReviewRoutes(router *gin.Engine, client *mongo.Client) {
	router.POST("/reviews/create", controllers.CreateReviewHandler(client))
	router.GET("/reviews/accommodation/:accommodationID", controllers.GetReviewsByAccommodationHandler(client))
	router.GET("/reviews/average-rating/:accommodationID", controllers.GetAverageRatingHandler(client))
	router.PUT("/reviews/update/:reviewID", controllers.UpdateReviewHandler(client))
	router.DELETE("/reviews/delete/:reviewID", controllers.DeleteReviewHandler(client))
	router.GET("/reviews/user/:userID", controllers.GetReviewsByUserHandler(client))
	router.POST("/reviews/create-host", controllers.CreateHostReviewHandler(client))
	router.GET("/reviews/user/:userID/hosts", controllers.GetReviewsByUserAndHostsHandler(client))
	router.GET("/reviews/host/:hostID", controllers.GetHostReviewsHandler(client))
	router.GET("/reviews/accommodation/:accommodationID/hosts/:hostID", controllers.GetReviewsByAccommodationAndHostHandler(client))
	router.GET("/reviews/all", controllers.GetAllReviewsHandler(client))
	router.GET("/reviews/:reviewID", controllers.GetReviewByIDHandler(client))

}
