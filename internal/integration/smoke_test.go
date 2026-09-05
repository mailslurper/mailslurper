package integration

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/p0vidl0/mylslurper/internal/api"
	"github.com/p0vidl0/mylslurper/internal/config"
	"github.com/p0vidl0/mylslurper/internal/mail"
	"github.com/p0vidl0/mylslurper/internal/smtp"
	"github.com/p0vidl0/mylslurper/internal/storage"
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

func sendLine(t *testing.T, rw *bufio.ReadWriter, line string) {
	t.Helper()
	if _, err := rw.WriteString(line + "\r\n"); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := rw.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
}

func readLine(t *testing.T, rw *bufio.ReadWriter) string {
	t.Helper()
	line, err := rw.ReadString('\n')
	if err != nil {
		t.Fatalf("read: %v", err)
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
	return strings.Join(lines, "\n")
}

func TestSMTPToAPI(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := storage.NewSQLiteStorage(filepath.Join(t.TempDir(), "integration.db"))
	if err := store.Connect(ctx); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)

	listener := &smtp.Listener{
		Address: reserveAddr(t),
		Deliver: func(item *mail.Item) error {
			_, err := store.StoreMail(context.Background(), item)
			return err
		},
		Log: log,
	}

	done := make(chan error, 1)
	go func() { done <- listener.Serve(ctx) }()
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Error("smtp listener did not stop")
		}
	})

	deadline := time.Now().Add(2 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", listener.Address, 50*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for smtp listener")
		}
		time.Sleep(10 * time.Millisecond)
	}

	conn, err := net.Dial("tcp", listener.Address)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	client := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	readLine(t, client)
	sendLine(t, client, "EHLO integration.test")
	readSMTPReply(t, client)
	sendLine(t, client, "MAIL FROM:<sender@example.com>")
	readLine(t, client)
	sendLine(t, client, "RCPT TO:<recipient@example.com>")
	readLine(t, client)
	sendLine(t, client, "DATA")
	readLine(t, client)
	sendLine(t, client, "Subject: Integration smoke")
	sendLine(t, client, "")
	sendLine(t, client, "End-to-end body")
	sendLine(t, client, ".")
	readLine(t, client)
	sendLine(t, client, "QUIT")
	readLine(t, client)

	apiServer := &api.API{
		Store:  store,
		Config: config.Default(),
		Log:    log,
	}
	httpServer := httptest.NewServer(apiServer.Router())
	t.Cleanup(httpServer.Close)

	resp, err := http.Get(httpServer.URL + "/api/mail")
	if err != nil {
		t.Fatalf("GET /api/mail: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/mail = %d", resp.StatusCode)
	}

	var list struct {
		TotalRecords int `json:"totalRecords"`
		MailItems    []struct {
			Subject string `json:"subject"`
		} `json:"mailItems"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if list.TotalRecords != 1 || len(list.MailItems) != 1 {
		t.Fatalf("unexpected list: %+v", list)
	}
	if list.MailItems[0].Subject != "Integration smoke" {
		t.Fatalf("Subject = %q", list.MailItems[0].Subject)
	}
}
