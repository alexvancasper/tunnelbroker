package tunnels

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGetEndpoints(t *testing.T) {
	var MyLogger = logrus.New()

	ipv6str := "3a03:abcd:1805:0000:0000:0000:0000:001a/127"
	start, end := GetEndpoints(ipv6str, MyLogger)
	startWant := "3a03:abcd:1805::1a/127"
	endWant := "3a03:abcd:1805::1b/127"
	if start != startWant {
		t.Errorf("got %q, wanted start %q\n", start, startWant)
	}
	if end != endWant {
		t.Errorf("got %q, wanted end %q\n", end, endWant)
	}
}

func TestGetNetworkAddr(t *testing.T) {
	var MyLogger = logrus.New()

	ipv6str := "3a03:abcd:1805:0000:0000:0000:0000:0005/127"
	start := GetNetworkAddr(ipv6str, MyLogger)
	startWant := "3a03:abcd:1805::4/127"
	if start != startWant {
		t.Errorf("got %q, wanted start %q\n", start, startWant)
	}
}
