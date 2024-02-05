package users

import (
	"net/http"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetUser(c *gin.Context) {
	id, err := getIDfromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var user models.User

	if result := h.DB.Model(&models.User{}).Preload("Tunnels").First(&user, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	// c.HTML(http.StatusOK, "user.html", gin.H{
	// 	"Title": "User room",
	// 	"User":  user,
	// })
	c.JSON(http.StatusOK, &user)
}