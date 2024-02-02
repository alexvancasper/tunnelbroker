package tunnels

import (
	"fmt"
	"io"
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
	tunnel.PD = getPrefix(64)
	tunnel.P2P = getPrefix(127)
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

func getPrefix(prefixlen int) string {
	requestURL := fmt.Sprintf("http://%s/acquire?prefixlen=%d", os.Getenv("IPAM"), prefixlen)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return ""
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return ""
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return ""
	}
	fmt.Printf("client: response body: %s\n", resBody)
	return string(resBody)
}
