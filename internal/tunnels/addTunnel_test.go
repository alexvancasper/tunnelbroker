package tunnels

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var hasher = func(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	shaHash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s%s", shaHash[:3], shaHash[len(shaHash)-3:])
}

func TestGenerateName(t *testing.T) {
	req := require.New(t)

	tests := map[string]struct {
		addr string
		want string
	}{
		"pre-test-1": {addr: fmt.Sprintf("%s%s", "1", "185.60.45.135"), want: hasher(fmt.Sprintf("%s%s", "1", "185.60.45.135"))},
		"pre-test-2": {addr: fmt.Sprintf("%s%s", "2", "185.60.45.135"), want: hasher(fmt.Sprintf("%s%s", "2", "185.60.45.135"))},
		"pre-test-3": {addr: fmt.Sprintf("%s%s", "1", "185.60.45.136"), want: hasher(fmt.Sprintf("%s%s", "1", "185.60.45.136"))},
		"pre-test-4": {addr: fmt.Sprintf("%s%s", "2", "185.60.45.136"), want: hasher(fmt.Sprintf("%s%s", "2", "185.60.45.136"))},
		"pre-test-5": {addr: fmt.Sprintf("%s%s", "172", "185.60.45.137"), want: hasher(fmt.Sprintf("%s%s", "172", "185.60.45.137"))},
		"pre-test-6": {addr: fmt.Sprintf("%s%s", "173", "185.60.45.137"), want: hasher(fmt.Sprintf("%s%s", "173", "185.60.45.137"))},
		"pre-test-7": {addr: fmt.Sprintf("%s%s", "174", "185.60.45.137"), want: hasher(fmt.Sprintf("%s%s", "174", "185.60.45.137"))},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			res := generateName(testCase.addr)
			req.Equal(testCase.want, res)
		})
	}
}
