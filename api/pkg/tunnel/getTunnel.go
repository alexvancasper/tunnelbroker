package tunnels

import (
	"net/http"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTunnel(c *gin.Context) {
	id := c.Param("id")
	api := c.Param("api")

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

	c.JSON(http.StatusOK, &tunnel)
}
