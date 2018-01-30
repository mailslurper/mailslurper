// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/*
SMTPListener sets up a server that listens on a TCP socket for connections.
When a connection is received a worker is created to handle processing
the mail on this connection
*/
type SMTPListener struct {
	certificate         tls.Certificate
	config              *Configuration
	killListenerChannel chan bool
	killRecieverChannel chan bool
	listener            net.Listener
	logger              *logrus.Entry
	mailItemChannel     chan *MailItem
	receivers           []IMailItemReceiver
	serverPool          *ServerPool
}

/*
NewSMTPListener creates an SMTPListener struct
*/
func NewSMTPListener(logger *logrus.Entry, config *Configuration, serverPool *ServerPool, receivers []IMailItemReceiver) (*SMTPListener, error) {
	var err error

	result := &SMTPListener{
		config:              config,
		killListenerChannel: make(chan bool, 1),
		killRecieverChannel: make(chan bool, 1),
		logger:              logger,
		mailItemChannel:     make(chan *MailItem, 1000),
		receivers:           receivers,
		serverPool:          serverPool,
	}

	if config.CertFile != "" && config.KeyFile != "" {
		if result.certificate, err = tls.LoadX509KeyPair(config.CertFile, config.KeyFile); err != nil {
			return result, errors.Wrapf(err, "Error loading X509 certificate key pair while setting up SMTP listener")
		}
	}

	return result, nil
}

/*
Start establishes a listening connection to a socket on an address.
*/
func (s *SMTPListener) Start() error {
	var err error
	var tcpAddress *net.TCPAddr

	if s.config.CertFile != "" && s.config.KeyFile != "" {
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{s.certificate}}

		if s.listener, err = tls.Listen("tcp", s.config.GetFullSMTPBindingAddress(), tlsConfig); err != nil {
			return errors.Wrapf(err, "Unable to start SMTP listener using TLS")
		}

		s.logger.Infof("SMTP listener running on SSL %s", s.config.GetFullSMTPBindingAddress())
	} else {
		if tcpAddress, err = net.ResolveTCPAddr("tcp", s.config.GetFullSMTPBindingAddress()); err != nil {
			return errors.Wrapf(err, "Error resolving address %s starting SMTP listener", s.config.GetFullSMTPBindingAddress())
		}

		if s.listener, err = net.ListenTCP("tcp", tcpAddress); err != nil {
			return errors.Wrapf(err, "Unable to start SMTP listener")
		}

		s.logger.Infof("SMTP listener running on %s", s.config.GetFullSMTPBindingAddress())
	}

	return nil
}

/*
Dispatch starts the process of handling SMTP client connections.
The first order of business is to setup a channel for writing
parsed mails, in the form of MailItemStruct variables, to our
database. A goroutine is setup to listen on that
channel and handles storage.

Meanwhile this method will loop forever and wait for client connections (blocking).
When a connection is recieved a goroutine is started to create a new MailItemStruct
and parser and the parser process is started. If the parsing is successful
the MailItemStruct is added to a channel. An receivers passed in will be
listening on that channel and may do with the mail item as they wish.
*/
func (s *SMTPListener) Dispatch(ctx context.Context) {
	/*
	 * Setup our receivers. These guys are basically subscribers to
	 * the MailItem channel.
	 */
	var worker *SMTPWorker

	go func() {
		s.logger.Infof("%d receiver(s) listening", len(s.receivers))

		for {
			select {
			case item := <-s.mailItemChannel:
				for _, r := range s.receivers {
					go r.Receive(item)
				}

			case <-ctx.Done():
				s.logger.Infof("Shutting down receiver channel...")
				break
			}
		}
	}()

	/*
	 * Now start accepting connections for SMTP
	 */
	go func() {
		for {
			select {
			case <-ctx.Done():
				break

			default:
				connection, err := s.listener.Accept()
				if err != nil {
					s.logger.Errorf("Problem accepting SMTP requests - %s", err.Error())
					break
				}

				if worker, err = s.serverPool.NextWorker(connection, s.mailItemChannel); err != nil {
					connection.Close()
					s.logger.Errorf("Error getting next SMTP worker: %s", err.Error())
					continue
				}

				go worker.Work()
			}
		}
	}()
}
