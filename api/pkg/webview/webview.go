package webview

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "TunnelBroker 6in4",
		"Token": csrf.GetToken(c),
	})
}

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "adduser.html", gin.H{
		"Title": "TunnelBroker register",
		"Token": csrf.GetToken(c),
	})
}

func Index(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func IP(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"origin": c.RemoteIP()})
	c.Abort()
}
