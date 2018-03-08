package main

import (
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
)

func setupDatabase() {
	/*
	 * Setup global database connection handle
	 */
	storageType, databaseConnection := config.GetDatabaseConfiguration()

	if database, err = mailslurper.ConnectToStorage(storageType, databaseConnection, logger); err != nil {
		logger.WithError(err).Fatalf("Error connecting to storage type '%d' with a connection string of %s", int(storageType), databaseConnection.String())
	}

	defer database.Disconnect()
}
