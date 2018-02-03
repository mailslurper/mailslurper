// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

//go:generate esc -o ./www/www.go -pkg www -ignore DS_Store|README\.md|LICENSE|www\.go -prefix /www/ ./www

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mailslurper/mailslurper/cmd/mailslurper/www"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mailslurper/mailslurper/cmd/mailslurper/controllers"
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
	"github.com/mailslurper/mailslurper/pkg/ui"
	"github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
)

const (
	// Version of the MailSlurper Server application
	SERVER_VERSION string = "1.12.0"

	// Set to true while developing
	DEBUG_ASSETS bool = false

	CONFIGURATION_FILE_NAME string = "config.json"
)

var config *mailslurper.Configuration
var database mailslurper.IStorage
var logger *logrus.Entry
var serviceTierConfig *mailslurper.ServiceTierConfiguration
var renderer *ui.TemplateRenderer
var smtpListener *mailslurper.SMTPListener

var logFormat = flag.String("logformat", "simple", "Format for logging. 'simple' or 'json'. Default is 'simple'")
var logLevel = flag.String("loglevel", "info", "Level of logs to write. Valid values are 'debug', 'info', or 'error'. Default is 'info'")

func main() {
	var err error
	flag.Parse()

	logger = mailslurper.GetLogger(*logLevel, *logFormat, "MailSlurper")
	logger.Infof("Starting MailSlurper Server v%s", SERVER_VERSION)

	renderer = loadRenderer()

	/*
	 * Load configuration
	 */
	if config, err = mailslurper.LoadConfigurationFromFile(CONFIGURATION_FILE_NAME); err != nil {
		logger.Errorf("There was an error reading the configuration file '%s': %s", CONFIGURATION_FILE_NAME, err.Error())
		os.Exit(-1)
	}

	/*
	 * Setup global database connection handle
	 */
	storageType, databaseConnection := config.GetDatabaseConfiguration()

	if database, err = mailslurper.ConnectToStorage(storageType, databaseConnection, logger); err != nil {
		logger.Errorf("Error connecting to storage type '%d' with a connection string of %s: %s", int(storageType), databaseConnection.String(), err.Error())
		os.Exit(-1)
	}

	defer database.Disconnect()

	/*
	 * Setup the server pool
	 */
	pool := mailslurper.NewServerPool(mailslurper.GetLogger(*logLevel, *logFormat, "SMTP Server Pool"), config.MaxWorkers)

	/*
	 * Setup receivers (subscribers) to handle new mail items.
	 */
	receivers := []mailslurper.IMailItemReceiver{
		mailslurper.NewDatabaseReceiver(database, mailslurper.GetLogger(*logLevel, *logFormat, "Database Receiver")),
	}

	/*
	 * Setup the SMTP listener
	 */
	smtpListenerContext, smtpListenerCancel := context.WithCancel(context.Background())

	if smtpListener, err = mailslurper.NewSMTPListener(
		mailslurper.GetLogger(*logLevel, *logFormat, "SMTP Listener"),
		config,
		pool,
		receivers,
	); err != nil {
		logger.Errorf("There was a problem starting the SMTP listener: %s", err.Error())
		os.Exit(0)
	}

	/*
	 * Start the SMTP listener
	 */
	if err = smtpListener.Start(); err != nil {
		logger.Fatalf("Error starting SMTP listener: %s", err.Error())
	}

	smtpListener.Dispatch(smtpListenerContext)

	/*
	 * Setup and start the HTTP listener for the application site
	 */
	adminController := controllers.NewAdminController(mailslurper.GetLogger(*logLevel, *logFormat, "AdminController"), renderer, SERVER_VERSION, config, CONFIGURATION_FILE_NAME, DEBUG_ASSETS)
	admin := echo.New()
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
		logger.Infof("HTTP admin listener running on %s", config.GetFullWWWBindingAddress())

		if err := admin.Start(config.GetFullWWWBindingAddress()); err != nil {
			logger.Info("Shutting down HTTP admin listener")
		}
	}()

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

	if config.AutoStartBrowser {
		startBrowser(config)
	}

	/*
	 * Block this thread until we get an interrupt signal. Once we have that
	 * start shutting everything down
	 */
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT)

	<-quit

	ctx, cancel := context.WithTimeout(smtpListenerContext, 10*time.Second)
	defer cancel()

	smtpListenerCancel()

	if err = admin.Shutdown(ctx); err != nil {
		logger.Fatalf("Error shutting down admin listener: %s", err.Error())
	}

	if err = service.Shutdown(ctx); err != nil {
		logger.Fatalf("Error shutting down service listener: %s", err.Error())
	}
}

func startBrowser(config *mailslurper.Configuration) {
	timer := time.NewTimer(time.Second)

	go func() {
		<-timer.C
		logger.Infof("Opening web browser to http://%s:%d", config.WWWAddress, config.WWWPort)
		err := open.Start(fmt.Sprintf("http://%s:%d", config.WWWAddress, config.WWWPort))
		if err != nil {
			logger.Infof("ERROR - Could not open browser - %s", err.Error())
		}
	}()
}

func loadRenderer() *ui.TemplateRenderer {
	return ui.NewTemplateRenderer(DEBUG_ASSETS)
}
