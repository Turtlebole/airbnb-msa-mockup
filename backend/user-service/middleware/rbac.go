package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mikespook/gorbac"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func Rbac(c *gin.Context, address string) {

	rbac := gorbac.New()
	rUUser := gorbac.NewStdRole("uuser")
	rHost := gorbac.NewStdRole("host")
	rGuest := gorbac.NewStdRole("guest")

	pCreateAccomodation := gorbac.NewStdPermission("accommodationsCreate")

	rHost.Assign(pCreateAccomodation)

	rbac.Add(rHost)
	rbac.Add(rUUser)
	rbac.Add(rGuest)

	rbac.SetParent("uuser", "host")
	rbac.SetParent("uuser", "guest")

	ur := getUserRole(c)

	if ur == "" {
		c.JSON(http.StatusUnauthorized, "Error")
		c.Next()
	}
	// if rbac.IsGranted(ur, address) {

	// }

	c.Next()

}

func getUserRole(c *gin.Context) string {

	authHeader := c.Request.Header["Authorization"]

	if len(authHeader) == 0 {
		c.JSON(http.StatusUnauthorized, "No header")
		return ""
	}
	authString := strings.Join(authHeader, "")
	tokenString := strings.Split(authString, "Bearer ")[1]

	// Check that the token string is not empty
	if len(tokenString) == 0 {
		c.JSON(http.StatusUnauthorized, "Token empty")
		return ""
	}
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return ""
	}
	// Handle token parsing errors

	claims := token.Claims.(jwt.MapClaims)

	userRole := claims["User_type"].(string)
	return userRole

}
