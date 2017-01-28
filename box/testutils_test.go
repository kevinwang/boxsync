package box

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func newTestServerClient(endpointPath, responseBody string) (*httptest.Server, Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != endpointPath {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, responseBody)
	}))

	client := &client{
		client:     &http.Client{},
		apiBaseUrl: server.URL,
	}

	return server, client
}
