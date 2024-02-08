package webview

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "TunnelBroker 6in4",
	})
}

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "adduser.html", gin.H{
		"Title": "TunnelBroker register",
	})
}

func Index(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}
