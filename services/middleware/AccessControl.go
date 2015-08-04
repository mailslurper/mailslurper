package middleware

import (
	"net/http"
)

/*
AccessControl is a middleware that tells a browser abour CORS, allowed verbs,
and accepted headers. Modify this to change these security features.
*/
func (ctx *AppContext) AccessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CRSF-Token")
		h.ServeHTTP(writer, request)
	})
}
