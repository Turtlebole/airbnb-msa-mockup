package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"profile-service/database"
	"profile-service/repositories"
)

var profileCollection *mongo.Collection = database.OpenCollection(database.Client, "profiles")

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
		profileID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
			return
		}

		// Get profile by ID
		profile, err := repositories.GetProfileByID(client, profileID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}

		// Return profile data
		c.JSON(http.StatusOK, profile)
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
