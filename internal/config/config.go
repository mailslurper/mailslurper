// Package config loads MylSlurper's runtime configuration from a JSON file,
// then applies environment variable overrides on top of it.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const (
	AuthSchemeNone  = "none"
	AuthSchemeBasic = "basic"
)

// Config holds every runtime setting for the SMTP listener, the HTTP
// server, storage, and authentication.
type Config struct {
	// Web UI (classic MailSlurper "www" listener).
	HTTPAddress string `json:"httpAddress" env:"HTTP_ADDRESS"`
	HTTPPort    int    `json:"httpPort" env:"HTTP_PORT"`
	PublicURL   string `json:"publicURL" env:"PUBLIC_URL"`

	// Service REST API (classic MailSlurper "service" listener).
	ServiceAddress string `json:"serviceAddress" env:"SERVICE_ADDRESS"`
	ServicePort    int    `json:"servicePort" env:"SERVICE_PORT"`

	SMTPAddress string `json:"smtpAddress" env:"SMTP_ADDRESS"`
	SMTPPort    int    `json:"smtpPort" env:"SMTP_PORT"`

	DBFile string `json:"dbFile" env:"DB_FILE"`

	// MaxConnections limits concurrent SMTP connections. "maxWorkers" is
	// accepted as a deprecated alias for backward compatibility with old
	// config.json files.
	MaxConnections int `json:"maxConnections" env:"MAX_CONNECTIONS"`
	MaxWorkers     int `json:"maxWorkers"`

	AutoStartBrowser bool `json:"autoStartBrowser" env:"AUTO_START_BROWSER"`

	CertFile string `json:"certFile" env:"CERT_FILE"`
	KeyFile  string `json:"keyFile" env:"KEY_FILE"`

	AuthenticationScheme string            `json:"authenticationScheme" env:"AUTH_SCHEME"`
	AuthSecret           string            `json:"authSecret" env:"AUTH_SECRET"`
	AuthTimeoutInMinutes int               `json:"authTimeoutInMinutes" env:"AUTH_TIMEOUT_MINUTES"`
	Credentials          map[string]string `json:"credentials"`

	// DevCORSOrigin, when set, allows cross-origin requests from a single
	// origin. Only meant for running the frontend from a separate dev
	// server (e.g. a live-reload server on another port); the SPA and API
	// otherwise share one origin and need no CORS at all.
	DevCORSOrigin string `json:"devCORSOrigin" env:"DEV_CORS_ORIGIN"`
}

// legacyFileConfig captures MailSlurper config.json field names so existing
// files work without changes.
type legacyFileConfig struct {
	WWWAddress     string `json:"wwwAddress"`
	WWWPort        int    `json:"wwwPort"`
	WWWPublicURL   string `json:"wwwPublicURL"`
	ServiceAddress string `json:"serviceAddress"`
	ServicePort    int    `json:"servicePort"`
	DBDatabase     string `json:"dbDatabase"`
}

// Default returns a Config populated with Ory MailSlurper-compatible defaults.
func Default() *Config {
	return &Config{
		HTTPAddress:          "0.0.0.0",
		HTTPPort:             4436,
		ServiceAddress:       "0.0.0.0",
		ServicePort:          4437,
		SMTPAddress:          "0.0.0.0",
		SMTPPort:             1025,
		DBFile:               "./mailslurper.db",
		MaxConnections:       100,
		AuthenticationScheme: AuthSchemeNone,
		AuthTimeoutInMinutes: 60,
		Credentials:          map[string]string{},
	}
}

