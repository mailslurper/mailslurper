package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mailslurper/mailslurper/cmd/mailslurper/controllers"
	"github.com/mailslurper/mailslurper/pkg/auth/authscheme"
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
)

func setupServicesListener() {
	middlewares := make([]echo.MiddlewareFunc, 0, 5)

	/*
	 * Start the services server
	 */
	serviceController := controllers.NewServiceController(mailslurper.GetLogger(*logLevel, *logFormat, "ServiceController"), SERVER_VERSION, config, database)
	service = echo.New()
	service.HideBanner = true

	if config.AuthenticationScheme != authscheme.NONE {
		middlewares = append(middlewares, serviceAuthorization)
	}

	service.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
	}))

	service.GET("/mail/:id", serviceController.GetMail, middlewares...)
	service.GET("/mail/:id/message", serviceController.GetMailMessage, middlewares...)
	service.DELETE("/mail", serviceController.DeleteMail, middlewares...)
	service.GET("/mail", serviceController.GetMailCollection, middlewares...)
	service.GET("/mailcount", serviceController.GetMailCount, middlewares...)
	service.GET("/mail/:mailID/attachment/:attachmentID", serviceController.DownloadAttachment, middlewares...)
	service.GET("/version", serviceController.Version, middlewares...)
	service.GET("/pruneoptions", serviceController.GetPruneOptions, middlewares...)
	service.POST("/login", serviceController.Login)

	go func() {
		var err error

		if config.IsServiceSSL() {
			err = service.StartTLS(config.GetFullServiceAppAddress(), config.CertFile, config.KeyFile)
		} else {
			err = service.Start(config.GetFullServiceAppAddress())
		}

		if err != nil {
			logger.WithError(err).Info("Shutting down HTTP service listener")
		} else {
			logger.Infof("Service listener running on %s", config.GetFullServiceAppAddress())
		}
	}()
}
