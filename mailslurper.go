// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/adampresley/GoHttpService"
	"github.com/adampresley/sigint"
	"github.com/mailslurper/mailslurper/services/listener"
	"github.com/mailslurper/mailslurper/services/middleware"

	"github.com/mailslurper/libmailslurper/configuration"
	"github.com/mailslurper/libmailslurper/receiver"
	"github.com/mailslurper/libmailslurper/server"
	"github.com/mailslurper/libmailslurper/storage"
	"github.com/mailslurper/mailslurper/global"
	serviceListener "github.com/mailslurper/mailslurperservice/listener"
)

func main() {
	var err error
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Printf("SERVER - INFO - Starting MailSlurper Server v%s\n", global.SERVER_VERSION)
	/*
	 * Prepare SIGINT handler (CTRL+C)
	 */
	sigint.ListenForSIGINT(func() {
		log.Println("SERVER - INFO - Shutting down via SIGINT...")
		os.Exit(0)
	})

	/*
	 * Load configuration
	 */
	config, err := configuration.LoadConfigurationFromFile(configuration.CONFIGURATION_FILE_NAME)
	if err != nil {
		log.Println("SERVER - ERROR - There was an error reading your configuration file:", err)
		os.Exit(0)
	}

	/*
	 * Setup global database connection handle
	 */
	databaseConnection := config.GetDatabaseConfiguration()

	if err = storage.ConnectToStorage(databaseConnection); err != nil {
		log.Println("SERVER - ERROR - There was an error connecting to your data storage:", err)
		os.Exit(0)
	}

	defer storage.DisconnectFromStorage()

	/*
	 * Setup the server pool
	 */
	pool := server.NewServerPool(config.MaxWorkers)

	/*
	 * Setup the SMTP listener
	 */
	smtpServer, err := server.SetupSmtpServerListener(config.GetFullSmtpBindingAddress())
	if err != nil {
		log.Println("SERVER - ERROR - There was a problem starting the SMTP listener:", err)
		os.Exit(0)
	}

	defer server.CloseSmtpServerListener(smtpServer)

	/*
	 * Setup receivers (subscribers) to handle new mail items.
	 */
	receivers := []receiver.IMailItemReceiver{
		receiver.DatabaseReceiver{},
	}

	/*
	 * Start the SMTP dispatcher
	 */
	go server.Dispatcher(pool, smtpServer, receivers)

	/*
	 * Pre-load layout information
	 */
	layout, err := GoHttpService.NewLayout("./www/", []string{
		"assets/mailslurper/layouts/mainLayout",
	})

	if err != nil {
		log.Printf("SERVER - ERROR - Error setting up layout: %s\n", err.Error())
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
		if err := httpListener.StartHTTPListener(); err != nil {
			log.Printf("SERVER - ERROR - Error starting HTTP listener: %s\n", err.Error())
			os.Exit(1)
		}
	}()

	/*
	 * Start the services server
	 */
	err = serviceListener.StartHttpListener(serviceListener.NewHttpListener(config.ServiceAddress, config.ServicePort))

	if err != nil {
		log.Printf("SERVER - ERROR - Error starting MailSlurper services server: %s\n", err.Error())
	}
}
