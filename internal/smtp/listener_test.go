package smtp

import (
	"bufio"
	"context"
	"encoding/base64"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/p0vidl0/mylslurper/internal/mail"
	"github.com/sirupsen/logrus"
)

func reserveAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	addr := ln.Addr().String()
	if err := ln.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	return addr
}

func dialSMTP(t *testing.T, addr string) *bufio.ReadWriter {
	t.Helper()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
}

func startTestListener(t *testing.T, deliver Deliver, maxConnections int) (*Listener, context.CancelFunc) {
	t.Helper()
	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)

	ctx, cancel := context.WithCancel(context.Background())
	listener := &Listener{
		Address:        reserveAddr(t),
		MaxConnections: maxConnections,
		Deliver:        deliver,
		Log:            log,
	}

	done := make(chan error, 1)
	go func() { done <- listener.Serve(ctx) }()

	deadline := time.Now().Add(2 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", listener.Address, 50*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for SMTP listener")
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Cleanup(func() {
		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("Serve returned error: %v", err)
			}
		case <-time.After(3 * time.Second):
			t.Error("listener did not shut down in time")
		}
	})

	return listener, cancel
}

func smtpHappyPath(t *testing.T, client *bufio.ReadWriter, subject string) {
	t.Helper()
	readLine(t, client) // greeting
	sendLine(t, client, "EHLO test.local")
	readSMTPReply(t, client)
	sendLine(t, client, "MAIL FROM:<sender@example.com>")
	readLine(t, client)
	sendLine(t, client, "RCPT TO:<recipient@example.com>")
	readLine(t, client)
	sendLine(t, client, "DATA")
	readLine(t, client)
	sendLine(t, client, "Subject: "+subject)
	sendLine(t, client, "")
	sendLine(t, client, "Hello from listener")
	sendLine(t, client, ".")
	readLine(t, client)
	sendLine(t, client, "QUIT")
	readLine(t, client)
}

func TestListenerAcceptsAndDelivers(t *testing.T) {
	delivered := make(chan *mail.Item, 1)
	listener, _ := startTestListener(t, func(item *mail.Item) error {
		delivered <- item
		return nil
	}, 10)

	client := dialSMTP(t, listener.Address)
	smtpHappyPath(t, client, "Listener test")

	select {
	case item := <-delivered:
		if item.Subject != "Listener test" {
			t.Errorf("Subject = %q", item.Subject)
		}
		if item.From != "sender@example.com" {
			t.Errorf("From = %q", item.From)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for delivered item")
	}
}

func TestListenerShutdownOnCancel(t *testing.T) {
	delivered := make(chan *mail.Item, 1)
	listener, cancel := startTestListener(t, func(item *mail.Item) error {
		delivered <- item
		return nil
	}, 10)

	client := dialSMTP(t, listener.Address)
	smtpHappyPath(t, client, "Shutdown test")

	select {
	case <-delivered:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for delivered item")
	}

	cancel()
}

func TestListenerMaxConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow max-connections test in short mode")
	}

	block := make(chan struct{})
	listener, _ := startTestListener(t, func(item *mail.Item) error {
		<-block
		return nil
	}, 1)

	hold := dialSMTP(t, listener.Address)
	readLine(t, hold) // greeting
	sendLine(t, hold, "EHLO test.local")
	readSMTPReply(t, hold)
	sendLine(t, hold, "MAIL FROM:<a@b.com>")
	readLine(t, hold)
	sendLine(t, hold, "RCPT TO:<c@d.com>")
	readLine(t, hold)
	sendLine(t, hold, "DATA")
	readLine(t, hold)

	second, err := net.DialTimeout("tcp", listener.Address, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("second dial: %v", err)
	}
	defer second.Close()
	secondRW := bufio.NewReadWriter(bufio.NewReader(second), bufio.NewWriter(second))

	if _, err := secondRW.ReadString('\n'); err == nil {
		t.Fatal("expected second connection to be closed when pool is full")
	}

	close(block)
	sendLine(t, hold, "Subject: hold")
	sendLine(t, hold, "")
	sendLine(t, hold, "body")
	sendLine(t, hold, ".")
	readLine(t, hold)
	sendLine(t, hold, "QUIT")
	readLine(t, hold)
}

func TestSessionAuthPlainInline(t *testing.T) {
	client, _ := newTestSession(t)
	readLine(t, client)
	sendLine(t, client, "EHLO test.local")
	readSMTPReply(t, client)

	creds := base64.StdEncoding.EncodeToString([]byte("\x00user\x00pass"))
	sendLine(t, client, "AUTH PLAIN "+creds)
	if got := readLine(t, client); got[:3] != "235" {
		t.Fatalf("AUTH PLAIN reply = %q", got)
	}
}

