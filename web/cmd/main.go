package main

import (
	"net/http"

	"github.com/alexvancasper/TunnelBroker/web/pkg/common/db"
	"github.com/alexvancasper/TunnelBroker/web/pkg/controllers"
	"github.com/alexvancasper/TunnelBroker/web/pkg/middleware"
	tunnels "github.com/alexvancasper/TunnelBroker/web/pkg/tunnel"
	"github.com/alexvancasper/TunnelBroker/web/pkg/users"
	"github.com/alexvancasper/TunnelBroker/web/pkg/webview"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	port := viper.Get("PORT").(string)
	dbUrl := viper.Get("DB_URL").(string)

	r := gin.Default()
	h := db.Init(dbUrl)

	r.Static("/static", "./pkg/webview/static")
	r.LoadHTMLGlob("pkg/webview/templates/*")
	r.GET("/", index)
	r.GET("/login", webview.Login)
	r.GET("/signup", middleware.NotRequireAuth, webview.Register)
	r.GET("/logout", controllers.Logout)
	// r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	users.RegisterRoutes(r, h)
	tunnels.RegisterRoutes(r, h)

	r.Run(port)
}

func index(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}
