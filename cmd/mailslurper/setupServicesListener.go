package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mailslurper/mailslurper/cmd/mailslurper/controllers"
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
)

func setupServicesListener() {
	/*
	 * Start the services server
	 */
	serviceController := controllers.NewServiceController(mailslurper.GetLogger(*logLevel, *logFormat, "ServiceController"), SERVER_VERSION, config, database)
	service := echo.New()
	service.HideBanner = true

	service.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
	}))

	service.GET("/mail/:id", serviceController.GetMail)
	service.GET("/mail/:id/message", serviceController.GetMailMessage)
	service.DELETE("/mail", serviceController.DeleteMail)
	service.GET("/mail", serviceController.GetMailCollection)
	service.GET("/mailcount", serviceController.GetMailCount)
	service.GET("/mail/:mailID/attachment/:attachmentID", serviceController.DownloadAttachment)
	service.GET("/version", serviceController.Version)
	service.GET("/pruneoptions", serviceController.GetPruneOptions)

	go func() {
		var err error

		if config.CertFile != "" && config.KeyFile != "" {
			err = service.StartTLS(config.GetFullServiceAppAddress(), config.CertFile, config.KeyFile)
		} else {
			err = service.Start(config.GetFullServiceAppAddress())
		}

		if err != nil {
			logger.Info("Shutting down HTTP service listener")
		} else {
			logger.Infof("Service listener running on %s", config.GetFullServiceAppAddress())
		}
	}()
}
