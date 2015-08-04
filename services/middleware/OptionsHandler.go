package middleware

import (
	"fmt"
	"net/http"
)

/*
OptionsHandler is a middleware for handling OPTIONS requests.
*/
func (ctx *AppContext) OptionsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "OPTIONS" {
			fmt.Fprintf(writer, "")
			return
		}

		h.ServeHTTP(writer, request)
	})
}
