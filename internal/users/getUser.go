package users

import (
	"net/http"

	"github.com/alexvancasper/TunnelBroker/web/internal/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (h handler) GetUser(c *gin.Context) {
	id, err := getIDfromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var user models.User

	if result := h.DB.Model(&models.User{}).Preload("Tunnels").First(&user, id); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Token": csrf.GetToken(c),
		"User":  user,
	})
}
