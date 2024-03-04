package tunnels

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCorrectIP(t *testing.T) {
	req := require.New(t)

	tests := map[string]struct {
		addr string
		want bool
	}{
		"empty":              {addr: "", want: false},
		"1.1.1.1":            {addr: "1.1.1.1", want: true},
		"private-192":        {addr: "192.168.0.0", want: false},
		"private-172-16":     {addr: "172.16.0.1", want: false},
		"private-10":         {addr: "10.195.134.1", want: false},
		"public-192":         {addr: "192.169.10.12", want: true},
		"wrong-172-17":       {addr: "172.17.345.1", want: false},
		"private-172-17":     {addr: "172.17.245.1", want: false},
		"public-172-32":      {addr: "172.32.245.1", want: true},
		"numeric":            {addr: "0000000000", want: false},
		"public-1space":      {addr: " 185.60.45.135", want: true},
		"public-2space":      {addr: "185.60.45.135 ", want: true},
		"public-3space":      {addr: " 185.60.45.135 ", want: true},
		"public-many-space":  {addr: "    185.60.45.135     ", want: true},
		"public-space-space": {addr: "    185. 60. 45 . 135     ", want: false},
		"pre-test-1":         {addr: "My IP address: 192.168.0.1! Yay!", want: false},
		"pre-test-2":         {addr: "185.60.45.135", want: true},
		"pre-test-3":         {addr: "255.255.255.255", want: false},
		"pre-test-4":         {addr: "0.0.0.0", want: false},
		"pre-test-5":         {addr: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", want: false},
		"pre-test-6":         {addr: "185.60.45.135 185.60.45.136", want: true},
		"pre-test-7":         {addr: "2001:0db8:85a3:0000:0000:8a2e:0370:7334 185.60.45.135", want: true},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			res := IsIPv4Public(testCase.addr)
			req.Equal(testCase.want, res)
		})
	}
}
