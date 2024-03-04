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
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var user models.User

	if result := h.DB.Model(&models.User{}).Preload("Tunnels").First(&user, id); result.Error != nil {
		// c.AbortWithError(http.StatusNotFound, result.Error)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": result.Error})
		return
	}

	ip := c.GetHeader("X-Real-Ip")
	c.HTML(http.StatusOK, "user.html", gin.H{
		"Title":    "User room",
		"Token":    csrf.GetToken(c),
		"User":     user,
		"ClientIP": ip,
	})
}
