package users

import (
	"net/http"
	"strings"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type AddUserRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h handler) AddUser(c *gin.Context) {
	body := AddUserRequestBody{}

	// получаем тело запроса
	if err := c.Bind(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var user models.User

	user.Login = body.Login
	user.Password = body.Password
	user.API = generateAPI()

	if result := h.DB.Create(&user); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	c.JSON(http.StatusCreated, &user)
}

func generateAPI() string {
	uuid := uuid.NewV1().String()
	uuid = strings.ReplaceAll(uuid, "-", "")
	return uuid
}
