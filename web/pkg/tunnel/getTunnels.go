package tunnels

import (
	"net/http"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTunnels(c *gin.Context) {
	var tunnels []models.Tunnel

	if result := h.DB.Find(&tunnels); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusOK, &tunnels)
}
