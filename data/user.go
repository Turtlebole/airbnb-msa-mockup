package data

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"regexp"
)

// Defining the main struct for our API
type User struct {
	ID         int     `json:"id"`                       //specifies that in the incoming Body the field to map to this will be called "id"
	First_Name *string `json:"name" validate:"required"` //there are some integrated validation, for eg. this specifies that a value for name must be provided, otherwise it will not be valid
	Last_Name  *string `json:"description"`
	email      *string `json:"price" validate:"gt=0"`
	address    *string `json:"address"`
	SKU        string  `json:"sku" validate:"required,sku"` //the tag "sku" is there so we can add custom validation
	Created_On string  `json:"created_On"`
	Updated_On string  `json:"updated_On"`
	Deleted_On string  `json:"deleted_On"`
}

type Users []*User

// Function to validate the incoming object from front.
// NOTE: if the tag "sku" is not present in the struct anotations we will get an error
func (u *User) Validate() error {
	validate := validator.New()

	err := validate.RegisterValidation("sku", validateSKU)
	if err != nil {
		return err
	}

	return validate.Struct(u)
}

// We use a regex to validate a custom look of sku-s
// For eg: abc-abc-abc
func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("[a-z]+-[a-z]+-[a-z]")
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

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
