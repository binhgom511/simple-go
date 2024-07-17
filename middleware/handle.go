package middleware

import (
	"log"
	"net/http"
	"time"
)

// TimeRequest is a middleware function that logs the time taken to process each request.
func TimeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		log.Printf("Request took %s", elapsed)
	})
}
