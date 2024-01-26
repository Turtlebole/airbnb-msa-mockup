// accommodationController.go

package controllers

import (
	"accommodation-service/repositories"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func GetAllAccommodationsByHostLocal(hostID primitive.ObjectID) ([]models.Accommodation, error) {

	filter := bson.M{"host_id": hostID}

	cursor, err := accommodationCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var accommodations []models.Accommodation
	for cursor.Next(context.Background()) {
		var accommodation models.Accommodation
		if err := cursor.Decode(&accommodation); err != nil {
			return nil, err
		}
		accommodations = append(accommodations, accommodation)
	}
	if len(accommodations) == 0 {
		return nil, fmt.Errorf("no accommodations found")
	}

	return accommodations, nil
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
		accID := c.Param("_id")
		l := log.New(gin.DefaultWriter, "Accommodation controller: ", log.LstdFlags)
		reservationURL := fmt.Sprintf("http://reservation-service:8002/reservations/by_room/%s", accID)

		reservationResponse, err := http.Get(reservationURL)
		if err != nil {
			l.Printf("Error making GET request for reservations of accommodation %s: %v\n", accID, err)
			return
		}
		defer reservationResponse.Body.Close()

		reservationBody, err := io.ReadAll(reservationResponse.Body)
		if err != nil {
			l.Printf("Error reading reservations response body for accommodation %s: %v\n", accID, err)
			return
		}
		if reservationResponse.StatusCode == http.StatusOK {
			l.Printf("Reservations for accommodation %s: %s\n", accID, string(reservationBody))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete an accommodation that has active reservations"})
			return
		} else if reservationResponse.StatusCode == http.StatusNotFound {
			l.Printf("No reservations for accommodation %s\n", accID)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": reservationResponse})
		}

		accommodationRepo, err := repositories.New(context.Background(), l)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
			return
		}

		if err := accommodationRepo.Delete(accID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "Accommodation deleted successfully")
	}
}

func DeleteAccommodationLocal(accid string) (bool, error) {
	l := log.New(gin.DefaultWriter, "Accommodation controller: ", log.LstdFlags)
	reservationURL := fmt.Sprintf("http://reservation-service:8002/reservations/by_room/%s", accid)

	reservationResponse, err := http.Get(reservationURL)
	if err != nil {
		l.Printf("Error making GET request for reservations of accommodation %s: %v\n", accid, err)
		return false, err
	}
	defer reservationResponse.Body.Close()

	reservationBody, err := io.ReadAll(reservationResponse.Body)
	if err != nil {
		l.Printf("Error reading reservations response body for accommodation %s: %v\n", accid, err)
		return false, err
	}
	if reservationResponse.StatusCode == http.StatusOK {
		l.Printf("Reservations for accommodation %s: %s\n", accid, string(reservationBody))
		return true, fmt.Errorf("reservations have been found for accommodation: %s", accid)
	} else if reservationResponse.StatusCode == http.StatusNotFound {
		l.Printf("No reservations for accommodation %s: %s\n", accid)
	} else {
		return false, fmt.Errorf("error: %v", reservationResponse)
	}
	accommodationRepo, err := repositories.New(context.Background(), l)
	if err != nil {
		l.Printf("Failed to connect to database")
		return false, err
	}

	if err := accommodationRepo.Delete(accid); err != nil {
		return false, err
	} else {
		return false, nil
	}
}

func DeleteAccommodationsByHost() gin.HandlerFunc {
	return func(c *gin.Context) {
		hostID, err := primitive.ObjectIDFromHex(c.Param("_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		accs, err := GetAllAccommodationsByHostLocal(hostID)
		if err != nil {
			if err.Error() == "no accommodations found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "no accommodations found for host"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}

		accsDeleted := 0
		for _, accommodation := range accs {
			hasActiveReservations, err := DeleteAccommodationLocal(accommodation.Id.Hex())
			if hasActiveReservations {
				c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete if there are active reservations present for one or more accommodations"})
				return
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			accsDeleted = accsDeleted + 1
		}

		c.JSON(http.StatusOK, gin.H{"Success": fmt.Sprintf("%d accommodation(s) deleted successfully", accsDeleted)})
	}
}

func parseAccommodationsResponse(body []byte) []models.Accommodation {
	var accommodations []models.Accommodation
	err := json.Unmarshal(body, &accommodations)
	if err != nil {
		fmt.Println("Error parsing accommodations response:", err)
		return nil
	}
	return accommodations
}
