// accommodationController.go

package controllers

import (
	"accommodation-service/repositories"
	"context"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"accommodation-service/database" // Import your database package
	"accommodation-service/models"   // Import your models package

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var accommodationCollection *mongo.Collection = database.OpenCollection(database.Client, "accommodations")

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

func UpdateAccommodation() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("_id")
		var accommodation models.Accommodation

		if err := c.ShouldBindJSON(&accommodation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logger := log.New(os.Stdout, "ACCOMMODATION-LOG: ", log.LstdFlags)
		accommodationRepo, err := repositories.New(context.Background(), logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
			return
		}

		defer accommodationRepo.Disconnect(context.Background())

		// Convert string ID to ObjectID
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		// Assign the ObjectID to the accommodation ID field
		accommodation.Id = objID

		// Attempt accommodation update
		err = accommodationRepo.Update(id, &accommodation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update accommodation"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Accommodation updated successfully"})
	}
}

func GetAccommodation() gin.HandlerFunc {
	return func(c *gin.Context) {
		accommodationID := c.Param("_id")

		var accommodation models.Accommodation

		// Convert string ID to ObjectID
		objID, err := primitive.ObjectIDFromHex(accommodationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		// Find accommodation by ObjectID
		err = accommodationCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&accommodation)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Accommodation not found"})
			return
		}

		c.JSON(http.StatusOK, accommodation)
	}
}

func GetAllAccommodationsByHost() gin.HandlerFunc {
	return func(c *gin.Context) {
		hostID, err := primitive.ObjectIDFromHex(c.Param("host_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID format"})
			return
		}

		filter := bson.M{"host_id": hostID}

		cursor, err := accommodationCollection.Find(context.Background(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying accommodations"})
			return
		}
		defer cursor.Close(context.Background())

		var accommodations []models.Accommodation
		for cursor.Next(context.Background()) {
			var accommodation models.Accommodation
			if err := cursor.Decode(&accommodation); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding accommodations"})
				return
			}
			accommodations = append(accommodations, accommodation)
		}
		if len(accommodations) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No accommodations found for host"})
			return
		}

		c.JSON(http.StatusOK, accommodations)
	}
}

func GetAllAccommodations() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accommodations []models.Accommodation

		cursor, err := accommodationCollection.Find(context.Background(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accommodations"})
			return
		}
		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &accommodations); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode accommodations"})
			return
		}

		c.JSON(http.StatusOK, accommodations)
	}
}

func DeleteAccommodation() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
