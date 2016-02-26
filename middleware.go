// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

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
