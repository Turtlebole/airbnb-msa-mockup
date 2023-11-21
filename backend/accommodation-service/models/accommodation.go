package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
