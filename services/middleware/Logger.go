// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package middleware

import (
	"log"
	"net/http"
	"time"
)

/*
Logger is a middleware which logs requests to the console. It also includes the
time it takes for the request to complete.
*/
func (ctx *AppContext) Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(writer, request)
		log.Printf("MailSlurper: INFO - %s - %s (%v)\n", request.Method, request.URL.String(), time.Since(startTime))
	})
}
