// Command mylslurper runs the SMTP capture server and its web UI/API.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/p0vidl0/mylslurper/internal/api"
	"github.com/p0vidl0/mylslurper/internal/auth"
	"github.com/p0vidl0/mylslurper/internal/config"
	"github.com/p0vidl0/mylslurper/internal/events"
	"github.com/p0vidl0/mylslurper/internal/logging"
	"github.com/p0vidl0/mylslurper/internal/mail"
	"github.com/p0vidl0/mylslurper/internal/serviceapi"
	"github.com/p0vidl0/mylslurper/internal/smtp"
	"github.com/p0vidl0/mylslurper/internal/storage"
	"github.com/p0vidl0/mylslurper/web"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config.json")
	logLevel := flag.String("loglevel", "info", "log level: debug, info, warn, error")
	logFormat := flag.String("logformat", "text", "log format: text or json")
	flag.Parse()

	log := logging.New(*logLevel, *logFormat)

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.WithError(err).Fatal("failed to load configuration")
	}

	store := storage.NewSQLiteStorage(cfg.DBFile)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := store.Connect(ctx); err != nil {
		log.WithError(err).Fatal("failed to connect to storage")
	}
	defer store.Close()

	jwtSvc := auth.NewJWTService(cfg.AuthSecret, time.Duration(cfg.AuthTimeoutInMinutes)*time.Minute)
	sessionSvc := auth.NewSessionService(cfg.AuthSecret, time.Duration(cfg.AuthTimeoutInMinutes)*time.Minute)
	eventHub := events.NewHub()

	assets, err := web.Assets()
	if err != nil {
		log.WithError(err).Fatal("failed to load embedded frontend assets")
	}

	apiServer := &api.API{
		Store:    store,
		Config:   cfg,
		JWT:      jwtSvc,
		Sessions: sessionSvc,
		Log:      log,
		Assets:   assets,
		Events:   eventHub,
	}

	wwwServer := &http.Server{
		Addr:    cfg.HTTPListenAddress(),
		Handler: apiServer.Router(),
	}

	serviceServer := &http.Server{
		Addr:    cfg.ServiceListenAddress(),
		Handler: (&serviceapi.Server{Store: store, Log: log}).Router(),
	}

	smtpListener := &smtp.Listener{
		Address:        cfg.SMTPListenAddress(),
		MaxConnections: cfg.MaxConnections,
		Log:            log,
		Deliver: func(item *mail.Item) error {
			if _, err := store.StoreMail(context.Background(), item); err != nil {
				return err
			}
			eventHub.Publish(events.MailReceived(item))
			return nil
		},
	}
	if cfg.UsesSMTPTLS() {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			log.WithError(err).Fatal("failed to load TLS certificate")
		}
		smtpListener.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	go func() {
		if err := smtpListener.Serve(ctx); err != nil {
			log.WithError(err).Fatal("smtp listener stopped")
		}
	}()

	go func() {
		log.WithField("address", wwwServer.Addr).Info("web ui listening")
		if err := wwwServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("web ui stopped")
		}
	}()

	go func() {
		log.WithField("address", serviceServer.Addr).Info("service api listening")
		if err := serviceServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("service api stopped")
		}
	}()

	if cfg.AutoStartBrowser {
		go func() {
			time.Sleep(500 * time.Millisecond)
			_ = open.Start(fmt.Sprintf("http://%s", wwwServer.Addr))
		}()
	}

	log.WithField("smtpAddress", smtpListener.Address).Info("smtp listener listening")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Info("shutting down")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer shutdownCancel()
	if err := wwwServer.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("error during web ui shutdown")
	}
	if err := serviceServer.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("error during service api shutdown")
	}
}
