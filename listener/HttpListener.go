// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package listener

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mailslurper/mailslurper/controllers/homeController"
	"github.com/mailslurper/mailslurper/controllers/settingsController"
	"github.com/mailslurper/mailslurper/controllers/webSocketController"
	"github.com/mailslurper/mailslurper/middleware"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

/*
Sets up and returns an HTTP server structure. This configures
middleware for logging and access control.
*/
func NewHttpListener(address string, port int) *http.Server {
	router := setupHttpRouter()

	server := alice.New(
		middleware.AccessControl,
		middleware.OptionsHandler,
		middleware.Logger).Then(router)

	listener := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: server,
	}

	return listener
}

/*
Sets the HTTP routes
*/
func setupHttpRouter() http.Handler {
	router := mux.NewRouter()

	// Home
	router.HandleFunc("/", homeController.Index).Methods("GET")
	router.HandleFunc("/servicesettings", settingsController.GetServiceSettings).Methods("GET", "OPTIONS")

	// Static requests
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./www/assets"))))

	// Web-sockets
	router.HandleFunc("/ws", webSocketController.WebSocketHandler)

	/*

		// Configuration
		requestRouter.HandleFunc("/configuration", controllers.Config).Methods("GET")
		requestRouter.HandleFunc("/config", controllers.GetConfig).Methods("GET")
		requestRouter.HandleFunc("/config", controllers.SaveConfig).Methods("PUT")

		// Web-sockets
		requestRouter.HandleFunc("/ws", smtp.WebsocketHandler)

	*/
	return router
}

/*
Starts the HTTP listener and serves Service requests
*/
func StartHttpListener(httpListener *http.Server) error {
	log.Println("INFO - HTTP App listener started on", httpListener.Addr)
	return httpListener.ListenAndServe()
}
