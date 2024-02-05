package webview

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "Tunnel Broker 6in4",
	})
}

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "adduser.html", gin.H{
		"Title": "Tunnel Broker register",
	})
}

func Index(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}
