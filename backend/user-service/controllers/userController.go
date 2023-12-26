package controllers

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"backend/database"

	helper "backend/helpers"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()
var SECRET_KEY string = os.Getenv("SECRET_KEY")
var blacklistedPasswords map[string]bool

// Function to generate a random token (used for password recovery)
func generateRandomToken() string {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", token)
}

// SendEmail sends an email with the given subject and body to the specified recipient
func SendEmail(to, subject, body string) error {
	// Set up authentication information
	auth := smtp.PlainAuth("", "agustina21@ethereal.email", "JKfDUgkC9tNCUZSK9u", "smtp.ethereal.email")

	// Compose the email message
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	// Connect to the server, authenticate, and send the email
	err := smtp.SendMail("smtp.ethereal.email:587", auth, "your-email@example.com", []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

// Endpoint for initiating password reset (password recovery)
func RequestPasswordReset() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user's email from the request
		email := c.PostForm("email")

		// Check if the email exists in the database
		var user models.User
		err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Generate a unique token
		token := generateRandomToken()

		// Save the token and timestamp in the database
		resetToken := models.PasswordResetToken{
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour * 2), // Token is valid for 2 hours (adjust as needed)
		}

		update := bson.M{
			"$set": bson.M{"PasswordResetToken": resetToken},
		}

		_, err = userCollection.UpdateOne(
			context.TODO(),
			bson.M{"email": email},
			update,
		)

		if err != nil {
			log.Println("Error updating user with password reset token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate password reset"})
			return
		}

		// Email to the user with a link containing the token
		emailSubject := "Password Reset Instructions"
		emailBody := fmt.Sprintf("Click the following link to reset your password: https://your-app.com/reset?token=%s", token)

		err = SendEmail(email, emailSubject, emailBody)
		if err != nil {
			log.Println("Error sending email:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Password reset initiated. Check your email for instructions."})
	}
}

// Endpoint for resetting the password
func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user's email, token, and new password from the request
		email := c.PostForm("email")
		token := c.PostForm("token")
		newPassword := c.PostForm("newPassword")

		// Retrieve the user from the database by email
		var user models.User
		err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Check if the user has a valid password reset token
		if user.PasswordResetToken == nil || user.PasswordResetToken.Token != token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Check if the token is still valid (not expired)
		if time.Now().After(user.PasswordResetToken.ExpiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			return
		}

		// Update the user's password in the database
		hashedPassword := HashPassword(newPassword)
		update := bson.M{
			"$set": bson.M{
				"password":             hashedPassword,
				"password_reset_token": nil, // Clear the password reset token after use
			},
		}

		_, err = userCollection.UpdateOne(
			context.TODO(),
			bson.M{"email": email},
			update,
		)

		if err != nil {
			log.Println("Error updating user password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
	}
}

func init() {
	// Read blacklisted passwords from the external file
	file, err := os.Open("controllers/blacklistedPasswords.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	blacklistedPasswords = make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		blacklistedPasswords[scanner.Text()] = true
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func isBlacklistedPassword(password string) bool {
	return blacklistedPasswords[password]
}

// HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

// CreateUser is the api used to tget a single user
func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		l := log.New(gin.DefaultWriter, "User controller: ", log.LstdFlags)
		l.Println(c.GetString("Authorization"))

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		if isBlacklistedPassword(*user.Password) {
			c.JSON(http.StatusOK, gin.H{"message": "The chosen password is blacklisted"})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

// Login is the api used to tget a single user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Request headers:", c.Errors)
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.JSON(http.StatusBadRequest, gin.H{"error": "TEST LOL"})
			c.JSON(http.StatusBadRequest, gin.H{"error": ""})
			fmt.Println("I am here")
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
			c.JSON(http.StatusBadRequest, gin.H{"error": "TEST LOL"})
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		passwordIsValid, _ := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			// c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			c.JSON(http.StatusBadRequest, gin.H{"error": "TEST LOL"})
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect password"})

			return

		}

		if foundUser.Email == nil {
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			c.JSON(http.StatusBadRequest, gin.H{"error": "TEST LOL"})
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "TEST LOL"})
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)

	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.New(gin.DefaultWriter, "User controller: ", log.LstdFlags)
		// Get the refresh token
		refreshToken := helper.ExtractRefreshToken(c)
		l.Println("primljen je: ", c.GetString("Authorization"))
		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "TEST LOL"})
		} else {

			c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
		}
	}
}

