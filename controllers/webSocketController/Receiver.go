// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package webSocketController

import (
	"net/http"

	ws "github.com/gorilla/websocket"
	"github.com/mailslurper/libmailslurper/model/mailitem"
	"github.com/mailslurper/libmailslurper/websocket"
)

/*
This function handles the handshake for our websocket connection.
It sets up a goroutine to handle sending MailItemStructs to the
other side.
*/
func WebSocketHandler(writer http.ResponseWriter, request *http.Request) {
	socket, err := ws.Upgrade(writer, request, nil, 1024, 1024)
	if _, ok := err.(ws.HandshakeError); ok {
		http.Error(writer, "Invalid handshake", 400)
		return
	} else if err != nil {
		return
	}

	/*
	 * Create a new websocket connection struct and add it's pointer
	 * address to our web socket tracking map.
	 */
	connection := &websocket.WebSocketConnection{WS: socket, SendChannel: make(chan mailitem.MailItem, 256)}
	websocket.ActivateSocket(connection)
	defer websocket.DestroyConnection(connection)

	for {
		for message := range connection.SendChannel {
			err := connection.WS.WriteJSON(message)
			if err != nil {
				break
			}
		}
	}

	connection.WS.Close()
}
