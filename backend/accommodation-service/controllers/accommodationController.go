// accommodationController.go

package controllers

import (
	"context"
	"net/http"

	"accommodation-service/database" // Import your database package
	"accommodation-service/models"   // Import your models package
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var accommodationCollection *mongo.Collection = database.OpenCollection(database.Client, "accommodation")

func CreateAccommodation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accommodation models.Accommodation
		if err := c.BindJSON(&accommodation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			return
		}

		// Assuming you have the authenticated user ID available
		// Accommodation owner can be set accordingly
		// accommodation.Owner = userID

		// Insert the accommodation into the database
		insertResult, err := accommodationCollection.InsertOne(context.Background(), &accommodation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create accommodation"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Accommodation created", "id": insertResult.InsertedID})
	}
}

func GetAccommodation() gin.HandlerFunc {
	return func(c *gin.Context) {
		accommodationID := c.Param("accommodation_id")

		var accommodation models.Accommodation
		err := accommodationCollection.FindOne(context.Background(), bson.M{"_id": accommodationID}).Decode(&accommodation)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Accommodation not found"})
			return
		}

		c.JSON(http.StatusOK, accommodation)
	}
}
