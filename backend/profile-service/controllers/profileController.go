package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"profile-service/database"
	"profile-service/repositories"
)

var profileCollection *mongo.Collection = database.OpenCollection(database.Client, "profiles")
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func CreateProfileHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProfile repositories.Profile
		if err := c.ShouldBindJSON(&newProfile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repositories.CreateProfile(client, &newProfile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func GetAllProfilesHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		profiles, err := repositories.GetAllProfiles(client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, profiles)
	}
}

func GetProfileByIDHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract profile ID from URL parameters
		userID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
			return
		}

		// Get profile by ID
		profile, err := repositories.GetProfileByID(client, userID.Hex()) // Convert ObjectID to string using Hex()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}

		// Return profile data
		c.JSON(http.StatusOK, profile)
	}
}

func GetIdFromToken(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header["Authorization"]

		if len(authHeader) == 0 {
			c.JSON(http.StatusUnauthorized, "No header")
			return
		}
		authString := strings.Join(authHeader, "")
		tokenString := strings.Split(authString, "Bearer ")[1]

		// Check that the token string is not empty
		if len(tokenString) == 0 {
			c.JSON(http.StatusUnauthorized, "Token empty")
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(SECRET_KEY), nil
		})

		// Handle token parsing errors
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Extract the claims from the parsed token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalid"})
			return
		}

		userID, ok := claims["Uid"].(string) //Check which claims you're selecting
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	}
}

func UpdateProfileHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract profile ID from URL parameters
		profileID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
			return
		}

		// Decode the JSON request body
		var updatedProfile repositories.Profile
		if err := c.ShouldBindJSON(&updatedProfile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the profile
		if err := repositories.UpdateProfile(client, profileID, &updatedProfile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}

func DeleteProfileHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract profile ID from URL parameters
		profileID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
			return
		}

		// Delete the profile
		if err := repositories.DeleteProfile(client, profileID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}

func GetProfileByEmailHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract email from URL parameters or request body
		email := c.Query("email")
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter missing"})
			return
		}

		// Get profile by email
		profile, err := repositories.GetProfileByEmail(client, email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}

		// Return profile data
		c.JSON(http.StatusOK, profile)
	}
}
