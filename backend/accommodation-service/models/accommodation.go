package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accommodation struct {
	Id           primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name         string                 `json:"name"`
	Location     string                 `json:"location"`
	Amenities    []string               `json:"amenities"`
	MinGuests    int                    `json:"min_guests"`
	MaxGuests    int                    `json:"max_guests"`
	Images       []string               `json:"images"`
	Availability []AvailabilityInterval `json:"availability"`
	PriceType    string                 `json:"price_type" validate:"oneof=Whole Per_Guest"`
	HostID       primitive.ObjectID     `bson:"host_id,omitempty" json:"host_id"`
}
type AvailabilityInterval struct {
	Start           string  `json:"start"`
	End             string  `json:"end"`
	PricePerNight   float64 `json:"price_per_night"`
	PriceOnWeekends float64 `json:"price_on_weekends"`
}
