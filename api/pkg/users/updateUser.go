package users

import (
	"net/http"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
)

type UpdateUserRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h handler) UpdateUser(c *gin.Context) {
	// id := c.Param("id")
	id, err := getIDfromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	body := UpdateUserRequestBody{}

	// получаем тело запроса
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var user models.User

	if result := h.DB.First(&user, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	user.Login = body.Login
	user.Password = body.Password

	h.DB.Save(&user)

	c.JSON(http.StatusOK, &user)
}
