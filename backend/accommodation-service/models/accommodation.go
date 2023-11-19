package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Accommodation struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `json:"name"`
	Location      string             `json:"location"`
	Amenities     []string           `json:"amenities"`
	MinGuests     int                `json:"min_guests"`
	MaxGuests     int                `json:"max_guests"`
	Images        []string           `json:"images"`
	Availability  string             `json:"availability"`
	PricePerNight float64            `json:"price_per_night"`
}

var accommodations []Accommodation

func CreateAccommodation(w http.ResponseWriter, r *http.Request) {
	var newAccommodation Accommodation

	err := json.NewDecoder(r.Body).Decode(&newAccommodation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Store the new accommodation in your data structure or database
	accommodations = append(accommodations, newAccommodation)

	// Respond with a success message or the created accommodation
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAccommodation)
}
