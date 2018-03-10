// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"bytes"
	"context"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/*
An SMTPReader is a simple object for reading commands and responses
from a connected TCP client
*/
type SMTPReader struct {
	Connection net.Conn

	logger            *logrus.Entry
	killServerContext context.Context
}

/*
The Read function reads the raw data from the socket connection to our client. This will
read on the socket until there is nothing left to read and an error is generated.
This method blocks the socket for the number of milliseconds defined in CONN_TIMEOUT_MILLISECONDS.
It then records what has been read in that time, then blocks again until there is nothing left on
the socket to read. The final value is stored and returned as a string.
*/
func (smtpReader *SMTPReader) Read() (string, error) {
	var raw bytes.Buffer
	var bytesRead int
	var err error

	bytesRead = 1

	for bytesRead > 0 {
		select {
		case <-smtpReader.killServerContext.Done():
			return "", nil

		default:
			if err = smtpReader.Connection.SetReadDeadline(time.Now().Add(time.Minute * CONNECTION_TIMEOUT_MINUTES)); err != nil {
				return raw.String(), nil
			}

			buffer := make([]byte, RECEIVE_BUFFER_LEN)
			bytesRead, err = smtpReader.Connection.Read(buffer)

			if err != nil {
				return raw.String(), err
			}

			if bytesRead > 0 {
				raw.WriteString(string(buffer[:bytesRead]))
				if strings.HasSuffix(raw.String(), "\r\n") {
					return raw.String(), nil
				}
			}
		}
	}

	return raw.String(), nil
}

/*
ReadDataBlock is used by the SMTP DATA command. It will read data from the connection
until the terminator is sent.
*/
func (smtpReader *SMTPReader) ReadDataBlock() (string, error) {
	var dataBuffer bytes.Buffer

	for {
		dataResponse, err := smtpReader.Read()
		if err != nil {
			smtpReader.logger.WithError(err).Errorf("Error reading in DATA block")
			return dataBuffer.String(), errors.Wrapf(err, "Error reading in DATA block")
		}

		dataBuffer.WriteString(dataResponse)
		terminatorPos := strings.Index(dataBuffer.String(), SMTP_DATA_TERMINATOR)

		if terminatorPos > -1 {
			break
		}
	}

	result := dataBuffer.String()
	result = result[:len(result)-3]

	return result, nil
}
