package tunnels

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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
	ctx, ctxCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer ctxCancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		l.Errorf("client: could not create request: %s", err)
		return ""
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Errorf("client: error making http request: %s", err)
		return ""
	}
	defer res.Body.Close()
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
	h := sha256.New()
	h.Write([]byte(s))
	shaHash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s%s", shaHash[:3], shaHash[len(shaHash)-3:])
}
