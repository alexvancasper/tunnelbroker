package tunnels

import (
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"

	"github.com/seancfoley/ipaddress-go/ipaddr"
)

func GetEndpoints(ipv6str string) (string, string) {
	val := strings.Split(ipv6str, "/")
	ipv6mask := mask2int(val[1])
	ip6 := ip2byte(val[0])
	mask := ipv6MyMask(ipv6mask)
	firstIP := make([]byte, 16)
	for i := 0; i < len(ip6); i++ {
		firstIP[i] = ip6[i] & mask[i]
	}
	return fmt.Sprintf("%s/%s", formatIPv6(firstIP), val[1]), fmt.Sprintf("%s/%s", formatIPv6(incv6IP(firstIP)), val[1])
}

func formatIPv6(buf []byte) string {
	addr := fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x",
		buf[0:2], buf[2:4], buf[4:6], buf[6:8], buf[8:10], buf[10:12], buf[12:14], buf[14:])
	return shorterIPv6(addr)
}

func shorterIPv6(addr string) string {
	return ipaddr.NewIPAddressString(addr).GetAddress().ToCompressedString()
}

func ipv6MyMask(maskLen int) []byte {
	bufLen := 16
	buf := make([]byte, bufLen)
	b := bufLen - 1
	totalBits := bufLen * 8
	zeroBits := totalBits - maskLen
	for i := 0; i < totalBits; i++ {
		offset := i % 8
		if i > 0 && i%8 == 0 {
			b--
		}
		if i < zeroBits {
			buf[b] |= (0 << offset)
		} else {
			buf[b] |= (1 << offset)
		}
	}
	return buf
}

func incv6IP(ip []byte) []byte {
	n := len(ip) - 1
	if ip[n] == 255 {
		ip[n] = 0
		incv6IP(ip[:n])
	} else {
		ip[n]++
	}
	return ip
}

func ip2byte(ipaddr string) []byte {
	ip := net.ParseIP(ipaddr)
	buf := big.NewInt(0)
	buf.SetBytes(ip)
	v, err := buf.GobEncode()
	if err != nil {
		fmt.Printf("not able to convert err %s", err)
		return []byte{}
	}
	return v[1:]
}

func mask2int(m string) int {
	mask, err := strconv.Atoi(m)
	if err != nil {
		fmt.Printf("not able to convert to integer err %s", err)
		return 0
	}
	return mask
}
