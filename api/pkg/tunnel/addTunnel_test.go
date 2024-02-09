package tunnels

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateName(t *testing.T) {
	req := require.New(t)

	tests := map[string]struct {
		addr string
		want string
	}{
		"pre-test-1": {addr: fmt.Sprintf("%s%s", "1", "185.60.45.135"), want: "a58c2"},
		"pre-test-2": {addr: fmt.Sprintf("%s%s", "2", "185.60.45.135"), want: "3025c"},
		"pre-test-3": {addr: fmt.Sprintf("%s%s", "1", "185.60.45.136"), want: "d40bb"},
		"pre-test-4": {addr: fmt.Sprintf("%s%s", "2", "185.60.45.136"), want: "9b9e4"},
		"pre-test-5": {addr: fmt.Sprintf("%s%s", "172", "185.60.45.137"), want: "9fff6"},
		"pre-test-6": {addr: fmt.Sprintf("%s%s", "173", "185.60.45.137"), want: "4db3d"},
		"pre-test-7": {addr: fmt.Sprintf("%s%s", "174", "185.60.45.137"), want: "e865b"},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			res := GenerateName(testCase.addr)
			req.Equal(testCase.want, res)
		})
	}

}
