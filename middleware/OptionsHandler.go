package middleware

import (
	"fmt"
	"net/http"
)

func OptionsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "OPTIONS" {
			fmt.Fprintf(writer, "")
			return
		} else {
			h.ServeHTTP(writer, request)
		}
	})
}
