package tunnels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddTunnelRequestBody struct {
	// IPv4Remote - client's ipv4 address
	IPv4Remote string `json:"ipv4remote"`
	// TunnelName - Name of the tunnel
	TunnelName string `json:"tunnelname"`
}

func (h handler) AddTunnel(c *gin.Context) {
	l := h.Logf.WithFields(logrus.Fields{
		"function": "AddTunnel",
	})
	api := c.Param("api")
	body := AddTunnelRequestBody{}

	// получаем тело запроса
	if err := c.Bind(&body); err != nil {
		l.Errorf("Not able to bind POST form data err %s", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var user models.User
	if apiExist := h.DB.Where("api = ?", api).First(&user); apiExist.Error != nil {
		l.Errorf("Not able to find user with provided API key %s", api)
		c.AbortWithError(http.StatusNotFound, apiExist.Error)
		return
	}

	var tunnel models.Tunnel

	tunnel.TunnelName = body.TunnelName
	tunnel.IPv4Remote = body.IPv4Remote
	tunnel.UserID = user.ID
	tunnel.IPv4Local = getLocalIPv4(h.Logf)
	tunnel.PD = getPrefix(64, h.Logf)
	p2pPrefix := getPrefix(127, h.Logf)
	tunnel.IPv6ClientEndpoint, tunnel.IPv6ServerEndpoint = GetEndpoints(p2pPrefix, h.Logf)
	tunnel.Configured = false

	if result := h.DB.Create(&tunnel); result.Error != nil {
		l.Errorf("Not able to insert tunnel in to DB %s", result.Error)
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	if apiExist := h.DB.Where("api = ?", api).First(&user); apiExist.Error != nil {
		l.Errorf("Not able to find user with provided API (2) key %s", api)
		c.AbortWithError(http.StatusNotFound, apiExist.Error)
		return
	}
	l.Debugf("Tunnel data: %+v", tunnel)
	l.Infof("Tunnel created tunnel_id %d", tunnel.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "tunnel is created"})
}

func getLocalIPv4(logf *logrus.Logger) string {
	l := logf.WithFields(logrus.Fields{
		"function": "getPrefix",
	})
	addr := os.Getenv("IPv4LOCALADDR")
	if len(addr) == 0 {
		l.Errorf("Environment variable IPv4LOCALADDR is empty, got default value")
		return "185.185.58.180"
	}
	return addr
}

func getPrefix(prefixlen int, logf *logrus.Logger) string {
	l := logf.WithFields(logrus.Fields{
		"function": "getPrefix",
	})
	requestURL := fmt.Sprintf("http://%s/acquire?prefixlen=%d", os.Getenv("IPAM"), prefixlen)
	l.Debugf("requestURL %s", requestURL)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		l.Errorf("client: could not create request: %s", err)
		return ""
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Errorf("client: error making http request: %s", err)
		return ""
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		l.Errorf("client: could not read response body: %s", err)
		return ""
	}
	l.Debugf("client: response body: %s", resBody)
	var Prefix struct {
		Prefix string `json:"prefix"`
	}
	json.Unmarshal(resBody, &Prefix)
	return Prefix.Prefix
}
