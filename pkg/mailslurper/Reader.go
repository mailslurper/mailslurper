// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"bytes"
	"net"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

/*
An SMTPReader is a simple object for reading commands and responses
from a connected TCP client
*/
type SMTPReader struct {
	Connection net.Conn

	logger *logrus.Entry
}

/*
The Read function reads the raw data from the socket connection to our client. This will
read on the socket until there is nothing left to read and an error is generated.
This method blocks the socket for the number of milliseconds defined in CONN_TIMEOUT_MILLISECONDS.
It then records what has been read in that time, then blocks again until there is nothing left on
the socket to read. The final value is stored and returned as a string.
*/
func (smtpReader *SMTPReader) Read() string {
	var raw bytes.Buffer
	var bytesRead int

	bytesRead = 1

	for bytesRead > 0 {
		smtpReader.Connection.SetReadDeadline(time.Now().Add(time.Millisecond * CONN_TIMEOUT_MILLISECONDS))

		buffer := make([]byte, RECEIVE_BUFFER_LEN)
		bytesRead, err := smtpReader.Connection.Read(buffer)

		if err != nil {
			break
		}

		if bytesRead > 0 {
			raw.WriteString(string(buffer[:bytesRead]))
		}
	}

	return raw.String()
}

/*
ReadDataBlock is used by the SMTP DATA command. It will read data from the connection
until the terminator is sent.
*/
func (smtpReader *SMTPReader) ReadDataBlock() string {
	var dataBuffer bytes.Buffer
	timeLimit := time.Now().Add(time.Second * COMMAND_TIMEOUT_SECONDS)

	for {
		dataResponse := smtpReader.Read()

		if len(dataResponse) > 0 {
			timeLimit = time.Now().Add(time.Second * COMMAND_TIMEOUT_SECONDS)
		}

		if time.Now().After(timeLimit) {
			break
		}

		terminatorPos := strings.Index(dataResponse, SMTP_DATA_TERMINATOR)
		if terminatorPos <= -1 {
			dataBuffer.WriteString(dataResponse)
		} else {
			dataBuffer.WriteString(dataResponse[0:terminatorPos])
			break
		}
	}

	return dataBuffer.String()
}
