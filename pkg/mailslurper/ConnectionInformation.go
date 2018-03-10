// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "fmt"

/*
ConnectionInformation contains data necessary to establish a connection
to a database server.
*/
type ConnectionInformation struct {
	Address  string
	Port     int
	Database string
	UserName string
	Password string
	Filename string
}

/*
NewConnectionInformation returns a new ConnectionInformation structure with
the address and port filled in.
*/
func NewConnectionInformation(address string, port int) *ConnectionInformation {
	return &ConnectionInformation{
		Address: address,
		Port:    port,
	}
}

/*
SetDatabaseInformation fills in the name of a database to connect to, and the user
credentials necessary to do so
*/
func (information *ConnectionInformation) SetDatabaseInformation(database, userName, password string) {
	information.Database = database
	information.UserName = userName
	information.Password = password
}

/*
SetDatabaseFile sets the name of a file-base database. This is used for SQLite
*/
func (information *ConnectionInformation) SetDatabaseFile(filename string) {
	information.Filename = filename
}

func (information *ConnectionInformation) String() string {
	if information.Filename != "" {
		return information.Filename
	}

	return fmt.Sprintf("%s:%s@%s:%d/%s", information.UserName, information.Password, information.Address, information.Port, information.Database)
}
