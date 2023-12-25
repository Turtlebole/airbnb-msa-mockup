package controllers

import (
	"fmt"
	"net/http"
	"review-service/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var reviewCollection *mongo.Collection // Initialize this with the appropriate collection

func CreateReviewHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newReview repositories.Review

		if err := c.ShouldBindJSON(&newReview); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Received Review Data: %+v\n", newReview)

		reviewRepo := repositories.NewReviewRepo(client)
		err := reviewRepo.CreateReview(c.Request.Context(), &newReview)
		if err != nil {
			fmt.Printf("Error Creating Review: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func GetReviewsByAccommodationHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		accommodationID := c.Param("accommodationID")
		fmt.Printf("Received request for reviews with Accommodation ID: %s\n", accommodationID)

		// Fetch reviews by accommodation ID
		reviewRepo := repositories.NewReviewRepo(client)
		reviews, err := reviewRepo.GetReviewsByAccommodation(c.Request.Context(), accommodationID)
		if err != nil {
			fmt.Println("Error fetching reviews:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("Fetched Reviews:", reviews)
		c.JSON(http.StatusOK, reviews)
	}
}
