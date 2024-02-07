package main

import (
	"os"
	"time"

	"github.com/alexvancasper/TunnelBroker/web/pkg/common/db"
	"github.com/alexvancasper/TunnelBroker/web/pkg/controllers"
	"github.com/alexvancasper/TunnelBroker/web/pkg/middleware"
	tunnels "github.com/alexvancasper/TunnelBroker/web/pkg/tunnel"
	"github.com/alexvancasper/TunnelBroker/web/pkg/users"
	"github.com/alexvancasper/TunnelBroker/web/pkg/webview"
	formatter "github.com/fabienm/go-logrus-formatters"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	//Initialize Logging connections
	var MyLogger = logrus.New()

	gelfFmt := formatter.NewGelf("API")
	MyLogger.SetFormatter(gelfFmt)
	MyLogger.SetOutput(os.Stdout)
	loglevel, err := logrus.ParseLevel("debug")
	if err != nil {
		MyLogger.WithField("function", "main").Fatalf("error %v", err)
	}
	MyLogger.SetLevel(loglevel)

	viper.SetConfigFile("/pkg/common/envs/.env")
	viper.ReadInConfig()

	// port := viper.Get("PORT").(string)
	// dbUrl := viper.Get("DB_URL").(string)
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	r := gin.Default()
	h := db.Init(dbUrl)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:8080", "http://127.0.0.1:8000", "http://127.0.0.1:3000", "http://localhost:3000"}
	config.AllowHeaders = []string{"content-type", "content-lenght", "authorization", "origin", "Set-Cookie"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowCredentials = true
	config.MaxAge = 1 * time.Minute
	r.Use(cors.New(config))
	// r.Static("/static", "/pkg/webview/static")
	// r.LoadHTMLGlob("/pkg/webview/templates/*")
	r.GET("/", webview.Index)
	// r.GET("/login", webview.Login)
	r.GET("/signup", middleware.NotRequireAuth, webview.Register)
	r.GET("/logout", controllers.Logout)
	// r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	users.RegisterRoutes(r, h, MyLogger)
	tunnels.RegisterRoutes(r, h, MyLogger)

	r.Run(port)
}

func CORSMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
