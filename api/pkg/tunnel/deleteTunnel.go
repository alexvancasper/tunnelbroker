package tunnels

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/alexvancasper/TunnelBroker/web/internal/broker"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h handler) DeleteTunnel(c *gin.Context) {
	l := h.Logf.WithFields(logrus.Fields{
		"function": "DeleteTunnel",
	})
	id := c.Param("id")
	api := c.Param("api")

	var user models.User

	if apiExist := h.DB.Where("api = ?", api).First(&user); apiExist.Error != nil {
		l.Errorf("Not able to find provided API key %s", apiExist.Error)
		// c.AbortWithError(http.StatusNotFound, apiExist.Error)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": apiExist.Error})
		return
	}

	var tunnel models.Tunnel

	if result := h.DB.Where("user_id = ?", user.ID).First(&tunnel, id); result.Error != nil {
		l.Errorf("Not able to find a user with provided tunnel id %s", result.Error)
		// c.AbortWithError(http.StatusNotFound, result.Error)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": result.Error})
		return
	}

	data, err := tunnel.Marshal()
	if err != nil {
		l.Errorf("Error of marshalling tunnel data %s", err)
		// c.AbortWithError(http.StatusInternalServerError, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	go func() {
		err = h.Msg.PublishMsg(data, broker.DELETE)
		if err != nil {
			l.Errorf("Error of sending tunnel data to server %s", err)
			// c.AbortWithError(http.StatusInternalServerError, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	}()

	l.Debugf("Tunnel data %+v", tunnel)
	h.DB.Delete(&tunnel)
	l.Infof("Tunnel deleted from DB with id %d", tunnel.ID)

	newTunnelCount := user.TunnelCount - 1
	l.Debugf("user tunnel count decreased %d->%d", user.TunnelCount, newTunnelCount)
	h.DB.Model(&user).Update("tunnel_count", newTunnelCount)

	networkAddr := GetNetworkAddr(tunnel.IPv6ClientEndpoint, h.Logf)
	releaseIPv6Prefixes(networkAddr, h.Logf)
	l.Infof("IPv6Endpoint network released %s", networkAddr)
	releaseIPv6Prefixes(tunnel.PD, h.Logf)
	l.Infof("IPv6 PD released %s", tunnel.PD)

	c.Status(http.StatusOK)
}

func releaseIPv6Prefixes(prefix string, logf *logrus.Logger) {
	l := logf.WithFields(logrus.Fields{
		"function": "releaseIPv6Prefixes",
	})
	prefixData := strings.Split(prefix, "/")
	requestURL := fmt.Sprintf("http://%s/release?prefix=%s&prefixlen=%s", os.Getenv("IPAM"), prefixData[0], prefixData[1])
	l.Debugf("requestURL %s", requestURL)
	req, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	if err != nil {
		l.Errorf("client: could not create request: %s\n", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Errorf("client: error making http request: %s\n", err)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		l.Errorf("client: could not read response body: %s\n", err)
	}
	l.Debugf("client: response body: %s\n", string(resBody))
}