// GetUser is the api used to tget a single user
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)

	}

}
func GetUsers(c *gin.Context) {
	//sluzi za logovanje u konzolu, korisno
	l := log.New(gin.DefaultWriter, "User Controller ", log.LstdFlags)
	// Retrieve the JWT token from the Authorization header
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
		l.Println("Parsing token..")
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
	l.Println("Extract the claims from the parsed token")
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		l.Println("Token invalid")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalid"})
		return
	}

	// Decode the token without verifying the signature
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		// Handle error
		l.Println("Error decoding token without verification:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "error decoding token"})
		return
	}

	// Log the token claims
	l.Println("Token claims:", parsedToken.Claims)

	// Retrieve the user ID from the token claims
	l.Println("Retrieving user id..")
	userID, ok := claims["Uid"].(string) //Check which claims you're selecting
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in token"})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user models.User
	userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "user": user})
}

func BecomeHost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Get the user ID from the request parameters
		userID := c.Param("user_id")

		// Verify that the user making the request has the necessary permissions
		if err := helper.MatchUserTypeToUid(c, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			cancel()
			return
		}

		// Update the user's user_type to "Host" in the database
		update := bson.M{"$set": bson.M{"user_type": "Host"}}
		result, err := userCollection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user type"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "User became a host"})
	}
}
func EditProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		type EditUser struct {
			First_name *string   `json:"first_name" validate:min=2,max=100"`
			Last_name  *string   `json:"last_name" validate:"min=2,max=100"`
			Email      *string   `json:"email" validate:"email"`
			Phone      *string   `json:"phone" validate:"min=6"`
			Address    *string   `json:"address" validate:"min=3"`
			Token      *string   `json:"token"`
			User_type  *string   `json:"user_type" validate:"oneof=Guest Host UUser"`
			Updated_at time.Time `json:"updated_at"`
		}
		//	recieves user data and a token
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user EditUser
		var foundUser models.User
		tokenEmail := c.MustGet("email").(string)

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"bind error": err.Error()})
			defer cancel()
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": tokenEmail}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"not found error": err.Error()})
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			defer cancel()
			c.JSON(http.StatusBadRequest, gin.H{"struct validation error": validationErr.Error()})
			return
		}

		if *foundUser.Email != *user.Email {
			count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
			defer cancel()
			if err != nil {
				log.Panic(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
				return
			}
			if count > 0 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
				return
			}
		}

		if foundUser.Phone != user.Phone {

			count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
			defer cancel()
			if err != nil {
				log.Panic(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
				return
			}

			if count > 0 {
				defer cancel()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "this phone number already exists"})
				return
			}
		}
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filter := bson.M{"user_id": foundUser.User_id}

		update := bson.M{
			"$set": bson.M{
				"first_name": user.First_name,
				"last_name":  user.Last_name,
				"email":      user.Email,
				"phone":      user.Phone,
				"address":    user.Address,
				"user_type":  user.User_type,
				"updated_at": user.Updated_at,
			},
		}
		updateOptions := options.Update().SetUpsert(false)
		_, err = userCollection.UpdateOne(context.TODO(), filter, update, updateOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		defer cancel()
		c.JSON(http.StatusOK, "Profile updated successfully")

	}
}
func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Passwords struct {
			Old_password *string `json:"old_password"`
			New_password *string `json:"new_password" validate:"min=12"`
		}
		l := log.New(gin.DefaultWriter, "User controller ChangePassword() : ", log.LstdFlags)
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		tokenEmail := c.MustGet("email").(string)

		var passwords Passwords
		var foundUser models.User

		if err := c.BindJSON(&passwords); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}
		hashedOldPassword := ""
		hashedNewPassword := ""
		if passwords.Old_password != nil && passwords.New_password != nil {
			hashedOldPassword = HashPassword(*passwords.Old_password)
			hashedNewPassword = HashPassword(*passwords.New_password)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
		}

		err := userCollection.FindOne(ctx, bson.M{"email": tokenEmail}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		l.Println("unhashed OldPassword: ", *passwords.Old_password)
		l.Println("comparing passwords, founduser: ", *foundUser.Password)
		l.Println("comparing passwords, hashedOldPassword ", hashedOldPassword)
		passwordIsValid, _ := VerifyPassword(*passwords.Old_password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect password"})

			return
		}

		if !strings.ContainsAny(*passwords.New_password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			defer cancel()
			c.JSON(http.StatusBadRequest, gin.H{"password validation error": "Password needs to contain at least one Capital letter"})
			return
		}
		specialCharRegex := regexp.MustCompile(`[!@#$%^&*()-_=+{};:'",.<>?/\\|[\]~]`)
		if !specialCharRegex.MatchString(*passwords.New_password) {
			defer cancel()
			c.JSON(http.StatusBadRequest, gin.H{"password validation error": "Password needs to contain at least one special character"})
			return
		}
		update := bson.M{
			"$set": bson.M{
				"password": hashedNewPassword,
			},
		}

		_, err = userCollection.UpdateOne(
			context.TODO(),
			bson.M{"email": tokenEmail},
			update,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		defer cancel()
		c.JSON(http.StatusOK, "Password changed successfully")
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userID := c.Param("user_id")
		result, err := userCollection.DeleteOne(ctx, bson.M{"user_id": userID})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
