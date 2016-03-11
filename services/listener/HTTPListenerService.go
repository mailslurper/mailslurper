// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package listener

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/mailslurper/libmailslurper/configuration"
	"github.com/mailslurper/mailslurper/services/middleware"
	"github.com/mailslurper/mailslurper/www"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

/*
HTTPListenerService is a structure which provides an HTTP listener to service
requests. This structure offers methods to add routes and middlewares. Typical
usage would first call NewHTTPListenerService(), add routes, then call
StartHTTPListener.
*/
type HTTPListenerService struct {
	Address string
	Port    int
	Context *middleware.AppContext

	Router                 *mux.Router
	BaseMiddlewareHandlers alice.Chain
}

/*
NewHTTPListenerService creates a new instance of the HTTPListenerService
*/
func NewHTTPListenerService(
	address string,
	port int,
	appContext *middleware.AppContext,
) *HTTPListenerService {
	return &HTTPListenerService{
		Address: address,
		Port:    port,
		Context: appContext,

		Router: mux.NewRouter(),
	}
}

/*
AddMiddleware adds a new middleware handler to the request chain.
*/
func (service *HTTPListenerService) AddMiddleware(middlewareHandler alice.Constructor) *HTTPListenerService {
	service.BaseMiddlewareHandlers = service.BaseMiddlewareHandlers.Append(middlewareHandler)
	return service
}

/*
AddRoute adds a HTTP handler route to the HTTP listener.
*/
func (service *HTTPListenerService) AddRoute(
	path string,
	handlerFunc http.HandlerFunc, methods ...string,
) *HTTPListenerService {
	service.Router.Handle(path, service.BaseMiddlewareHandlers.ThenFunc(handlerFunc)).Methods(methods...)
	return service
}

/*
AddRouteWithMiddleware adds a HTTP handler route that goes through an additional
middleware handler, to the HTTP listener.
*/
func (service *HTTPListenerService) AddRouteWithMiddleware(
	path string,
	handlerFunc http.HandlerFunc,
	middlewareHandler alice.Constructor,
	methods ...string,
) *HTTPListenerService {
	service.Router.Handle(
		path,
		service.BaseMiddlewareHandlers.Append(middlewareHandler).ThenFunc(handlerFunc),
	).Methods(methods...)

	return service
}

/*
AddStaticRoute adds a HTTP handler route for static assets.
*/
func (service *HTTPListenerService) AddStaticRoute(pathPrefix string, directory string) *HTTPListenerService {
	/*
		fileServer := http.FileServer(http.Dir("./www/assets"))
		service.Router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix, fileServer))
	*/
	service.Router.PathPrefix(pathPrefix).Handler(http.FileServer(www.FS(false)))
	return service
}

func (service *HTTPListenerService) gzipFileServer(dir http.FileSystem) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		root := fmt.Sprintf("%s", dir)

		file, err := ioutil.ReadFile(path.Join(root, request.URL.Path))
		if err == nil {
			extension := filepath.Ext(request.URL.Path)
			mimeType := mime.TypeByExtension(extension)

			writer.Header().Set("Content-Type", mimeType)

			if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
				writer.Write(file)
			} else {
				writer.Header().Set("Content-Encoding", "gzip")

				gz := gzip.NewWriter(writer)
				defer gz.Close()

				gz.Write(file)
			}
		} else {
			http.NotFound(writer, request)
		}
	})
}

/*
StartHTTPListener starts the HTTP listener and servicing requests.
*/
func (service *HTTPListenerService) StartHTTPListener(config *configuration.Configuration) error {
	listener := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", service.Address, service.Port),
		Handler: alice.New().Then(service.Router),
	}

	if config.CertFile != "" && config.KeyFile != "" {
		log.Printf("MailSlurper: INFO - HTTPS listener started on %s:%d\n", service.Address, service.Port)
		return listener.ListenAndServeTLS(config.CertFile, config.KeyFile)
	}

	log.Printf("MailSlurper: INFO - HTTP listener started on %s:%d\n", service.Address, service.Port)
	return listener.ListenAndServe()
}
