package tunnels

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteTunnel(c *gin.Context) {
	id := c.Param("id")
	api := c.Param("api")

	var user models.User

	if apiExist := h.DB.Where("api = ?", api).First(&user); apiExist.Error != nil {
		c.AbortWithError(http.StatusNotFound, apiExist.Error)
		return
	}

	var tunnel models.Tunnel

	if result := h.DB.Where("user_id = ?", user.ID).First(&tunnel, id); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	h.DB.Delete(&tunnel)
	releaseIPv6Prefixes(tunnel.P2P)
	releaseIPv6Prefixes(tunnel.PD)

	c.Status(http.StatusOK)
}

func releaseIPv6Prefixes(prefix string) error {
	prefix = prefix[1 : len(prefix)-2]
	prefixData := strings.Split(prefix, "/")
	requestURL := fmt.Sprintf("http://%s/release?prefix=%s&prefixlen=%s", os.Getenv("IPAM"), prefixData[0], prefixData[1])

	req, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return err
	}
	fmt.Printf("client: response body: %s\n", string(resBody))
	fmt.Printf("IPV6 %s is released\n", prefix)
	return nil
}
