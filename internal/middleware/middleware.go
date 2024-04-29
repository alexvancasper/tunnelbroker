package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexvancasper/TunnelBroker/web/internal/common/db"
	"github.com/alexvancasper/TunnelBroker/web/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var signTokenFunc = func(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	return []byte(os.Getenv("SECRET")), nil
}

func NotRequireAuth(c *gin.Context) {
	// Get the cookie off the request
	tokenString, err := c.Cookie("Authorization")

	if errors.Is(err, http.ErrNoCookie) {
		c.Next()
	}

	// Decode/validate it
	token, err := jwt.Parse(tokenString, signTokenFunc)
	if err == nil || token != nil {
		fmt.Printf("error, already authorized\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func RequireAuth(c *gin.Context) {
	// Get the cookie off the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		fmt.Printf("error, not authorized %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}

	// Decode/validate it
	token, err := jwt.Parse(tokenString, signTokenFunc)

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	var claims jwt.MapClaims
	var ok bool

	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		fmt.Printf("error, not authorized\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	expiryAt, err := claims.GetExpirationTime()
	if err != nil {
		fmt.Printf("error, not authorized %s\n", "token issue, not able to get expiration time")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	if expiryAt.Before(time.Now()) {
		fmt.Print("token expired\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Find the user with token Subject
	var user models.User
	db.DB.First(&user, claims["sub"])

	if user.ID == 0 {
		fmt.Printf("error, not authorized\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
	c.Set("user", user)

	c.Next()
}
