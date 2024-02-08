package main

import (
	"net/http"
	"os"

	"github.com/alexvancasper/TunnelBroker/web/pkg/common/db"
	"github.com/alexvancasper/TunnelBroker/web/pkg/controllers"
	"github.com/alexvancasper/TunnelBroker/web/pkg/middleware"
	tunnels "github.com/alexvancasper/TunnelBroker/web/pkg/tunnel"
	"github.com/alexvancasper/TunnelBroker/web/pkg/users"
	"github.com/alexvancasper/TunnelBroker/web/pkg/webview"
	formatter "github.com/fabienm/go-logrus-formatters"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	csrf "github.com/utrack/gin-csrf"
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

	// viper.SetConfigFile("./pkg/common/envs/.env")
	// viper.ReadInConfig()

	// port := viper.Get("PORT").(string)
	// dbUrl := viper.Get("DB_URL").(string)
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	r := gin.Default()
	h := db.Init(dbUrl)

	store := cookie.NewStore([]byte("secret"))
	option := csrf.Options{
		Secret: "secret123",
		ErrorFunc: func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSRF token mismatch"})
			c.Abort()
		},
	}

	r.Use(sessions.Sessions("session", store))
	r.Use(csrf.Middleware(option))

	// r.StaticFS("/static", http.Dir("/pkg/webview/static"))
	r.Static("/static", "./pkg/webview/static")
	r.LoadHTMLGlob("./pkg/webview/templates/*")
	r.GET("/", webview.Index)
	r.GET("/login", webview.Login)
	r.GET("/signup", middleware.NotRequireAuth, webview.Register)
	r.GET("/logout", controllers.Logout)
	r.GET("/ip", webview.IP)
	// r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	users.RegisterRoutes(r, h, MyLogger)
	tunnels.RegisterRoutes(r, h, MyLogger)

	r.Run(port)
}

func csrfCheck() {

}
