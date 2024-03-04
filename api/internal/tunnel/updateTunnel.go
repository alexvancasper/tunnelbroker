package tunnels

import (
	"net/http"

	"github.com/alexvancasper/TunnelBroker/web/internal/models"
	"github.com/gin-gonic/gin"
)

type UpdateTunnelRequestBody struct {
	// TunnelName- Name of the tunnel
	TunnelName string `json:"tunnelname"`
	// IPv4Remote- Client's ipv4 address
	IPv4Remote string `json:"ipv4remote"`
}

func (h handler) UpdateTunnel(c *gin.Context) {
	id := c.Param("id")
	api := c.Param("api")
	body := UpdateTunnelRequestBody{}

	// получаем тело запроса
	if err := c.BindJSON(&body); err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})

		return
	}

	var user models.User

	if apiExist := h.DB.Where("api = ?", api).First(&user); apiExist.Error != nil {
		// c.AbortWithError(http.StatusNotFound, apiExist.Error)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": apiExist.Error})

		return
	}

	var tunnel models.Tunnel

	if result := h.DB.Where("user_id = ?", user.ID).First(&tunnel, id); result.Error != nil {
		// c.AbortWithError(http.StatusNotFound, result.Error)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": result.Error})
		return
	}

	tunnel.TunnelName = body.TunnelName
	tunnel.IPv4Remote = body.IPv4Remote

	h.DB.Save(&tunnel)

	c.JSON(http.StatusOK, &tunnel)
}
