package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTPPort != 4436 || cfg.ServicePort != 4437 || cfg.SMTPPort != 1025 {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"httpPort": 9090, "dbFile": "/tmp/x.db"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTPPort != 9090 {
		t.Fatalf("expected httpPort 9090, got %d", cfg.HTTPPort)
	}
	if cfg.DBFile != "/tmp/x.db" {
		t.Fatalf("expected overridden dbFile, got %q", cfg.DBFile)
	}
}

func TestEnvOverride(t *testing.T) {
	t.Setenv("HTTP_PORT", "1234")
	cfg, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTPPort != 1234 {
		t.Fatalf("expected env override to win, got %d", cfg.HTTPPort)
	}
}

func TestMaxWorkersAlias(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"maxWorkers": 5}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.MaxConnections != 5 {
		t.Fatalf("expected maxWorkers alias to set MaxConnections, got %d", cfg.MaxConnections)
	}
}

func TestLoadLegacyFileFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{
		"wwwAddress": "127.0.0.1",
		"wwwPort": 9001,
		"serviceAddress": "127.0.0.1",
		"servicePort": 9002,
		"smtpPort": 9003,
		"dbDatabase": "/tmp/legacy.db"
	}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTPAddress != "127.0.0.1" || cfg.HTTPPort != 9001 {
		t.Fatalf("www settings = %s:%d", cfg.HTTPAddress, cfg.HTTPPort)
	}
	if cfg.ServicePort != 9002 || cfg.SMTPPort != 9003 {
		t.Fatalf("service/smtp ports = %d/%d", cfg.ServicePort, cfg.SMTPPort)
	}
	if cfg.DBFile != "/tmp/legacy.db" {
		t.Fatalf("dbFile = %q", cfg.DBFile)
	}
}

func TestValidateRejectsBadAuthScheme(t *testing.T) {
	cfg := Default()
	cfg.AuthenticationScheme = "bogus"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid authenticationScheme")
	}
}

func TestValidateRequiresCredentialsForBasicAuth(t *testing.T) {
	cfg := Default()
	cfg.AuthenticationScheme = AuthSchemeBasic
	cfg.AuthSecret = "shh"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when credentials map is empty")
	}
}
