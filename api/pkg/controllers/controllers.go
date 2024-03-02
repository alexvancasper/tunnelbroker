package controllers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexvancasper/TunnelBroker/web/pkg/common/db"
	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func PostSignup(c *gin.Context) {
	var body struct {
		Email    string `json:"login"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read request"})
		c.Abort()
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read data"})
		c.Abort()
		return
	}
	user := models.User{Login: body.Email, Password: string(hash), API: generateAPI()}

	var userDB models.User
	if userExist := db.DB.Where("login = ?", body.Email).First(&userDB); userExist.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User already exist"})
		c.Abort()
		return
	}

	result := db.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User is created"})
}

func PostLogin(c *gin.Context) {
	var body struct {
		Email    string `json:"login"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read request"})
		c.Abort()
		return
	}
	var user models.User
	db.DB.First(&user, "login = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		c.Abort()
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		c.Abort()
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session"})
		c.Abort()
		return
	}

	// Respond
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*1, "", "", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "logged in successfully"})
}

func Logout(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil || len(tokenString) <= 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cookie is not found"})
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, -1, "", "", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "logout done"})
}

// func Validate(c *gin.Context) {
// 	user, _ := c.Get("user")
// 	// user.(models.User).Email    -->   to access specific data
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": user,
// 	})
// }

func generateAPI() string {
	uuid := uuid.NewV1().String()
	uuid = strings.ReplaceAll(uuid, "-", "")
	return uuid
}
