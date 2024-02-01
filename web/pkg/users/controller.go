package users

import (
	"fmt"

	"github.com/alexvancasper/TunnelBroker/web/pkg/middleware"
	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	routes := r.Group("/user", middleware.RequireAuth)
	routes.GET("/", h.GetUser)
	routes.PUT("/", h.UpdateUser)
	routes.DELETE("/", h.DeleteUser)
	// routes.GET("/", h.GetUsers)
	// routes.POST("/", h.AddUser)

}

func getIDfromToken(c *gin.Context) (uint, error) {

	user, exists := c.Get("user")
	if !exists {
		return 0, fmt.Errorf("user not found")
	}
	return user.(models.User).ID, nil
}
