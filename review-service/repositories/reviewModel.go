package repositories

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Review struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	UserName        string             `bson:"user_name,omitempty" json:"user_name"`
	AccommodationID primitive.ObjectID `bson:"accommodation_id,omitempty" json:"accommodation_id"`
	Rating          int                `json:"rating"`
	Comment         string             `json:"comment"`
}

func (r *Review) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r)
}

func (r *Review) FromJSON(reader io.Reader) error {
	d := json.NewDecoder(reader)
	return d.Decode(r)
}

type User struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `json:"name"`
}
