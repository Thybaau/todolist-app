package middleware

import (
	"log"
	"net/http"
	"time"
)

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Create a custom ResponseWriter to get request status code
		crw := &customResponseWriter{ResponseWriter: w}

		// Call the next handler in the chain
		next.ServeHTTP(crw, r)

		// Log request details
		duration := time.Since(start)
		log.Printf("[%s] %v - %s %s - %s\n", r.Method, crw.status, r.RequestURI, r.RemoteAddr, duration)
	})
}

//ResponseWriter to get request status code
type customResponseWriter struct {
	http.ResponseWriter
	status int
}

func (crw *customResponseWriter) WriteHeader(statusCode int) {
	crw.status = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}
