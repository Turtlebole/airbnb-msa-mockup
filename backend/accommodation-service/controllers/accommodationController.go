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
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"accommodation-service/database"
	"accommodation-service/models"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var accommodationCollection *mongo.Collection = database.OpenCollection(database.Client, "accommodations")

func isThereAnyOverlap(start1, end1, start2, end2 string) (bool, error) {
	const layout = "2006-01-02"

	start1Time, err := time.Parse(layout, start1)
	if err != nil {
		return false, err
	}
	end1Time, err := time.Parse(layout, end1)
	if err != nil {
		return false, err
	}
	start2Time, err := time.Parse(layout, start2)
	if err != nil {
		return false, err
	}
	end2Time, err := time.Parse(layout, end2)
	if err != nil {
		return false, err
	}

	if end1Time.Before(start2Time) || start1Time.After(end2Time) {
		return false, nil
	}
	return true, nil
}

func isAWithinB(start1, end1, start2, end2 string) (bool, error) {
	const layout = "2006-01-02"

	s1T, err := time.Parse(layout, start1)
	if err != nil {
		return false, err
	}
	e1T, err := time.Parse(layout, end1)
	if err != nil {
		return false, err
	}
	s2T, err := time.Parse(layout, start2)
	if err != nil {
		return false, err
	}
	e2T, err := time.Parse(layout, end2)
	if err != nil {
		return false, err
	}

	if s2T.Before(s1T) && e1T.Before(e2T) {
		return true, nil
	}
	return false, nil
}

func CreateAccommodation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accommodation models.Accommodation
		if err := c.BindJSON(&accommodation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		intervals := accommodation.Availability

		for i := 0; i < len(intervals)-1; i++ {
			for j := i + 1; j < len(intervals); j++ {
				overlap, err := isThereAnyOverlap(intervals[i].Start, intervals[i].End, intervals[j].Start, intervals[j].End)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				if overlap {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Availability intervals cannot be overlapping."})
					return
				}
			}
		}

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

		l := log.New(gin.DefaultWriter, "Accommodation controller: ", log.LstdFlags)
		// Convert string ID to ObjectID
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		var newAccommodation models.Accommodation
		if err := c.ShouldBindJSON(&newAccommodation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetch existing accommodation from MongoDB
		var existingAccommodation models.Accommodation
		err = accommodationCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&existingAccommodation)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Accommodation not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing accommodation"})
			}
			return
		}
		availabilityUpdated := !reflect.DeepEqual(existingAccommodation.Availability, newAccommodation.Availability)

		if availabilityUpdated {

			reservations, err := GetReservations(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			intervals := newAccommodation.Availability

			for i := 0; i < len(reservations)-1; i++ {
				overlapFound := false
				for j := i + 1; j < len(intervals); j++ {
					overlap, err := isAWithinB(reservations[i].CheckInDate, reservations[i].CheckOutDate, intervals[j].Start, intervals[j].End)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					if overlap {
						l.Printf("intervals overlap")
						overlapFound = true
						break
					} else {
						continue
					}
				}
				if !overlapFound {
					c.JSON(http.StatusBadRequest, gin.H{"error": "you cannot edit accommodataion availability for dates that are reserved"})
					return
				}
			}

			for i := 0; i < len(intervals)-1; i++ {
				for j := i + 1; j < len(intervals); j++ {
					overlap, err := isThereAnyOverlap(intervals[i].Start, intervals[i].End, intervals[j].Start, intervals[j].End)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					if overlap {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Availability intervals cannot overlap."})
						return
					}
				}
			}
		}

		// Update the accommodation in MongoDB
		updateResult, err := accommodationCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": objID},
			bson.M{"$set": newAccommodation},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update accommodation"})
			return
		}

		// Check if the document was successfully updated
		if updateResult.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Accommodation not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":             "Accommodation updated successfully",
			"availabilityUpdated": availabilityUpdated,
		})
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
func IsAvailable(accommodation models.Accommodation, checkIn, checkOut time.Time) bool {
	layout := "2006-01-02"
	for _, interval := range accommodation.Availability {
		intervalStart, _ := time.Parse(layout, interval.Start)
		intervalEnd, _ := time.Parse(layout, interval.End)
		if (checkIn.After(intervalStart) || checkIn.Equal(intervalStart)) && checkOut.Before(intervalEnd) {
			return true
		}
	}
	return false
}

