package mw

import (
	"log"
	"net/http"
)

// HTTPLogger Middleware HTTP request logger.
func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.URL, r.Method)
		next.ServeHTTP(w, r)
	})
}
