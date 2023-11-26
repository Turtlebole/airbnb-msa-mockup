package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	helper "backend/helpers"

	"github.com/gin-gonic/gin"
)

// Validates token
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.New(gin.DefaultWriter, "Authenticator ", log.LstdFlags)
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			authHeader := c.GetHeader("Authorization")
			clientToken = strings.Split(authHeader, "Bearer ")[1]
			if clientToken == "" {
				l.Println("No Authorization header provided")
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
				c.Abort()
				return
			}
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			l.Println("Validating token failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.User_type)

		c.Next()

	}
}
