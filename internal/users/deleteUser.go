package users

import (
	"net/http"

	"github.com/alexvancasper/TunnelBroker/web/internal/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteUser(c *gin.Context) {
	userID, err := getIDfromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "delete token issue"})
		return
	}

	var user models.User

	if result := h.DB.First(&user, userID); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": result.Error})
		return
	}

	h.DB.Delete(&user)

	c.Status(http.StatusOK)
}
