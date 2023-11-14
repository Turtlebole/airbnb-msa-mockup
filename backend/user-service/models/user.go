package models

import (
	"encoding/json"
	"io"
	"time"
)

// Defining the main struct for our API
type User struct {
	ID         int       `json:"id"`                       //specifies that in the incoming Body the field to map to this will be called "id"
	First_Name *string   `json:"name" validate:"required"` //there are some integrated validation, for eg. this specifies that a value for name must be provided, otherwise it will not be valid
	Last_Name  *string   `json:"description"`
	Email      *string   `json:"price" validate:"gt=0"`
	Address    *string   `json:"address"`
	Created_On time.Time `json:"created_On"`
	Updated_On time.Time `json:"updated_On"`
	Deleted_On time.Time `json:"deleted_On"`
	UserType   UserType  `json:"user_type"`
}

type UserType string

const (
	GuestUser        UserType = "guest"
	HostUser         UserType = "host"
	UnidentifiedUser UserType = "uuser"
)

type Users []*User

// Functions to encode and decode products to json and from json.
// If we inject an interface we achieve dependancy injection, meaning that anything that implements this interface can be passed down
// For us it will be ResponseWriter, but it also may be a file writer or something similar.
func (u *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (p *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
