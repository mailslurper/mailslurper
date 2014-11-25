package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(writer, request)
		log.Printf("APP - %s - %s (%v)\n", request.Method, request.URL.Path, time.Since(startTime))
	})
}
