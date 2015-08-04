package main

import (
	"github.com/mailslurper/mailslurper/services/listener"
	"github.com/mailslurper/mailslurper/services/middleware"
)

func setupMiddleware(httpListener *listener.HTTPListenerService, appContext *middleware.AppContext) {
	httpListener.
		AddMiddleware(appContext.Logger).
		AddMiddleware(appContext.StartAppContext).
		AddMiddleware(appContext.AccessControl).
		AddMiddleware(appContext.OptionsHandler)
}
