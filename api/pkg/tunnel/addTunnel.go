package tunnels

import (
	"fmt"
	"net/http"

	"github.com/alexvancasper/TunnelBroker/web/internal/broker"
	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddTunnelRequestBody struct {
	// IPv4Remote - client's ipv4 address
	IPv4Remote string `json:"ipv4remote"`
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
		// c.AbortWithError(http.StatusBadRequest, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	var user models.User
	if apiExist := h.DB.Where("api = ?", api).First(&user); apiExist.Error != nil {
		l.Errorf("Not able to find user with provided API key %s", api)
		// c.AbortWithError(http.StatusNotFound, apiExist.Error)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": apiExist.Error})
		return
	}

	//Check for user tunnel limit
	if user.TunnelCount >= user.TunnelLimit {
		l.Debugf("You have reached tunnel limit %d", user.TunnelLimit)
		//c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("You have reached tunnel limit %d", user.TunnelLimit)})
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("You have reached tunnel limit %d", user.TunnelLimit)})
		return
	}

	var tunnel models.Tunnel

	if !IsIPv4Public(body.IPv4Remote) {
		l.WithField("input ipv4 address", body.IPv4Remote).Error("Not able to parse")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Wrong IPv4 address"})
		return
	}

	tunnel.TunnelName = generateName(fmt.Sprintf("%d%s", user.ID, body.IPv4Remote))
	tunnel.IPv4Remote = body.IPv4Remote
	tunnel.UserID = user.ID
	tunnel.IPv4Local = getLocalIPv4(h.Logf)
	tunnel.PD = getPrefix(64, h.Logf)
	p2pPrefix := getPrefix(127, h.Logf)
	tunnel.IPv6ClientEndpoint, tunnel.IPv6ServerEndpoint = GetEndpoints(p2pPrefix, h.Logf)
	tunnel.Configured = false

	if apiExist := h.DB.Where("api = ?", api).First(&user); apiExist.Error != nil {
		l.Errorf("Not able to find user with provided API (2) key %s", api)
		// c.AbortWithError(http.StatusNotFound, apiExist.Error)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": apiExist.Error})

		return
	}

	if result := h.DB.Create(&tunnel); result.Error != nil {
		l.Errorf("Not able to insert tunnel in to DB %s", result.Error)
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
		err = h.Msg.PublishMsg(data, broker.ADD)
		if err != nil {
			l.Errorf("Error of sending tunnel data to server %s", err)
			// c.AbortWithError(http.StatusInternalServerError, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	}()

	newTunnelCount := user.TunnelCount + 1
	l.Debugf("user tunnel count increased %d->%d", user.TunnelCount, newTunnelCount)
	h.DB.Model(&user).Update("tunnel_count", newTunnelCount)

	l.Debugf("Tunnel data: %+v", tunnel)
	l.Infof("Tunnel created tunnel_id %d", tunnel.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "tunnel is created"})
}
