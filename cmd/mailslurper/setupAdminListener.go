package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mailslurper/mailslurper/cmd/mailslurper/controllers"
	"github.com/mailslurper/mailslurper/cmd/mailslurper/www"
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
)

func setupAdminListener() {
	/*
	 * Setup and start the HTTP listener for the application site
	 */
	adminController := controllers.NewAdminController(mailslurper.GetLogger(*logLevel, *logFormat, "AdminController"), renderer, SERVER_VERSION, config, CONFIGURATION_FILE_NAME, DEBUG_ASSETS)
	admin = echo.New()
	admin.HideBanner = true
	admin.Renderer = renderer

	assetHandler := http.FileServer(www.FS(DEBUG_ASSETS))
	admin.GET("/www/*", echo.WrapHandler(assetHandler))

	admin.GET("/", adminController.Index)
	admin.GET("/admin", adminController.Admin)
	admin.GET("/savedsearches", adminController.ManageSavedSearches)
	admin.GET("/servicesettings", adminController.GetServiceSettings)
	admin.GET("/version", adminController.GetVersion)
	admin.GET("/masterversion", adminController.GetVersionFromMaster)
	admin.POST("/theme", adminController.ApplyTheme)

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
