// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
)

func setupDatabase() {
	var err error

	/*
	 * Setup global database connection handle
	 */
	storageType, databaseConnection := config.GetDatabaseConfiguration()
	xssService := mailslurper.NewXSSService(config)

	if database, err = mailslurper.ConnectToStorage(storageType, databaseConnection, xssService, logger); err != nil {
		logger.WithError(err).Fatalf("Error connecting to storage type '%d' with a connection string of %s", int(storageType), databaseConnection.String())
	}
}
