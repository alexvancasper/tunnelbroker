package tunnels

import (
	"net"
	"regexp"
	"strings"
)

func getIPv4address(address string) string {
	ipv4Regexp := regexp.MustCompile(`(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
	address = strings.Trim(address, " ")
	return ipv4Regexp.FindString(address)
}

func isPublicIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch {
		case ip4[0] == 0 && ip4[1] == 0 && ip4[2] == 0 && ip4[3] == 0:
			return false
		case ip4[0] == 255:
			return false
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func IsIPv4Public(address string) bool {
	ipv4addr := getIPv4address(address)
	netAddr := net.ParseIP(ipv4addr)
	if netAddr == nil {
		return false
	}
	if isPublicIP(netAddr) {
		return true
	}
	return false
}
