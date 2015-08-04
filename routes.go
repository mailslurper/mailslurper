package main

import (
	"github.com/mailslurper/mailslurper/controllers"
	"github.com/mailslurper/mailslurper/services/listener"
	"github.com/mailslurper/mailslurper/services/middleware"
)

/*
Add routes here using AddRoute and AddRouteWithMiddleware.
*/
func setupRoutes(httpListener *listener.HTTPListenerService, appContext *middleware.AppContext) {
	httpListener.
		AddStaticRoute("/assets/", "./www/assets").
		AddRoute("/", controllers.Index, "GET").
		AddRoute("/savedsearches", controllers.ManageSavedSearches, "GET").
		AddRoute("/version", controllers.GetVersion, "GET", "OPTIONS").
		AddRoute("/servicesettings", controllers.GetServiceSettings, "GET", "OPTIONS")
}
