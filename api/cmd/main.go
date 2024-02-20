package main

import (
	"github.com/alexvancasper/TunnelBroker/web/internal/broker"
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

	// Initialize DB connection
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")
	h := db.Init(dbUrl)

	// Initialize message broker connection
	m, err := broker.MsgBrokerInit(os.Getenv("BROKER_CONN"), os.Getenv("QUEUENAME"))
	if err != nil {
		MyLogger.Errorf("Message broker error init: %s", err)
	}
	defer m.Close()

	store := cookie.NewStore([]byte(os.Getenv("COOKIEKEY1")))
	option := csrf.Options{
		Secret: os.Getenv("COOKIEKEY2"),
		ErrorFunc: func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSRF token mismatch"})
			c.Abort()
		},
	}

	r := gin.Default()
	r.Use(sessions.Sessions("session", store))
	r.Use(csrf.Middleware(option))

	r.StaticFS("/static", http.Dir("./pkg/webview/static"))
	r.LoadHTMLGlob("./pkg/webview/templates/*")
	r.GET("/", webview.Index)
	r.GET("/login", middleware.NotRequireAuth, webview.Login)
	r.GET("/signup", middleware.NotRequireAuth, webview.Register)
	r.GET("/logout", controllers.Logout)
	r.GET("/ip", webview.IP)

	r.POST("/signup", controllers.PostSignup)
	r.POST("/login", controllers.PostLogin)

	users.RegisterRoutes(r, h, MyLogger)
	tunnels.RegisterRoutes(r, h, m, MyLogger)

	r.Run(port)
}
