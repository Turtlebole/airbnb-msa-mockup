package controllers

import (
	"backend/database"
	"backend/models"
	"backend/services"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var userRepo services.UsersRepo

func init() {
	userRepo = services.NewMongoUsersRepo(database.Client, "dbConnection", "user")

}

func Login(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	foundUser, err := userRepo.GetUser(context.Background(), user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, foundUser)
}

func SignUp(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if EmailChecker(context.Background(), *user.Email) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This email already exists"})
		return
	}

	user.Created_On = time.Now()
	user.Updated_On = time.Now()

	if _, err := userRepo.RegisterUser(context.Background(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
		return
	}

	c.JSON(http.StatusOK, "User created successfully")
}

// Check if an email is already in use
func EmailChecker(ctx context.Context, email string) bool {
	count, err := userRepo.CountByEmail(ctx, email)
	if err != nil {
		return true
	}
	return count > 0
}
