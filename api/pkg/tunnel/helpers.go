package tunnels

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func getLocalIPv4(logf *logrus.Logger) string {
	l := logf.WithFields(logrus.Fields{
		"function": "getPrefix",
	})
	addr := os.Getenv("IPv4LOCALADDR")
	if len(addr) == 0 {
		l.Errorf("Environment variable IPv4LOCALADDR is empty")
		return ""
	}
	return addr
}

func getPrefix(prefixlen int, logf *logrus.Logger) string {
	l := logf.WithFields(logrus.Fields{
		"function": "getPrefix",
	})
	requestURL := fmt.Sprintf("http://%s/acquire?prefixlen=%d", os.Getenv("IPAM"), prefixlen)
	l.Debugf("requestURL %s", requestURL)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		l.Errorf("client: could not create request: %s", err)
		return ""
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Errorf("client: error making http request: %s", err)
		return ""
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		l.Errorf("client: could not read response body: %s", err)
		return ""
	}
	l.Debugf("client: response body: %s", resBody)
	var Prefix struct {
		Prefix string `json:"prefix"`
	}
	err = json.Unmarshal(resBody, &Prefix)
	if err != nil {
		l.Errorf("client: unmarshalling error: %s", err)
		return ""
	}
	return Prefix.Prefix
}

func generateName(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s%s", sha1_hash[:3], sha1_hash[len(sha1_hash)-3:])
}
