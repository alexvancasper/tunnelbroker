package webview

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"token": csrf.GetToken(c),
	})
}

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"token": csrf.GetToken(c),
	})
}

func Index(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func IP(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"origin":    c.RemoteIP(),
		"x-real-ip": c.GetHeader("X-REAL-IP"),
	})
	c.Abort()
}
