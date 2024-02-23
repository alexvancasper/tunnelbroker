package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShorterIPv6(t *testing.T) {
	req := require.New(t)
	tests := map[string]struct {
		addr string
		want string
	}{
		"empty":      {addr: "0000:0000:0000:0000:0000:0000:0000:0000", want: "::"},
		"pre-test-5": {addr: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", want: "2001:db8:85a3::8a2e:370:7334"},
		"pre-test-6": {addr: "2001:0db8:05a3:0000:0000:0000:0370:0334", want: "2001:db8:5a3::370:334"},
		"pre-test-7": {addr: "2001:0db8:05a3::", want: "2001:db8:5a3::"},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			res := shorterIPv6(testCase.addr)
			req.Equal(testCase.want, res)
		})
	}
}

func TestExpandIP(t *testing.T) {
	req := require.New(t)
	tests := map[string]struct {
		addr string
		want string
	}{
		"empty":      {addr: "::", want: "0000:0000:0000:0000:0000:0000:0000:0000"},
		"pre-test-5": {addr: "2001:db8:85a3::8a2e:370:7334", want: "2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
		"pre-test-6": {addr: "2001:db8:5a3::370:334", want: "2001:0db8:05a3:0000:0000:0000:0370:0334"},
		"pre-test-7": {addr: "2001:db8:5a3::", want: "2001:0db8:05a3:0000:0000:0000:0000:0000"},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			res := expandIPv6(testCase.addr)
			req.Equal(testCase.want, res)
		})
	}
}
