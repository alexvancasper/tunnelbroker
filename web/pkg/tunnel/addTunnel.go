package tunnels

import (
	"net/http"
	"os"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
)

type AddTunnelRequestBody struct {
	// IPv4Remote - client's ipv4 address
	IPv4Remote string `json:"ipv4remote"`
	// TunnelName - Name of the tunnel
	TunnelName string `json:"tunnelname"`
}

func (h handler) AddTunnel(c *gin.Context) {
	api := c.Param("api")
	body := AddTunnelRequestBody{}

	// получаем тело запроса
	if err := c.Bind(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var user models.User
	if apiExist := h.DB.Where("api = ?", api).First(&user); apiExist.Error != nil {
		c.AbortWithError(http.StatusNotFound, apiExist.Error)
		return
	}

	var tunnel models.Tunnel

	tunnel.TunnelName = body.TunnelName
	tunnel.IPv4Remote = body.IPv4Remote
	tunnel.UserID = user.ID
	tunnel.IPv4Local = getLocalIPv4()
	tunnel.PD = "pd prefix"
	tunnel.P2P = "p2p prefix"
	tunnel.Configured = false

	if result := h.DB.Create(&tunnel); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	if apiExist := h.DB.Where("api = ?", api).First(&user); apiExist.Error != nil {
		c.AbortWithError(http.StatusNotFound, apiExist.Error)
		return
	}

	// c.JSON(http.StatusCreated, &tunnel)
	c.JSON(http.StatusCreated, gin.H{"message": "tunnel is created"})
}

func getLocalIPv4() string {
	addr := os.Getenv("IPv4LOCALADDR")
	if len(addr) == 0 {
		return "185.185.58.180"
	}
	return addr
}

func getPD() string {
	return "pd/64 prefix"
}

func getP2P() string {
	return "p2p/127 prefix"
}
