// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "net"

/*
ConnectionPool is a map of remote address to TCP connections and their workers
*/
type ConnectionPool map[string]*ConnectionPoolItem

/*
ConnectionPoolItem is a single item in the pool. This tracks a connection
to its worker
*/
type ConnectionPoolItem struct {
	Connection net.Conn
	Worker     *SMTPWorker
}

/*
NewConnectionPool creates a new empty map
*/
func NewConnectionPool() ConnectionPool {
	return make(ConnectionPool)
}

/*
NewConnectionPoolItem create a new object
*/
func NewConnectionPoolItem(connection net.Conn, worker *SMTPWorker) *ConnectionPoolItem {
	return &ConnectionPoolItem{
		Connection: connection,
		Worker:     worker,
	}
}
