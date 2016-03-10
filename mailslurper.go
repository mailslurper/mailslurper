// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adampresley/GoHttpService"
	"github.com/adampresley/sigint"
	"github.com/mailslurper/libmailslurper"
	"github.com/mailslurper/libmailslurper/configuration"
	"github.com/mailslurper/libmailslurper/receiver"
	"github.com/mailslurper/libmailslurper/server"
	"github.com/mailslurper/libmailslurper/storage"
	"github.com/mailslurper/mailslurper/global"
	"github.com/mailslurper/mailslurper/services/listener"
	"github.com/mailslurper/mailslurper/services/middleware"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	var err error

	log.Printf("MailSlurper: INFO - Starting MailSlurper Server v%s\n", global.SERVER_VERSION)
	/*
	 * Prepare SIGINT handler (CTRL+C)
	 */
	sigint.ListenForSIGINT(func() {
		log.Println("MailSlurper: INFO - Shutting down via SIGINT.")
		os.Exit(0)
	})

	/*
	 * Load configuration
	 */
	config, err := configuration.LoadConfigurationFromFile(configuration.CONFIGURATION_FILE_NAME)
	if err != nil {
		log.Println("MailSlurper: ERROR - There was an error reading your configuration file:", err)
		os.Exit(0)
	}

	/*
	 * Setup global database connection handle
	 */
	storageType, databaseConnection := config.GetDatabaseConfiguration()

	if global.Database, err = storage.ConnectToStorage(storageType, databaseConnection); err != nil {
		log.Println("MailSlurper: ERROR - There was an error connecting to your data storage:", err.Error())
		os.Exit(0)
	}

	defer global.Database.Disconnect()

	/*
	 * Setup the server pool
	 */
	pool := server.NewServerPool(config.MaxWorkers)

	/*
	 * Setup the SMTP listener
	 */
	smtpServer, err := server.SetupSMTPServerListener(config)
	if err != nil {
		log.Println("MailSlurper: ERROR - There was a problem starting the SMTP listener:", err)
		os.Exit(0)
	}

	defer server.CloseSMTPServerListener(smtpServer)

	/*
	 * Setup receivers (subscribers) to handle new mail items.
	 */
	receivers := []receiver.IMailItemReceiver{
		receiver.NewDatabaseReceiver(global.Database),
	}

	/*
	 * Start the SMTP dispatcher
	 */
	go server.Dispatch(pool, smtpServer, receivers)

	/*
	 * Pre-load layout information
	 */
	layout, err := GoHttpService.NewLayout("./www/", []string{
		"assets/mailslurper/layouts/mainLayout",
	})

	if err != nil {
		log.Printf("MailSlurper: ERROR - Error setting up layout: %s\n", err.Error())
		os.Exit(1)
	}

	/*
	 * Application context gets passed around all over the place
	 */
	appContext := &middleware.AppContext{
		Config: config,
		Layout: layout,
	}

	httpListener := listener.NewHTTPListenerService(config.WWWAddress, config.WWWPort, appContext)

	setupMiddleware(httpListener, appContext)
	setupRoutes(httpListener, appContext)

	/*
	 * Setup the app HTTP listener
	 */
	go func() {
		if err := httpListener.StartHTTPListener(config); err != nil {
			log.Printf("MailSlurper: ERROR - Error starting HTTP listener: %s\n", err.Error())
			os.Exit(1)
		}
	}()

	if config.AutoStartBrowser {
		startBrowser(config)
	}

	/*
	 * Start the services server
	 */
	serviceTierConfiguration := &configuration.ServiceTierConfiguration{
		Address:  config.ServiceAddress,
		Port:     config.ServicePort,
		Database: global.Database,
		CertFile: config.CertFile,
		KeyFile:  config.KeyFile,
	}

	if err = libmailslurper.StartServiceTier(serviceTierConfiguration); err != nil {
		log.Printf("MailSlurper: ERROR - Error starting MailSlurper services server: %s\n", err.Error())
		os.Exit(1)
	}
}

func startBrowser(config *configuration.Configuration) {
	timer := time.NewTimer(time.Second)
	go func() {
		<-timer.C
		log.Printf("Opening web browser to http://%s:%d\n", config.WWWAddress, config.WWWPort)
		err := open.Start(fmt.Sprintf("http://%s:%d", config.WWWAddress, config.WWWPort))
		if err != nil {
			log.Printf("ERROR - Could not open browser - %s\n", err.Error())
		}
	}()
}