func IsAvailableNow(accommodation models.Accommodation) bool {
	now := time.Now()
	layout := "2006-01-02"
	for _, interval := range accommodation.Availability {
		intervalStart, _ := time.Parse(layout, interval.Start)
		intervalEnd, _ := time.Parse(layout, interval.End)
		if (now.After(intervalStart) || now.Equal(intervalStart)) && now.Before(intervalEnd) {
			return true
		}
	}
	return false
}
func GetAllAvailableAccommodations() gin.HandlerFunc {
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

		for i, accommodation := range accommodations {
			if !IsAvailableNow(accommodation) {
				accommodations = append(accommodations[:i], accommodations[i+1:]...)
			}
		}

		c.JSON(http.StatusOK, accommodations)
	}
}
func HasReservations(accID string) (bool, error) {
	l := log.New(gin.DefaultWriter, "Accommodation controller: ", log.LstdFlags)

	reservationURL := fmt.Sprintf("http://reservation-service:8002/reservations/by_room/%s", accID)

	reservationResponse, err := http.Get(reservationURL)
	if err != nil {
		l.Printf("Error making GET request for reservations of accommodation %s: %v\n", accID, err)
		return false, fmt.Errorf("error making GET request for reservations of accommodation")
	}
	defer reservationResponse.Body.Close()

	reservationBody, err := io.ReadAll(reservationResponse.Body)
	if err != nil {
		l.Printf("Error reading reservations response body for accommodation %s: %v\n", accID, err)
		return false, fmt.Errorf("error reading reservations response body for accommodation")
	}
	if reservationResponse.StatusCode == http.StatusOK {
		l.Printf("Reservations found for accommodation %s: %s\n", accID, string(reservationBody))
		return true, nil
	} else if reservationResponse.StatusCode == http.StatusNotFound {
		l.Printf("No reservations found for accommodation %s\n", accID)
		return false, nil
	} else {
		l.Println("error: ", string(reservationBody))
		return false, fmt.Errorf(string(reservationBody))
	}
}

type ReservationByRoom struct {
	RoomId          string     `json:"room_id"`
	ReservationId   gocql.UUID `json:"reservation_id"`
	GuestID         string     `json:"guest_id"`
	GuestUsername   string     `json:"guest_username"`
	NumberOfGuests  int        `json:"number_of_guests"`
	Price           float64    `json:"price"`
	ReservationDate string     `json:"reservation_date"`
	CheckInDate     string     `json:"checkin_date"`
	CheckOutDate    string     `json:"checkout_date"`
}

func GetReservations(accID string) ([]ReservationByRoom, error) {
	l := log.New(gin.DefaultWriter, "Accommodation controller: ", log.LstdFlags)

	reservationURL := fmt.Sprintf("http://reservation-service:8002/reservations/by_room/%s", accID)

	reservationResponse, err := http.Get(reservationURL)
	if err != nil {
		l.Printf("Error making GET request for reservations of accommodation %s: %v\n", accID, err)
		return nil, fmt.Errorf("error making GET request for reservations of accommodation")
	}
	defer reservationResponse.Body.Close()

	reservationBody, err := io.ReadAll(reservationResponse.Body)
	if err != nil {
		l.Printf("Error reading reservations response body for accommodation %s: %v\n", accID, err)
		return nil, fmt.Errorf("error reading reservations response body for accommodation")
	}

	var reservations []ReservationByRoom
	if reservationResponse.StatusCode == http.StatusOK {
		err = json.Unmarshal(reservationBody, &reservations)
		if err != nil {
			l.Printf("Error unmarshalling reservations response body for accommodation %s: %v\n", accID, err)
			return nil, fmt.Errorf("error unmarshalling reservations response body")
		}
		l.Printf("%v Reservations found for accommodation %s\n", len(reservations), accID)
		return reservations, nil
	} else if reservationResponse.StatusCode == http.StatusNotFound {
		l.Printf("No reservations found for accommodation %s\n", accID)
		return nil, nil
	} else {
		l.Println("error: ", string(reservationBody))
		return nil, fmt.Errorf(string(reservationBody))
	}
}

func DeleteAccommodation() gin.HandlerFunc {
	return func(c *gin.Context) {
		accID := c.Param("_id")
		l := log.New(gin.DefaultWriter, "Accommodation controller: ", log.LstdFlags)

		accommodationRepo, err := repositories.New(context.Background(), l)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
			return
		}
		if hasReservations, err := HasReservations(accID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else if hasReservations {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you cannot delete an accommodation that has active reservations present"})
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
		l.Printf("No reservations for accommodation %s\n", accid)
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

func GetAvailableDates() gin.HandlerFunc {
	return func(c *gin.Context) {
		accommodationID := c.Param("_id")

		var accommodation models.Accommodation

		objID, err := primitive.ObjectIDFromHex(accommodationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		err = accommodationCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&accommodation)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Accommodation not found"})
			return
		}

		c.JSON(http.StatusOK, accommodation.Availability)
	}
}
