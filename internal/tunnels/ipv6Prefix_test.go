package tunnels

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEndpoints(t *testing.T) {
	req := require.New(t)

	tests := map[string]struct {
		v6addr   string
		network  string
		nextAddr string
	}{
		"pre-test-5": {
			v6addr:   "2001:0db8:85a3:0000:0000:8a2e:0370:7334/127",
			network:  "2001:db8:85a3::8a2e:370:7334/127",
			nextAddr: "2001:db8:85a3::8a2e:370:7335/127",
		},
		"pre-test-7": {
			v6addr:   "2001:0db8:85a3:0000:0000:8a2e:0370:0001/127",
			network:  "2001:db8:85a3::8a2e:370:0/127",
			nextAddr: "2001:db8:85a3::8a2e:370:1/127",
		},
		"pre-test-8": {
			v6addr:   "2001:0db8:85a3::fafd/64",
			network:  "2001:db8:85a3::/64",
			nextAddr: "2001:db8:85a3::1/64",
		},
		"pre-test-9": {
			v6addr:   "2001:db8:85a3::/87",
			network:  "2001:db8:85a3::/87",
			nextAddr: "2001:db8:85a3::1/87",
		},
		"pre-test-10": {
			v6addr:   "2F7F:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF/3",
			network:  "2000::/3",
			nextAddr: "2000::1/3",
		},
		"pre-test-1": {
			v6addr:   "2001:0db8:85a3:0000:0000:8a2e:0370:7334/127",
			network:  "2001:db8:85a3::8a2e:370:7334/127",
			nextAddr: "2001:db8:85a3::8a2e:370:7335/127",
		},
		"pre-test-2": {
			v6addr:   "2001:0db8:85a3:0000:0000:8a2e:0370:FFFF/64",
			network:  "2001:db8:85a3::/64",
			nextAddr: "2001:db8:85a3::1/64",
		},
		"pre-test-3": {
			v6addr:   "2001:0db8:85a3::/87",
			network:  "2001:db8:85a3::/87",
			nextAddr: "2001:db8:85a3::1/87",
		},
		"pre-test-4": {
			v6addr:   "2001:0db8:85a3:FFFF:FFFF:FFFF:FFFF:FFFF/10",
			network:  "2000::/10",
			nextAddr: "2000::1/10",
		},
		"pre-test-11": {
			v6addr:   "2FFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF/3",
			network:  "2000::/3",
			nextAddr: "2000::1/3",
		},
		"pre-test-12": {
			v6addr:   "2FFF:FFFF:FFFF:FFFF:FFFF:FF00:FF00:FFFF/112",
			network:  "2fff:ffff:ffff:ffff:ffff:ff00:ff00::/112",
			nextAddr: "2fff:ffff:ffff:ffff:ffff:ff00:ff00:1/112",
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			network, nextAddr := GetEndpoints(testCase.v6addr)
			req.Equal(testCase.network, network)
			req.Equal(testCase.nextAddr, nextAddr)
		})
	}
}

func TestGetNetworkAddr(t *testing.T) {
	req := require.New(t)

	tests := map[string]struct {
		v6addr string
		want   string
	}{
		"pre-test-5":  {v6addr: "2001:0db8:85a3:0000:0000:8a2e:0370:7334/128", want: "2001:db8:85a3::8a2e:370:7334/128"},
		"pre-test-7":  {v6addr: "2001:0db8:85a3:0000:0000:8a2e:0370:0001/127", want: "2001:db8:85a3::8a2e:370:0/127"},
		"pre-test-8":  {v6addr: "2001:0db8:85a3::fafd/64", want: "2001:db8:85a3::/64"},
		"pre-test-9":  {v6addr: "2001:db8:85a3::/87", want: "2001:db8:85a3::/87"},
		"pre-test-10": {v6addr: "2F7F:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF/3", want: "2000::/3"},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			res, _ := GetEndpoints(testCase.v6addr)
			req.Equal(testCase.want, res)
		})
	}
}

func TestIncIPv6(t *testing.T) {
	req := require.New(t)

	tests := map[string]struct {
		ip   []byte
		want []byte
	}{
		"pre-test-5":  {ip: ip2byte("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), want: ip2byte("2001:0db8:85a3:0000:0000:8a2e:0370:7335")},
		"pre-test-7":  {ip: ip2byte("2001:0db8:85a3:0000:0000:8a2e:0370:FFFF"), want: ip2byte("2001:0db8:85a3:0000:0000:8a2e:0371:0000")},
		"pre-test-8":  {ip: ip2byte("2001:0db8:85a3::"), want: ip2byte("2001:0db8:85a3::1")},
		"pre-test-9":  {ip: ip2byte("2001:0db8:85a3:FFFF:FFFF:FFFF:FFFF:FFFF"), want: ip2byte("2001:0db8:85a4:0000:0000:0000:0000:0000")},
		"pre-test-10": {ip: ip2byte("2FFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF"), want: ip2byte("3000:0000:0000:0000:0000:0000:0000:0000")},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			res := incv6IP(testCase.ip)
			req.Equal(testCase.want, res)
		})
	}
}
