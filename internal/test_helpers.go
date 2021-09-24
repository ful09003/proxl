package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//test_helpers contains things like http servers used for unit tests.

func newTestHTTPServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/good_case" {
			// Send back a calid metrics string
			s := `# HELP joy_felt_total A counter of joy experienced.
# TYPE joy_felt_total counter
joy_felt_total{developer="me"} 9000
`
			rw.Write([]byte(s))
		}
	}))

	return server
}