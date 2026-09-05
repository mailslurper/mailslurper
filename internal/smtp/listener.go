// Package smtp implements a minimal SMTP server sufficient for capturing
// mail sent by applications under test: HELO/EHLO, MAIL FROM, RCPT TO,
// DATA, RSET, NOOP, and QUIT.
package smtp

import (
	"context"
	"crypto/tls"
	"net"
	"net/textproto"
	"time"

	"github.com/sirupsen/logrus"
)

// Listener accepts SMTP connections and hands each one to its own
// goroutine, bounded by a semaphore rather than a pre-allocated worker pool.
type Listener struct {
	Address        string
	MaxConnections int
	TLSConfig      *tls.Config
	Deliver        Deliver
	Log            *logrus.Logger

	ln  net.Listener
	sem chan struct{}
}

// Serve binds the listener and accepts connections until ctx is cancelled.
func (l *Listener) Serve(ctx context.Context) error {
	ln, err := net.Listen("tcp", l.Address)
	if err != nil {
		return err
	}
	l.ln = ln
	if l.MaxConnections <= 0 {
		l.MaxConnections = 100
	}
	l.sem = make(chan struct{}, l.MaxConnections)

	go func() {
		<-ctx.Done()
		_ = l.ln.Close()
	}()

	for {
		conn, err := l.ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			continue
		}

		if l.TLSConfig != nil {
			conn = tls.Server(conn, l.TLSConfig)
		}

		select {
		case l.sem <- struct{}{}:
			go func() {
				defer func() { <-l.sem }()
				l.handle(conn)
			}()
		case <-time.After(2 * time.Second):
			_ = conn.Close()
		}
	}
}

func (l *Listener) handle(conn net.Conn) {
	defer conn.Close()

	log := l.Log.WithField("remote", conn.RemoteAddr().String())

	s := &session{
		conn:    textproto.NewConn(conn),
		deliver: l.Deliver,
		log:     log,
		deadline: func() {
			_ = conn.SetDeadline(time.Now().Add(idleTimeout))
		},
	}
	s.run()
}
