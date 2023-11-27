package middleware

import (
	"fmt"
	"net/http"
	"strings"

	helper "backend/helpers"

	"github.com/gin-gonic/gin"
)

// Validates token
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		//clientToken := c.Request.Header.Get("token")
		clientToken := c.Request.Header.Get("Authorization")
		clientToken = strings.Replace(clientToken, "Bearer ", "", 1)
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		//c.JSON(http.StatusInternalServerError, gin.H{"errorlol": err})
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"errorlol": err})
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
