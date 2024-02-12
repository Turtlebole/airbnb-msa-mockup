package controllers

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"review-service/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var reviewCollection *mongo.Collection

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
func GetAverageRatingHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		accommodationID := c.Param("accommodationID")

		if _, err := primitive.ObjectIDFromHex(accommodationID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid accommodation ID"})
			return
		}

		// Call the repository function to get the average rating
		reviewRepo := repositories.NewReviewRepo(client)
		averageRating, err := reviewRepo.GetAverageRatingByAccommodation(context.Background(), accommodationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch average rating"})
			return
		}

		// Return the average rating in the response
		c.JSON(http.StatusOK, gin.H{"averageRating": averageRating})
	}
}
func GetReviewsByUserHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")

		// Fetch reviews by user ID
		reviewRepo := repositories.NewReviewRepo(client)
		reviews, err := reviewRepo.GetReviewsByUser(c.Request.Context(), userID)
		if err != nil {
			fmt.Println("Error fetching user reviews:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("Fetched User Reviews:", reviews)
		c.JSON(http.StatusOK, reviews)
	}
}

func UpdateReviewHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updatedReview repositories.Review

		if err := c.ShouldBindJSON(&updatedReview); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Extract review ID from the path parameters
		reviewID := c.Param("reviewID")
		if _, err := primitive.ObjectIDFromHex(reviewID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
			return
		}

		log.Printf("Received updated review: %+v", updatedReview)
		log.Printf("Review ID: %s", reviewID)

		// Call the repository function to update the review
		reviewRepo := repositories.NewReviewRepo(client)
		err := reviewRepo.UpdateReview(c.Request.Context(), reviewID, &updatedReview)
		if err != nil {
			log.Printf("Failed to update review. Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
			return
		}

		log.Println("Review updated successfully")

		c.Status(http.StatusOK)
	}
}
func DeleteReviewHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract review ID from the path parameters
		reviewID := c.Param("reviewID")
		if _, err := primitive.ObjectIDFromHex(reviewID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
			return
		}

		// Call the repository function to delete the review
		reviewRepo := repositories.NewReviewRepo(client)
		err := reviewRepo.DeleteReview(c.Request.Context(), reviewID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
			return
		}

		c.Status(http.StatusOK)
	}
}
func CreateHostReviewHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newReview repositories.Review

		if err := c.ShouldBindJSON(&newReview); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ensure that both UserID and HostID are set
		if newReview.UserID.IsZero() || newReview.HostID.IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserID and HostID are required"})
			return
		}

		fmt.Printf("Received Host Review Data: %+v\n", newReview)

		reviewRepo := repositories.NewReviewRepo(client)
		err := reviewRepo.CreateHostReview(c.Request.Context(), &newReview)
		if err != nil {
			fmt.Printf("Error Creating Host Review: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	}
}
func GetReviewsByUserAndHostsHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")

		// Fetch reviews by user ID and include host information
		reviewRepo := repositories.NewReviewRepo(client)
		reviews, err := reviewRepo.GetReviewsByUserAndHosts(c.Request.Context(), userID)
		if err != nil {
			fmt.Println("Error fetching user reviews:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("Fetched User Reviews:", reviews)
		c.JSON(http.StatusOK, reviews)
	}
}
func GetHostReviewsHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		hostID := c.Param("hostID")

		// Call the repository function to get host reviews
		reviewRepo := repositories.NewReviewRepo(client)
		reviews, err := reviewRepo.GetHostReviews(c.Request.Context(), hostID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch host reviews"})
			return
		}

		// Return the host reviews in the response
		c.JSON(http.StatusOK, reviews)
	}
}
func GetReviewsByAccommodationAndHostHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		accommodationID := c.Param("accommodationID")
		hostID := c.Param("hostID")

		reviews, err := repositories.GetReviewsByAccommodationAndHost(client, accommodationID, hostID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(200, gin.H{"reviews": reviews})
	}
}
func GetAllReviewsHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewRepo := repositories.NewReviewRepo(client)
		reviews, err := reviewRepo.GetAllReviews(c.Request.Context())
		if err != nil {
			fmt.Println("Error fetching all reviews:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("Fetched All Reviews:", reviews)
		c.JSON(http.StatusOK, reviews)
	}
}
func GetReviewByIDHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract review ID from the path parameters
		reviewID := c.Param("reviewID")
		if _, err := primitive.ObjectIDFromHex(reviewID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
			return
		}

		// Call the repository function to get the review by ID
		reviewRepo := repositories.NewReviewRepo(client)
		review, err := reviewRepo.GetReviewByID(c.Request.Context(), reviewID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
			return
		}

		// Return the review data
		c.JSON(http.StatusOK, review)
	}
}
