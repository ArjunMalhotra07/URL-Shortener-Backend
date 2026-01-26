package httpserver

import (
	"fmt"
	"net/http"
)

// NewMux returns a ServeMux with a basic health endpoint.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	return mux
}
