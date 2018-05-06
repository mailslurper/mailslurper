package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/mailslurper/mailslurper/pkg/auth/authscheme"

	"github.com/labstack/echo"
	"github.com/mailslurper/mailslurper/cmd/mailslurper/controllers"
	"github.com/mailslurper/mailslurper/cmd/mailslurper/www"
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
)

func setupAdminListener() {
	middlewares := make([]echo.MiddlewareFunc, 0, 5)

	/*
	 * Setup and start the HTTP listener for the application site
	 */
	adminController := controllers.NewAdminController(mailslurper.GetLogger(*logLevel, *logFormat, "AdminController"), renderer, SERVER_VERSION, config, CONFIGURATION_FILE_NAME, DEBUG_ASSETS)
	admin = echo.New()
	admin.HideBanner = true
	admin.Renderer = renderer

	assetHandler := http.FileServer(www.FS(DEBUG_ASSETS))
	admin.GET("/www/*", echo.WrapHandler(assetHandler))

	if config.AuthenticationScheme != authscheme.NONE {
		admin.Use(session.Middleware(sessions.NewCookieStore([]byte(config.AdminCookieSecret))))
		middlewares = append(middlewares, adminAuthorization)

		admin.GET("/login", adminController.Login)
		admin.POST("/perform-login", adminController.PerformLogin)
	}

	admin.GET("/", adminController.Index, middlewares...)
	admin.GET("/admin", adminController.Admin, middlewares...)
	admin.GET("/savedsearches", adminController.ManageSavedSearches, middlewares...)
	admin.GET("/servicesettings", adminController.GetServiceSettings, middlewares...)
	admin.GET("/version", adminController.GetVersion, middlewares...)
	admin.GET("/masterversion", adminController.GetVersionFromMaster, middlewares...)
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
