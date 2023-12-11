package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasswordResetToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("passwordStrength", passwordStrengthValidation)
}

type User struct {
	ID                 primitive.ObjectID  `bson:"_id"`
	First_name         *string             `json:"first_name" validate:"required,min=2,max=100"`
	Last_name          *string             `json:"last_name" validate:"required,min=2,max=100"`
	Email              *string             `json:"email" validate:"email,required"`
	Password           *string             `json:"password" validate:"required,min=8,passwordStrength"`
	Phone              *string             `json:"phone" validate:"required"`
	Address            *string             `json:"address" validate:"required"`
	Token              *string             `json:"token"`
	User_type          *string             `json:"user_type" validate:"required,oneof=Guest Host UUser"`
	Refresh_token      *string             `json:"refresh_token"`
	Created_at         time.Time           `json:"created_at"`
	Updated_at         time.Time           `json:"updated_at"`
	User_id            string              `json:"user_id"`
	PasswordResetToken *PasswordResetToken `json:"password_reset_token,omitempty"`
}

func passwordStrengthValidation(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	// Check for at least one capital letter
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}

	// Check for at least one special character
	specialCharRegex := regexp.MustCompile(`[!@#$%^&*()-_=+{};:'",.<>?/\\|[\]~]`)
	if !specialCharRegex.MatchString(password) {
		return false
	}

	return true
}
