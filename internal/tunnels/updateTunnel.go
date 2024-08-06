package tunnels

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/alexvancasper/TunnelBroker/web/internal/models"
	"github.com/alexvancasper/broker"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UpdateTunnelRequestBody struct {
	// IPv4Remote- Client's ipv4 address
	IPv4Remote string `json:"ipv4remote"`
}

func (u UpdateTunnelRequestBody) Validate() error {
	netaddr := net.ParseIP(u.IPv4Remote)
	if netaddr == nil {
		return fmt.Errorf("bad ipv4 address")
	}
	return nil
}

func (h handler) UpdateTunnel(c *gin.Context) {
	l := h.Logf.WithFields(logrus.Fields{
		"function": "UpdateTunnel",
	})

	idC, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "id does not exists"})
		return
	}
	api, ok := c.Params.Get("api")
	if !ok || len(api) != 32 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "api does not exists"})
		return
	}
	id, err := strconv.Atoi(idC)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "id should be numeric"})
		return
	}
	l.Debugf("input data id: %d, api: %s", id, api)

	body := UpdateTunnelRequestBody{}

	// получаем тело запроса
	if err := c.ShouldBindJSON(&body); err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = body.Validate()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User

	if apiExist := h.DB.Where("api = ?", api).First(&user); apiExist.Error != nil {
		// c.AbortWithError(http.StatusNotFound, apiExist.Error)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": apiExist.Error})

		return
	}

	var tunnel models.Tunnel

	if result := h.DB.Where("user_id = ?", user.ID).First(&tunnel, id); result.Error != nil {
		// c.AbortWithError(http.StatusNotFound, result.Error)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": result.Error})
		return
	}

	tunnel.IPv4Remote = body.IPv4Remote

	h.DB.Save(&tunnel)

	data, err := tunnel.Marshal()
	if err != nil {
		l.Errorf("Error of marshalling tunnel data %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	// TODO: нам нужно быть уверены, что сообщение поместилось в MQ
	// иначе будет неконсистентность между БД и реальными туннелями.
	go func() {
		err = h.Msg.PublishMsg(data, broker.UPDATE)
		if err != nil {
			l.Errorf("Error of updating tunnel data %s", err)
			// c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	}()
	l.Debugf("Tunnel data: %+v", tunnel)
	l.Infof("Tunnel updated tunnel_id %d", tunnel.ID)

	c.JSON(http.StatusOK, &tunnel)
}
