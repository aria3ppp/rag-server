package http

import (
	http_pkg "net/http"
)

type MockHandlers map[string]http_pkg.HandlerFunc

var _ http_pkg.Handler = (MockHandlers)(nil)

func (m MockHandlers) ServeHTTP(w http_pkg.ResponseWriter, r *http_pkg.Request) {
	mux := http_pkg.NewServeMux()
	for p, h := range m {
		mux.HandleFunc(p, h)
	}
	mux.ServeHTTP(w, r)
}
