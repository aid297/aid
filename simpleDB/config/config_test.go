package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultsWhenFileMissing(t *testing.T) {
	config, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}
	if config.Database.Path != Default().Database.Path {
		t.Fatalf("database path = %s, want %s", config.Database.Path, Default().Database.Path)
	}
	if config.Transport.HTTP.Address != Default().Transport.HTTP.Address {
		t.Fatalf("address = %s, want %s", config.Transport.HTTP.Address, Default().Transport.HTTP.Address)
	}
	if config.Transport.HTTP.Route.SQLExecute != Default().Transport.HTTP.Route.SQLExecute {
		t.Fatalf("sql route = %s, want %s", config.Transport.HTTP.Route.SQLExecute, Default().Transport.HTTP.Route.SQLExecute)
	}
}

func TestLoad_MergesOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := Save(path, Config{Database: DatabaseConfig{Path: "custom_db"}, Transport: TransportConfig{HTTP: HTTPConfig{Enabled: true, EnableAdmin: true, EnableReport: true, Address: ":19090", TokenTTL: "2h", Route: HTTPRouteConfig{Login: "/login"}}}}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	config, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if config.Database.Path != "custom_db" {
		t.Fatalf("database path = %s, want custom_db", config.Database.Path)
	}
	if config.Transport.HTTP.Address != ":19090" {
		t.Fatalf("address = %s, want :19090", config.Transport.HTTP.Address)
	}
	if config.Transport.HTTP.Route.Refresh == "" {
		t.Fatal("expected default refresh path to be applied")
	}
	if _, err = config.ParseTokenTTL(); err != nil {
		t.Fatalf("parse token ttl: %v", err)
	}
}

func TestLoad_PreservesExplicitFalse(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config_false.yaml")
	content := []byte("database:\n  path: explicit_false\ntransport:\n  http:\n    enabled: true\n    enableAdminRoute: false\n    enableReportRoute: false\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if config.Transport.HTTP.EnableAdmin {
		t.Fatal("expected enableAdminRoute=false to be preserved")
	}
	if config.Transport.HTTP.EnableReport {
		t.Fatal("expected enableReportRoute=false to be preserved")
	}
}
