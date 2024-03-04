package users

import (
	"fmt"

	"github.com/alexvancasper/TunnelBroker/web/internal/middleware"
	"github.com/alexvancasper/TunnelBroker/web/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type handler struct {
	DB   *gorm.DB
	Logf *logrus.Logger
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB, logf *logrus.Logger) {
	h := &handler{
		DB:   db,
		Logf: logf,
	}

	routes := r.Group("/user", middleware.RequireAuth)
	routes.GET("/", h.GetUser)
	routes.DELETE("/", h.DeleteUser)
}

func getIDfromToken(c *gin.Context) (uint, error) {
	user, exists := c.Get("user")
	if !exists {
		return 0, fmt.Errorf("user not found")
	}
	var u models.User
	var ok bool
	if u, ok = user.(models.User); !ok {
		return 0, fmt.Errorf("user assertion failed")
	}
	return u.ID, nil
}
