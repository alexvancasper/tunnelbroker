package tunnels

import (
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

	routes := r.Group("/tunnel/:api")
	routes.POST("/", h.AddTunnel)
	routes.GET("/:id", h.GetTunnel)
	routes.PUT("/:id", h.UpdateTunnel)
	routes.DELETE("/:id", h.DeleteTunnel)
	// routes.GET("/", h.GetTunnels)

}
