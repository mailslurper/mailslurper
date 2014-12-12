// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/adampresley/sigint"
	"github.com/skratchdot/open-golang/open"

	"github.com/mailslurper/libmailslurper/configuration"
	"github.com/mailslurper/libmailslurper/receiver"
	"github.com/mailslurper/libmailslurper/server"
	"github.com/mailslurper/libmailslurper/storage"
	"github.com/mailslurper/libmailslurper/websocket"
	"github.com/mailslurper/mailslurper/global"
	appListener "github.com/mailslurper/mailslurper/listener"
	serviceListener "github.com/mailslurper/mailslurperservice/listener"
)

func main() {
	var err error
	runtime.GOMAXPROCS(runtime.NumCPU())

	/*
	 * Prepare SIGINT handler (CTRL+C)
	 */
	sigint.ListenForSIGINT(func() {
		log.Println("Shutting down via SIGINT...")
		os.Exit(0)
	})

	/*
	 * Load configuration
	 */
	global.Config, err = configuration.LoadConfigurationFromFile(configuration.CONFIGURATION_FILE_NAME)
	if err != nil {
		log.Println("ERROR - There was an error reading your configuration file: ", err)
		os.Exit(0)
	}

	/*
	 * Setup global database connection handle
	 */
	databaseConnection := global.Config.GetDatabaseConfiguration()

	if err = storage.ConnectToStorage(databaseConnection); err != nil {
		log.Println("ERROR - There was an error connecting to your data storage: ", err)
		os.Exit(0)
	}

	defer storage.DisconnectFromStorage()

	/*
	 * Setup the server pool
	 */
	pool := server.NewServerPool(global.Config.MaxWorkers)

	/*
	 * Setup the SMTP listener
	 */
	smtpServer, err := server.SetupSmtpServerListener(global.Config.GetFullSmtpBindingAddress())
	if err != nil {
		log.Println("ERROR - There was a problem starting the SMTP listener: ", err)
		os.Exit(0)
	}

	defer server.CloseSmtpServerListener(smtpServer)

	/*
	 * Setup receivers (subscribers) to handle new mail items.
	 */
	receivers := []receiver.IMailItemReceiver{
		receiver.DatabaseReceiver{},
		receiver.WebSocketReceiver{},
	}

	/*
	 * Start the SMTP dispatcher
	 */
	go server.Dispatcher(pool, smtpServer, receivers)

	/*
	 * Setup the app HTTP listener
	 */
	go appListener.StartHttpListener(appListener.NewHttpListener(global.Config.WWWAddress, global.Config.WWWPort))

	/*
	 * Setup web socket bucket
	 */
	websocket.WebSocketConnections = make(map[*websocket.WebSocketConnection]bool)

	/*
	 * Fire the app up in the user's default browser.
	 */
	log.Printf("INFO - Opening browser to http://%s:%d", global.Config.WWWAddress, global.Config.WWWPort)
	open.Start(fmt.Sprintf("http://%s:%d", global.Config.WWWAddress, global.Config.WWWPort))

	/*
	 * Start the services server
	 */
	serviceListener.StartHttpListener(serviceListener.NewHttpListener(global.Config.ServiceAddress, global.Config.ServicePort))
}
