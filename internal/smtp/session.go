package smtp

import (
	"io"
	"net/textproto"
	"strings"
	"time"

	"github.com/p0vidl0/mylslurper/internal/mail"
	"github.com/sirupsen/logrus"
)

const idleTimeout = 5 * time.Minute

// Deliver is called with every fully-received mail item.
type Deliver func(item *mail.Item) error

type envelope struct {
	from string
	to   []string
}

func (e *envelope) reset() {
	e.from = ""
	e.to = nil
}

// session runs the SMTP protocol state machine for a single connection.
type session struct {
	conn     *textproto.Conn
	deliver  Deliver
	log      *logrus.Entry
	env      envelope
	greeted  bool
	authed   bool
	deadline func()
}

func (s *session) run() {
	s.deadline()
	if err := s.conn.PrintfLine("220 MylSlurper ESMTP ready"); err != nil {
		return
	}

	for {
		s.deadline()
		line, err := s.conn.ReadLine()
		if err != nil {
			if err != io.EOF {
				s.log.WithError(err).Debug("smtp: read error")
			}
			return
		}

		cmd := parseCommand(line)

		switch cmd.verb {
		case "HELO":
			s.env.reset()
			s.authed = false
			s.greeted = true
			s.reply(250, "Hello")
		case "EHLO":
			s.env.reset()
			s.authed = false
			s.greeted = true
			s.replyEhlo()
		case "AUTH":
			s.handleAuth(cmd.args)
		case "MAIL":
			s.handleMail(cmd.args)
		case "RCPT":
			s.handleRcpt(cmd.args)
		case "DATA":
			s.handleData()
		case "RSET":
			s.env.reset()
			s.reply(250, "OK")
		case "NOOP":
			s.reply(250, "OK")
		case "QUIT":
			s.reply(221, "Bye")
			return
		default:
			s.reply(500, "Unrecognized command")
		}
	}
}

func (s *session) handleMail(args string) {
	if !s.greeted {
		s.reply(503, "Send HELO/EHLO first")
		return
	}
	addr, err := parseAddress("FROM:", args)
	if err != nil {
		s.reply(501, "Malformed MAIL FROM command")
		return
	}
	s.env.reset()
	s.env.from = addr
	s.reply(250, "OK")
}

func (s *session) handleRcpt(args string) {
	if s.env.from == "" {
		s.reply(503, "Send MAIL FROM first")
		return
	}
	addr, err := parseAddress("TO:", args)
	if err != nil {
		s.reply(501, "Malformed RCPT TO command")
		return
	}
	s.env.to = append(s.env.to, addr)
	s.reply(250, "OK")
}

func (s *session) handleData() {
	if s.env.from == "" || len(s.env.to) == 0 {
		s.reply(503, "Send MAIL FROM and RCPT TO first")
		return
	}
	if err := s.conn.PrintfLine("354 Start mail input; end with <CRLF>.<CRLF>"); err != nil {
		return
	}

	s.deadline()
	raw, err := s.conn.ReadDotBytes()
	if err != nil {
		s.log.WithError(err).Debug("smtp: reading DATA body")
		s.reply(451, "Error reading message data")
		return
	}

	item, err := mail.Parse(string(raw), s.env.from, s.env.to)
	if err != nil {
		s.log.WithError(err).Warn("smtp: failed to parse message")
		s.reply(451, "Could not parse message")
		s.env.reset()
		return
	}

	if err := s.deliver(item); err != nil {
		s.log.WithError(err).Error("smtp: failed to store message")
		s.reply(451, "Could not store message")
		s.env.reset()
		return
	}

	s.reply(250, "OK: message accepted")
	s.env.reset()
}

func (s *session) replyEhlo() {
	_ = s.conn.PrintfLine("250-mylslurper ESMTP ready")
	_ = s.conn.PrintfLine("250 AUTH PLAIN LOGIN")
}

func (s *session) handleAuth(args string) {
	if !s.greeted {
		s.reply(503, "Send HELO/EHLO first")
		return
	}

	fields := strings.Fields(args)
	if len(fields) == 0 {
		s.reply(501, "AUTH requires a mechanism")
		return
	}

	mech := strings.ToUpper(fields[0])
	switch mech {
	case "PLAIN":
		if len(fields) > 1 {
			s.authed = true
			s.reply(235, "Authentication successful")
			return
		}
		if err := s.conn.PrintfLine("334 "); err != nil {
			return
		}
		if _, err := s.conn.ReadLine(); err != nil {
			return
		}
		s.authed = true
		s.reply(235, "Authentication successful")
	case "LOGIN":
		if err := s.conn.PrintfLine("334 VXNlcm5hbWU6"); err != nil {
			return
		}
		if _, err := s.conn.ReadLine(); err != nil {
			return
		}
		if err := s.conn.PrintfLine("334 UGFzc3dvcmQ6"); err != nil {
			return
		}
		if _, err := s.conn.ReadLine(); err != nil {
			return
		}
		s.authed = true
		s.reply(235, "Authentication successful")
	default:
		s.reply(504, "Unsupported authentication mechanism")
	}
}

func (s *session) reply(code int, message string) {
	_ = s.conn.PrintfLine("%d %s", code, message)
}
