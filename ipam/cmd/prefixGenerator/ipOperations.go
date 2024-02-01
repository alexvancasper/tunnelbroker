package main

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	psql "github.com/alexvancasper/TunnelBroker/ipam/internal/database"
	sh "github.com/alexvancasper/TunnelBroker/ipam/internal/servicesHandler"
	formatter "github.com/fabienm/go-logrus-formatters"

	"github.com/sirupsen/logrus"
)

// func LoadPrefix(ipv6addr string) bool {
// 	val := strings.Split(ipv6addr, "/")
// 	ipv6addr = val[0]
// 	ipv6mask, err := strconv.Atoi(val[1])
// 	if err != nil {
// 		return false
// 	}
// 	ipv6 := net.ParseIP(ipv6addr)
// 	IPv6Int := big.NewInt(0)
// 	IPv6Int.SetBytes(ipv6)
// 	v, err := IPv6Int.GobEncode()
// 	if err != nil {
// 		return false
// 	}
// 	d.basePrefix = v[1:]
// 	d.prefixLength = ipv6mask
// 	return true
// }

// func isExist(prefix string) bool {
// 	for _, allocatedPrefix := range d.db {
// 		if prefix == allocatedPrefix {
// 			return true
// 		}
// 	}
// 	return false
// }

// func GetPrefix(LenOfPrefix int) string {
// 	m := ipv6MyMask(LenOfPrefix)
// 	if len(d.db) == 0 {
// 		clientPrefix := GetFreeIPv6Prefix(d.basePrefix, m)
// 		clientPrefix = fmt.Sprintf("%s/%d", clientPrefix, LenOfPrefix)
// 		d.db = append(d.db, clientPrefix)
// 		return clientPrefix
// 	}

// 	clientPrefix := GetFreeIPv6Prefix(d.db[len(d.db)-1], m)
// 	for d.isExist(clientPrefix) {
// 		clientPrefix = incv6IP(clientPrefix)
// 		clientPrefix := GetFreeIPv6Prefix(clientPrefix, m)
// 	}
// 	return clientPrefix
// }

// func main() {
// ipv6str := "3A03:ABCD:1805:1000::/64"
// val := strings.Split(ipv6str, "/")
// ipv6addr := val[0]
// ipv6mask, _ := strconv.Atoi(val[1])
// ipv6 := net.ParseIP(ipv6addr)
// IPv6Int := big.NewInt(0)
// IPv6Int.SetBytes(ipv6)
// v, _ := IPv6Int.GobEncode()
// v = v[1:]

//count := math.Pow(2, float64(128-ipv6mask))
//fmt.Printf("Total amount of addresses per network: %.0f\n", count)

// m := ipv6MyMask(ipv6mask)
// pref := nextIPv6Range(v, m, 1)
// fmt.Printf("%s\n", pref)
// fmt.Println("----------------")
// var prefixDB IPv6DB
// prefixDB.LoadPrefix("3A03:ABCD:1805::/48")
// prefix := prefixDB.GetPrefix(127)
// fmt.Printf("allocated: %s\n", prefix)
// prefix = prefixDB.GetPrefix(127)
// fmt.Printf("allocated: %s\n", prefix)
// }

