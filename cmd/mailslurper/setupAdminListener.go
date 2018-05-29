// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/mailslurper/mailslurper/cmd/mailslurper/controllers"
	"github.com/mailslurper/mailslurper/cmd/mailslurper/www"
	"github.com/mailslurper/mailslurper/pkg/auth/authscheme"
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
)

func setupAdminListener() {
	middlewares := make([]echo.MiddlewareFunc, 0, 5)

	/*
	 * Setup and start the HTTP listener for the application site
	 */
	adminController := &controllers.AdminController{
		CacheService:   cacheService,
		Config:         config,
		ConfigFileName: CONFIGURATION_FILE_NAME,
		DebugMode:      DEBUG_ASSETS,
		Renderer:       renderer,
		Lock:           &sync.Mutex{},
		Logger:         mailslurper.GetLogger(*logLevel, *logFormat, "AdminController"),
		ServerVersion:  SERVER_VERSION,
	}

	admin = echo.New()
	admin.HideBanner = true
	admin.Renderer = renderer

	assetHandler := http.FileServer(www.FS(DEBUG_ASSETS))
	admin.GET("/www/*", echo.WrapHandler(assetHandler))

	if config.AuthenticationScheme != authscheme.NONE {
		admin.Use(session.Middleware(sessions.NewCookieStore([]byte(config.AuthSecret))))
		middlewares = append(middlewares, adminAuthorization)

		admin.GET("/login", adminController.Login)
		admin.POST("/perform-login", adminController.PerformLogin)
		admin.GET("/logout", adminController.Logout)
	}

	admin.GET("/", adminController.Index, middlewares...)
	admin.GET("/admin", adminController.Admin, middlewares...)
	admin.GET("/savedsearches", adminController.ManageSavedSearches, middlewares...)
	admin.GET("/servicesettings", adminController.GetServiceSettings)
	admin.GET("/version", adminController.GetVersion)
	admin.GET("/masterversion", adminController.GetVersionFromMaster)
	admin.POST("/theme", adminController.ApplyTheme, middlewares...)

	go func() {
		var err error

		if config.IsAdminSSL() {
			logger.Infof("HTTP admin listener running on SSL %s", config.GetFullWWWBindingAddress())

			if err = admin.StartTLS(config.GetFullWWWBindingAddress(), config.AdminCertFile, config.AdminKeyFile); err != nil {
				logger.WithError(err).Info("Shutting down HTTP admin listener")
			}
		} else {
			logger.Infof("HTTP admin listener running on %s", config.GetFullWWWBindingAddress())

			if err = admin.Start(config.GetFullWWWBindingAddress()); err != nil {
				logger.WithError(err).Info("Shutting down HTTP admin listener")
			}
		}
	}()
}
