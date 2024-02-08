package middleware

import (
	"fmt"

	"github.com/alexvancasper/TunnelBroker/web/pkg/common/db"

	"net/http"
	"os"
	"time"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NotRequireAuth(c *gin.Context) {
	// Get the cookie off the request
	tokenString, err := c.Cookie("Authorization")

	if err == nil {
		// c.AbortWithStatusJSON(http.StatusTemporaryRedirect, gin.H{"message": "already registered"})
		// c.Abort()
		fmt.Printf("error, already authorized %s\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/user/")
		c.Abort()
		return
	}

	// Decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err == nil || token != nil {
		fmt.Printf("error, already authorized %s\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/user/")
		c.Abort()
		return
	}
	c.Next()

}

func RequireAuth(c *gin.Context) {
	// Get the cookie off the request
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		fmt.Printf("error, not authorized %s\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}

	// Decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || token == nil {
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Chec k the expiry date
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Printf("error, not authorized %s\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
		}

		// Find the user with token Subject
		var user models.User
		db.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			fmt.Printf("error, not authorized %s\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
		}

		// Attach the request
		c.Set("user", user)

		//Continue
		c.Next()
	} else {
		fmt.Printf("error, not authorized %s\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}
}