func main() {
	//Initialize Logging connections
	var MyLogger = logrus.New()

	gelfFmt := formatter.NewGelf("Prefix Generator")
	MyLogger.SetFormatter(gelfFmt)
	MyLogger.SetOutput(os.Stdout)
	loglevel, err := logrus.ParseLevel("debug")
	if err != nil {
		MyLogger.WithField("function", "main").Fatalf("error %v", err)
	}
	MyLogger.SetLevel(loglevel)
	sh.SHandler.Log = MyLogger

	//Initialize database connections
	sh.SHandler.DB, err = psql.New(os.Getenv("APP_DSN"), 5)
	sh.SHandler.Timeout = time.Duration(1000)
	if err != nil {
		MyLogger.WithField("DSN", os.Getenv("APP_DSN"))
		MyLogger.Fatal(err)
	}
	defer sh.SHandler.DB.CloseConnection()
	MyLogger.Info("Database connection successfully")
	err = psql.MigrationUP(sh.SHandler.DB)
	if err != nil {
		MyLogger.Fatal(err)
	}
	MyLogger.Info("Migration is done")

	ipv6str := "3A03:ABCD:1805:0000::/127"

	val := strings.Split(ipv6str, "/")
	ipv6addr := val[0]
	ipv6mask, err := strconv.Atoi(val[1])
	if err != nil {
		return
	}
	ipv6 := net.ParseIP(ipv6addr)
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(ipv6)
	v, err := IPv6Int.GobEncode()
	if err != nil {
		return
	}
	prefixBytes := v[1:]
	maskBytes := ipv6MyMask(ipv6mask)
	start, end := nextIPv6Range(prefixBytes, maskBytes, 1)
	ctx := context.TODO()
	query := "INSERT INTO p2p (prefix, released) VALUES"
	n := 1000
	for i := 0; i < n; i++ {
		pref := fmt.Sprintf("%s/%d", formatIPv6(start), ipv6mask)
		if i == n-1 {
			query += fmt.Sprintf(" ('%s', true)", pref)
		} else {
			query += fmt.Sprintf(" ('%s', true),", pref)
		}
		start, end = nextIPv6Range(end, maskBytes, 1)
	}
	query += ";"
	err = sh.SHandler.DB.Conn.Ping(ctx)
	if err != nil {
		fmt.Println("error ping")
		return
	}
	_, err = sh.SHandler.DB.Conn.Query(ctx, query)
	if err != nil {
		fmt.Println("error insert")
		return
	}
	MyLogger.Infof("Prefix inserted successfully")
	fmt.Println(query)

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
			buf[b] = buf[b] | (0 << offset)
		} else {
			buf[b] = buf[b] | (1 << offset)
		}
	}
	return buf
}

func formatIPv6(buf []byte) string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x", buf[0:2], buf[2:4], buf[4:6], buf[6:8], buf[8:10], buf[10:12], buf[12:14], buf[14:])
}

func nextIPv6Range(ip, mask []byte, rangeCount int) ([]byte, []byte) {
	f, l := nextv6Range(ip, mask)
	for rangeCount > 0 {
		f = incv6IP(l)
		f, l = nextv6Range(f, mask)
		rangeCount--
	}
	return f, l
}
func nextv6Range(ip, mask []byte) ([]byte, []byte) {
	firstIP := make([]byte, 16)
	lastIP := make([]byte, 16)

	for idx, _ := range ip {
		firstIP[idx] = ip[idx] & mask[idx]
		lastIP[idx] = ip[idx] | ^mask[idx]
	}
	return firstIP, lastIP
}

func incv6IP(ip []byte) []byte {
	if ip[15] == 255 {
		ip[15] = 0
		if ip[14] == 255 {
			ip[14] = 0
			if ip[13] == 255 {
				ip[13] = 0
				if ip[12] == 255 {
					ip[12] = 0
					if ip[11] == 255 {
						ip[11] = 0
						if ip[10] == 255 {
							ip[10] = 0
							if ip[9] == 255 {
								ip[9] = 0
								if ip[8] == 255 {
									ip[8] = 0
									if ip[7] == 255 {
										ip[7] = 0
										if ip[6] == 255 {
											ip[6] = 0
											if ip[5] == 255 {
												ip[5] = 0
												if ip[4] == 255 {
													ip[4] = 0
													if ip[3] == 255 {
														ip[3] = 0
														if ip[2] == 255 {
															ip[2] = 0
															if ip[1] == 255 {
																ip[1] = 0
																if ip[0] == 255 {
																	ip[0] = 0
																} else {
																	ip[0] += 1
																}
															} else {
																ip[1] += 1
															}
														} else {
															ip[2] += 1
														}
													} else {
														ip[3] += 1
													}
												} else {
													ip[4] += 1
												}
											} else {
												ip[5] += 1
											}
										} else {
											ip[6] += 1
										}
									} else {
										ip[7] += 1
									}
								} else {
									ip[8] += 1
								}
							} else {
								ip[9] += 1
							}
						} else {
							ip[10] += 1
						}
					} else {
						ip[11] += 1
					}
				} else {
					ip[12] += 1
				}
			} else {
				ip[13] += 1
			}
		} else {
			ip[14] += 1
		}
	} else {
		ip[15] += 1
	}
	return ip
}
