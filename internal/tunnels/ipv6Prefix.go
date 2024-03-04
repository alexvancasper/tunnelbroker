package tunnels

import (
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"

	"github.com/seancfoley/ipaddress-go/ipaddr"
	"github.com/sirupsen/logrus"
)

func GetEndpoints(ipv6str string, logf *logrus.Logger) (string, string) {
	l := logf.WithFields(logrus.Fields{
		"function": "GetEndpoints",
	})

	l.Infof("Input prefix %s", ipv6str)
	val := strings.Split(ipv6str, "/")
	ipv6addr := val[0]
	ipv6mask, err := strconv.Atoi(val[1])
	if err != nil {
		l.Errorf("not able to convert to integer err %s", err)
		return "", ""
	}
	ipv6 := net.ParseIP(ipv6addr)
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(ipv6)
	v, err := IPv6Int.GobEncode()
	if err != nil {
		l.Errorf("not able to encode to bytes err %s", err)
		return "", ""
	}
	ip6 := v[1:]
	mask := ipv6MyMask(ipv6mask)
	firstIP := make([]byte, 16)
	// for idx, _ := range ip6 {
	for i := 0; i < len(ipv6); i++ {
		firstIP[i] = ip6[i] & mask[i]
		// secondIP[idx] = ip6[idx] | ^mask[idx] //it will find last address
	}
	return fmt.Sprintf("%s/%s", formatIPv6(firstIP), val[1]), fmt.Sprintf("%s/%s", formatIPv6(incv6IP(firstIP)), val[1])
}

func GetNetworkAddr(ipv6str string, logf *logrus.Logger) string {
	l := logf.WithFields(logrus.Fields{
		"function": "GetNetworkAddr",
	})
	l.Infof("Input prefix %s", ipv6str)
	val := strings.Split(ipv6str, "/")
	ipv6addr := val[0]
	ipv6mask, err := strconv.Atoi(val[1])
	l.Debugf("Splitted prefixes prefix:%s mask:%s", val[0], val[1])
	if err != nil {
		l.Errorf("not able to convert to integer err %s", err)
		return ""
	}
	ipv6 := net.ParseIP(ipv6addr)
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(ipv6)
	v, err := IPv6Int.GobEncode()
	if err != nil {
		l.Errorf("not able to encode to bytes err %s", err)
		return ""
	}
	ip6 := v[1:]
	mask := ipv6MyMask(ipv6mask)
	networkAddr := make([]byte, 16)
	// for idx, _ := range ip6 {
	for i := 0; i < len(ipv6); i++ {
		networkAddr[i] = ip6[i] & mask[i]
		// secondIP[idx] = ip6[idx] | ^mask[idx] //it will find last address
	}
	return fmt.Sprintf("%s/%s", formatIPv6(networkAddr), val[1])
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
		incv6IP(ip[:n-1])
	} else {
		ip[n]++
	}
	return ip
	// if ip[15] == 255 {
	// 	ip[15] = 0
	// 	if ip[14] == 255 {
	// 		ip[14] = 0
	// 		if ip[13] == 255 {
	// 			ip[13] = 0
	// 			if ip[12] == 255 {
	// 				ip[12] = 0
	// 				if ip[11] == 255 {
	// 					ip[11] = 0
	// 					if ip[10] == 255 {
	// 						ip[10] = 0
	// 						if ip[9] == 255 {
	// 							ip[9] = 0
	// 							if ip[8] == 255 {
	// 								ip[8] = 0
	// 								if ip[7] == 255 {
	// 									ip[7] = 0
	// 									if ip[6] == 255 {
	// 										ip[6] = 0
	// 										if ip[5] == 255 {
	// 											ip[5] = 0
	// 											if ip[4] == 255 {
	// 												ip[4] = 0
	// 												if ip[3] == 255 {
	// 													ip[3] = 0
	// 													if ip[2] == 255 {
	// 														ip[2] = 0
	// 														if ip[1] == 255 {
	// 															ip[1] = 0
	// 															if ip[0] == 255 {
	// 																ip[0] = 0
	// 															} else {
	// 																ip[0]++
	// 															}
	// 														} else {
	// 															ip[1]++
	// 														}
	// 													} else {
	// 														ip[2]++
	// 													}
	// 												} else {
	// 													ip[3]++
	// 												}
	// 											} else {
	// 												ip[4]++
	// 											}
	// 										} else {
	// 											ip[5]++
	// 										}
	// 									} else {
	// 										ip[6]++
	// 									}
	// 								} else {
	// 									ip[7]++
	// 								}
	// 							} else {
	// 								ip[8]++
	// 							}
	// 						} else {
	// 							ip[9]++
	// 						}
	// 					} else {
	// 						ip[10]++
	// 					}
	// 				} else {
	// 					ip[11]++
	// 				}
	// 			} else {
	// 				ip[12]++
	// 			}
	// 		} else {
	// 			ip[13]++
	// 		}
	// 	} else {
	// 		ip[14]++
	// 	}
	// } else {
	// 	ip[15]++
	// }
}