// Load reads config from path (if it exists), applies environment variable
// overrides, and validates the result. A missing file is not an error.
func Load(path string) (*Config, error) {
	cfg := Default()

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("reading config file %q: %w", path, err)
			}
		} else {
			if err := json.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("parsing config file %q: %w", path, err)
			}
			applyLegacyFileFields(cfg, data)
		}
	}

	if cfg.MaxWorkers > 0 && cfg.MaxConnections == Default().MaxConnections {
		cfg.MaxConnections = cfg.MaxWorkers
	}

	applyEnvOverrides(cfg)
	applyLegacyEnvAliases(cfg)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func applyLegacyFileFields(cfg *Config, data []byte) {
	var legacy legacyFileConfig
	if err := json.Unmarshal(data, &legacy); err != nil {
		return
	}
	if legacy.WWWAddress != "" {
		cfg.HTTPAddress = legacy.WWWAddress
	}
	if legacy.WWWPort > 0 {
		cfg.HTTPPort = legacy.WWWPort
	}
	if legacy.WWWPublicURL != "" && cfg.PublicURL == "" {
		cfg.PublicURL = legacy.WWWPublicURL
	}
	if legacy.ServiceAddress != "" {
		cfg.ServiceAddress = legacy.ServiceAddress
	}
	if legacy.ServicePort > 0 {
		cfg.ServicePort = legacy.ServicePort
	}
	if legacy.DBDatabase != "" && cfg.DBFile == Default().DBFile {
		cfg.DBFile = legacy.DBDatabase
	}
}

func applyLegacyEnvAliases(cfg *Config) {
	if v, ok := os.LookupEnv("WWW_ADDRESS"); ok {
		cfg.HTTPAddress = v
	}
	if v, ok := os.LookupEnv("WWW_PORT"); ok {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.HTTPPort = n
		}
	}
}

// applyEnvOverrides walks the Config struct's fields and, for any field with
// an `env:"NAME"` tag whose environment variable is set, overrides the
// field's value.
func applyEnvOverrides(cfg *Config) {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envName, ok := field.Tag.Lookup("env")
		if !ok {
			continue
		}

		raw, present := os.LookupEnv(envName)
		if !present {
			continue
		}

		fv := v.Field(i)
		switch fv.Kind() {
		case reflect.String:
			fv.SetString(raw)
		case reflect.Int:
			if n, err := strconv.Atoi(raw); err == nil {
				fv.SetInt(int64(n))
			}
		case reflect.Bool:
			if b, err := strconv.ParseBool(raw); err == nil {
				fv.SetBool(b)
			}
		}
	}
}

// Validate checks that the configuration is internally consistent.
func (c *Config) Validate() error {
	if c.HTTPPort <= 0 {
		return fmt.Errorf("httpPort must be a positive number")
	}
	if c.ServicePort <= 0 {
		return fmt.Errorf("servicePort must be a positive number")
	}
	if c.SMTPPort <= 0 {
		return fmt.Errorf("smtpPort must be a positive number")
	}
	if c.DBFile == "" {
		return fmt.Errorf("dbFile must be set")
	}
	if c.MaxConnections <= 0 {
		return fmt.Errorf("maxConnections must be a positive number")
	}

	switch c.AuthenticationScheme {
	case "", AuthSchemeNone, AuthSchemeBasic:
	default:
		return fmt.Errorf("authenticationScheme must be %q or %q", AuthSchemeNone, AuthSchemeBasic)
	}

	if c.AuthenticationScheme == AuthSchemeBasic {
		if c.AuthSecret == "" {
			return fmt.Errorf("authSecret is required when authenticationScheme is %q", AuthSchemeBasic)
		}
		if len(c.Credentials) == 0 {
			return fmt.Errorf("at least one entry in credentials is required when authenticationScheme is %q", AuthSchemeBasic)
		}
	}

	if (c.CertFile == "") != (c.KeyFile == "") {
		return fmt.Errorf("certFile and keyFile must both be set, or both left blank")
	}

	return nil
}

// IsAuthEnabled reports whether requests must be authenticated.
func (c *Config) IsAuthEnabled() bool {
	return c.AuthenticationScheme == AuthSchemeBasic
}

// HTTPListenAddress returns the address the web UI server should bind to.
func (c *Config) HTTPListenAddress() string {
	return fmt.Sprintf("%s:%d", c.HTTPAddress, c.HTTPPort)
}

// ServiceListenAddress returns the address the legacy service API should bind to.
func (c *Config) ServiceListenAddress() string {
	return fmt.Sprintf("%s:%d", c.ServiceAddress, c.ServicePort)
}

// SMTPListenAddress returns the address the SMTP server should bind to.
func (c *Config) SMTPListenAddress() string {
	return fmt.Sprintf("%s:%d", c.SMTPAddress, c.SMTPPort)
}

// UsesSMTPTLS reports whether implicit TLS should be enabled on the SMTP listener.
func (c *Config) UsesSMTPTLS() bool {
	return c.CertFile != "" && c.KeyFile != ""
}
