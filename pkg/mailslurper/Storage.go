// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "github.com/sirupsen/logrus"

/*
ConnectToStorage establishes a connection to the configured database engine and returns
an object.
*/
func ConnectToStorage(storageType StorageType, connectionInfo *ConnectionInformation, logger *logrus.Entry) (IStorage, error) {
	var err error
	var storageHandle IStorage

	logger.Infof("Connecting to database")

	switch storageType {
	case STORAGE_SQLITE:
		storageHandle = NewSQLiteStorage(connectionInfo, logger)

	case STORAGE_MSSQL:
		storageHandle = NewMSSQLStorage(connectionInfo, logger)

	case STORAGE_MYSQL:
		storageHandle = NewMySQLStorage(connectionInfo, logger)
	}

	if err = storageHandle.Connect(); err != nil {
		return storageHandle, err
	}

	if err = storageHandle.Create(); err != nil {
		return storageHandle, err
	}

	return storageHandle, nil
}
