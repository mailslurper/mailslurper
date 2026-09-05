package smtp

import (
	"bufio"
	"net"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"github.com/p0vidl0/mylslurper/internal/mail"
	"github.com/sirupsen/logrus"
)

func newTestSession(t *testing.T) (client *bufio.ReadWriter, delivered chan *mail.Item) {
	t.Helper()
	serverConn, clientConn := net.Pipe()
	delivered = make(chan *mail.Item, 1)

	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)

	s := &session{
		conn: textproto.NewConn(serverConn),
		deliver: func(item *mail.Item) error {
			delivered <- item
			return nil
		},
		log: log.WithField("test", true),
		deadline: func() {
			_ = serverConn.SetDeadline(time.Now().Add(2 * time.Second))
		},
	}
	go s.run()

	client = bufio.NewReadWriter(bufio.NewReader(clientConn), bufio.NewWriter(clientConn))
	return client, delivered
}

func readLine(t *testing.T, rw *bufio.ReadWriter) string {
	t.Helper()
	line, err := rw.ReadString('\n')
	if err != nil {
		t.Fatalf("reading line: %v", err)
	}
	return strings.TrimRight(line, "\r\n")
}

func readSMTPReply(t *testing.T, rw *bufio.ReadWriter) string {
	t.Helper()
	var lines []string
	for {
		line := readLine(t, rw)
		lines = append(lines, line)
		if len(line) >= 4 && line[3] == ' ' {
			break
		}
	}
	return lines[len(lines)-1]
}

func sendLine(t *testing.T, rw *bufio.ReadWriter, line string) {
	t.Helper()
	if _, err := rw.WriteString(line + "\r\n"); err != nil {
		t.Fatalf("writing line: %v", err)
	}
	if err := rw.Flush(); err != nil {
		t.Fatalf("flushing: %v", err)
	}
}

func TestSessionHappyPath(t *testing.T) {
	client, delivered := newTestSession(t)

	if got := readLine(t, client); !strings.HasPrefix(got, "220") {
		t.Fatalf("greeting = %q", got)
	}

	sendLine(t, client, "EHLO test.local")
	if got := readSMTPReply(t, client); !strings.HasPrefix(got, "250") {
		t.Fatalf("EHLO reply = %q", got)
	}

	sendLine(t, client, "MAIL FROM:<sender@example.com>")
	if got := readLine(t, client); !strings.HasPrefix(got, "250") {
		t.Fatalf("MAIL reply = %q", got)
	}

	sendLine(t, client, "RCPT TO:<recipient@example.com>")
	if got := readLine(t, client); !strings.HasPrefix(got, "250") {
		t.Fatalf("RCPT reply = %q", got)
	}

	sendLine(t, client, "DATA")
	if got := readLine(t, client); !strings.HasPrefix(got, "354") {
		t.Fatalf("DATA reply = %q", got)
	}

	sendLine(t, client, "Subject: Hi")
	sendLine(t, client, "")
	sendLine(t, client, "Body text")
	sendLine(t, client, ".")
	if got := readLine(t, client); !strings.HasPrefix(got, "250") {
		t.Fatalf("post-DATA reply = %q", got)
	}

	select {
	case item := <-delivered:
		if item.From != "sender@example.com" {
			t.Errorf("From = %q", item.From)
		}
		if len(item.To) != 1 || item.To[0] != "recipient@example.com" {
			t.Errorf("To = %v", item.To)
		}
		if item.Subject != "Hi" {
			t.Errorf("Subject = %q", item.Subject)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for delivered item")
	}

	sendLine(t, client, "QUIT")
	if got := readLine(t, client); !strings.HasPrefix(got, "221") {
		t.Fatalf("QUIT reply = %q", got)
	}
}

func TestSessionRejectsDataBeforeRcpt(t *testing.T) {
	client, _ := newTestSession(t)
	readLine(t, client) // greeting

	sendLine(t, client, "HELO test.local")
	readLine(t, client)

	sendLine(t, client, "MAIL FROM:<sender@example.com>")
	readLine(t, client)

	sendLine(t, client, "DATA")
	if got := readLine(t, client); !strings.HasPrefix(got, "503") {
		t.Fatalf("expected 503 for DATA before RCPT, got %q", got)
	}
}

func TestSessionRejectsMailBeforeHelo(t *testing.T) {
	client, _ := newTestSession(t)
	readLine(t, client) // greeting

	sendLine(t, client, "MAIL FROM:<sender@example.com>")
	if got := readLine(t, client); !strings.HasPrefix(got, "503") {
		t.Fatalf("expected 503 for MAIL before HELO, got %q", got)
	}
}
