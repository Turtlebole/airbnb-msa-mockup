package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Accommodation struct {
	Id                  primitive.ObjectID `bson:"_idomitempty" json:"id"`
	Owner               User               `bson:"owner,omitempty" json:"owner"`
	Name                string             `bson:"name" json:"name"`
	Location            string             `bson:"location" json:"location"`
	MinimalGuestsNumber int                `bson:"minGuestNum" json:"minimal_guests_number"`
	MaximalGuestsNumber int                `bson:"maxGuestNum" json:"maximal_guests_number"`
	Amenities           []string           `bson:"amenities,omitempty" json:"amenities"`
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
