package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// RequestLogger logs method, path, status code, and duration for every request.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, wrapped.status, time.Since(start))
	})
}

// responseWriter captures the status code written by downstream handlers.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// MethodOverride checks for a _method field in POST form bodies and
// rewrites r.Method to the specified value (PUT or DELETE).
// This lets plain HTML forms submit PUT and DELETE requests.
func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			override := r.FormValue("_method")
			if override == "" {
				// Try the X-HTTP-Method-Override header as a fallback.
				override = r.Header.Get("X-HTTP-Method-Override")
			}
			override = strings.ToUpper(override)
			if override == http.MethodPut || override == http.MethodDelete || override == http.MethodPatch {
				r.Method = override
			}
		}
		next.ServeHTTP(w, r)
	})
}
