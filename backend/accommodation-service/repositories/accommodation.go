package repositories

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Accommodation struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner         User               `bson:"owner,omitempty" json:"owner"`
	Name          string             `json:"name"`
	Location      string             `json:"location"`
	Amenities     []string           `json:"amenities"`
	MinGuests     int                `json:"min_guests"`
	MaxGuests     int                `json:"max_guests"`
	Images        []string           `json:"images"`
	Availability  string             `json:"availability"`
	PricePerNight float64            `json:"price_per_night"`
}
type Accommodations []*Accommodation

func (a Accommodations) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Accommodation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

type User struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

func (u User) Equals(user User) bool {
	return u.Id.String() == user.Id.String()
}

func (a Accommodation) Of(user User) bool {
	return a.Owner.Equals(user)
}