func TestSessionAuthPlainChallenge(t *testing.T) {
	client, _ := newTestSession(t)
	readLine(t, client)
	sendLine(t, client, "EHLO test.local")
	readSMTPReply(t, client)

	sendLine(t, client, "AUTH PLAIN")
	if got := readLine(t, client); got[:3] != "334" {
		t.Fatalf("AUTH PLAIN challenge = %q", got)
	}

	creds := base64.StdEncoding.EncodeToString([]byte("\x00user\x00pass"))
	sendLine(t, client, creds)
	if got := readLine(t, client); got[:3] != "235" {
		t.Fatalf("AUTH PLAIN reply = %q", got)
	}
}

func TestSessionAuthLogin(t *testing.T) {
	client, _ := newTestSession(t)
	readLine(t, client)
	sendLine(t, client, "EHLO test.local")
	readSMTPReply(t, client)

	sendLine(t, client, "AUTH LOGIN")
	if got := readLine(t, client); got[:3] != "334" {
		t.Fatalf("username challenge = %q", got)
	}
	sendLine(t, client, base64.StdEncoding.EncodeToString([]byte("user")))
	if got := readLine(t, client); got[:3] != "334" {
		t.Fatalf("password challenge = %q", got)
	}
	sendLine(t, client, base64.StdEncoding.EncodeToString([]byte("pass")))
	if got := readLine(t, client); got[:3] != "235" {
		t.Fatalf("AUTH LOGIN reply = %q", got)
	}
}

func TestSessionAuthUnsupported(t *testing.T) {
	client, _ := newTestSession(t)
	readLine(t, client)
	sendLine(t, client, "EHLO test.local")
	readSMTPReply(t, client)

	sendLine(t, client, "AUTH CRAM-MD5")
	if got := readLine(t, client); got[:3] != "504" {
		t.Fatalf("unsupported AUTH reply = %q", got)
	}
}

func TestSessionAuthBeforeHelo(t *testing.T) {
	client, _ := newTestSession(t)
	readLine(t, client)

	sendLine(t, client, "AUTH PLAIN dXNlcg==")
	if got := readLine(t, client); got[:3] != "503" {
		t.Fatalf("AUTH before HELO reply = %q", got)
	}
}

func TestSessionRset(t *testing.T) {
	client, _ := newTestSession(t)
	readLine(t, client)
	sendLine(t, client, "HELO test.local")
	readLine(t, client)
	sendLine(t, client, "MAIL FROM:<sender@example.com>")
	readLine(t, client)

	sendLine(t, client, "RSET")
	if got := readLine(t, client); got[:3] != "250" {
		t.Fatalf("RSET reply = %q", got)
	}

	sendLine(t, client, "DATA")
	if got := readLine(t, client); got[:3] != "503" {
		t.Fatalf("expected 503 for DATA after RSET, got %q", got)
	}
}

func TestSessionDotStuffing(t *testing.T) {
	client, delivered := newTestSession(t)
	readLine(t, client)
	sendLine(t, client, "EHLO test.local")
	readSMTPReply(t, client)
	sendLine(t, client, "MAIL FROM:<sender@example.com>")
	readLine(t, client)
	sendLine(t, client, "RCPT TO:<recipient@example.com>")
	readLine(t, client)
	sendLine(t, client, "DATA")
	readLine(t, client)
	sendLine(t, client, "Subject: Dots")
	sendLine(t, client, "")
	sendLine(t, client, "..not a terminator")
	sendLine(t, client, ".")
	readLine(t, client)

	select {
	case item := <-delivered:
		if strings.TrimSpace(item.TextBody) != ".not a terminator" {
			t.Errorf("TextBody = %q", item.TextBody)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for delivered item")
	}
}

func TestSessionMultipleRcpt(t *testing.T) {
	client, delivered := newTestSession(t)
	readLine(t, client)
	sendLine(t, client, "EHLO test.local")
	readSMTPReply(t, client)
	sendLine(t, client, "MAIL FROM:<sender@example.com>")
	readLine(t, client)
	sendLine(t, client, "RCPT TO:<one@example.com>")
	readLine(t, client)
	sendLine(t, client, "RCPT TO:<two@example.com>")
	readLine(t, client)
	sendLine(t, client, "DATA")
	readLine(t, client)
	sendLine(t, client, "Subject: Multi")
	sendLine(t, client, "")
	sendLine(t, client, "Both")
	sendLine(t, client, ".")
	readLine(t, client)

	select {
	case item := <-delivered:
		if len(item.To) != 2 {
			t.Fatalf("To = %v", item.To)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for delivered item")
	}
}
