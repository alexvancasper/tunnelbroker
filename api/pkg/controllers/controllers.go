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

func Signup(c *gin.Context) {
	title := "TunnelBroker 6in4 - Register new user"
	// Get the email/pass off req Body
	var body struct {
		Email    string `json:"login"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error": "Failed to read body",
		// })
		c.HTML(http.StatusBadRequest, "adduser.html", gin.H{
			"Title": title,
			"Error": "Failed to read request",
		})

		return
	}
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error": "Failed to hash password.",
		// })
		c.HTML(http.StatusBadRequest, "adduser.html", gin.H{
			"Title": title,
			"Error": "Failed to read data",
		})
		return
	}

	// Create the user
	user := models.User{Login: body.Email, Password: string(hash), API: generateAPI()}

	result := db.DB.Create(&user)

	if result.Error != nil {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error": "Failed to create user.",
		// })
		c.HTML(http.StatusBadRequest, "adduser.html", gin.H{
			"Title": title,
			"Error": "Failed to read user data",
		})
		c.Abort()
		return
	}

	// Respond
	// c.JSON(http.StatusCreated, gin.H{"mesage": "User created"})
	c.HTML(http.StatusCreated, "adduser.html", gin.H{
		"Title": title,
		"Error": "User created. Redirect to login page",
	})
	// c.Redirect(http.StatusCreated, "/login")
}

func Login(c *gin.Context) {
	title := "TunnelBroker 6in4"
	// Get email & pass off req body
	var body struct {
		Email    string `json:"login"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error": "Failed to read body",
		// })
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Title": title,
			"Error": "Failed to read request",
		})
		return
	}
	// Look up for requested user
	var user models.User

	db.DB.First(&user, "login = ?", body.Email)

	if user.ID == 0 {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error": "Invalid email or password",
		// })
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Title": title,
			"Error": "Invalid email or password",
		})
		return
	}

	// Compare sent in password with saved users password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error": "Invalid email or password",
		// })
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Title": title,
			"Error": "Invalid email or password",
		})
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error": "Failed to create token",
		// })
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Title": title,
			"Error": "Failed to create token",
		})
		return
	}

	// Respond
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*1, "", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "logged in successfully"})
	// c.Redirect(http.StatusTemporaryRedirect, "/user/")
}

func Logout(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil || len(tokenString) <= 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cookie is not found"})
	}
	c.SetCookie("Authorization", tokenString, 1, "", "", false, false)
	// c.JSON(http.StatusOK, gin.H{"message": "user logout"})
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	// user.(models.User).Email    -->   to access specific data

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func generateAPI() string {
	uuid := uuid.NewV1().String()
	uuid = strings.ReplaceAll(uuid, "-", "")
	return uuid
}
