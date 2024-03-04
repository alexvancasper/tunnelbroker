package tunnels

import (
	"github.com/alexvancasper/broker"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type handler struct {
	DB   *gorm.DB
	Logf *logrus.Logger
	Msg  *broker.MsgBroker
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB, m *broker.MsgBroker, logf *logrus.Logger) {
	h := &handler{
		DB:   db,
		Logf: logf,
		Msg:  m,
	}

	routes := r.Group("/tunnel/:api")
	routes.POST("/", h.AddTunnel)
	routes.GET("/:id", h.GetTunnel)
	// routes.PUT("/:id", h.UpdateTunnel)
	routes.DELETE("/:id", h.DeleteTunnel)
	// routes.GET("/", h.GetTunnels)
}
