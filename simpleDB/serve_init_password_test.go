package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	configpkg "github.com/aid297/aid/simpleDB/config"
)

func TestEnsureInitPassword_GeneratesAndPersistsWhenEmpty(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	cfg := configpkg.Default()
	cfg.Transport.HTTP.InitPassword = ""
	if err := configpkg.Save(configPath, cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	loaded, err := configpkg.Load(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if err = ensureInitPassword(configPath, &loaded, nil); err != nil {
		t.Fatalf("ensure init password: %v", err)
	}

	if strings.TrimSpace(loaded.Transport.HTTP.InitPassword) == "" {
		t.Fatal("expected generated initPassword in memory")
	}

	reloaded, err := configpkg.Load(configPath)
	if err != nil {
		t.Fatalf("reload config: %v", err)
	}
	if strings.TrimSpace(reloaded.Transport.HTTP.InitPassword) == "" {
		t.Fatal("expected generated initPassword persisted to config file")
	}
}

func TestEnsureInitPassword_KeepExistingPassword(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	cfg := configpkg.Default()
	cfg.Transport.HTTP.InitPassword = "already-set"
	if err := configpkg.Save(configPath, cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	loaded, err := configpkg.Load(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if err = ensureInitPassword(configPath, &loaded, nil); err != nil {
		t.Fatalf("ensure init password: %v", err)
	}
	if loaded.Transport.HTTP.InitPassword != "already-set" {
		t.Fatalf("initPassword changed unexpectedly: %s", loaded.Transport.HTTP.InitPassword)
	}
}

func TestResolveConfigPath_FallbackToLocalConfig(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := t.TempDir()
	if err = os.Chdir(dir); err != nil {
		t.Fatalf("chdir temp: %v", err)
	}
	defer func() { _ = os.Chdir(wd) }()

	if err = os.MkdirAll("simpleDB", 0o755); err != nil {
		t.Fatalf("mkdir simpleDB: %v", err)
	}
	if err = os.WriteFile("config.yaml", []byte("database:\n  path: demo\ntransport:\n  http:\n    enabled: true\n"), 0o644); err != nil {
		t.Fatalf("write config.yaml: %v", err)
	}

	resolved := resolveConfigPath("simpleDB/config.yaml")
	if resolved != "config.yaml" {
		t.Fatalf("resolved path = %s, want config.yaml", resolved)
	}
}

func TestRotateInitPasswordInConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	cfg := configpkg.Default()
	cfg.Transport.HTTP.InitPassword = "before-rotate"
	if err := configpkg.Save(configPath, cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	rotated, err := rotateInitPasswordInConfig(configPath)
	if err != nil {
		t.Fatalf("rotate init password: %v", err)
	}
	if strings.TrimSpace(rotated) == "" || rotated == "before-rotate" {
		t.Fatalf("unexpected rotated password: %s", rotated)
	}

	reloaded, err := configpkg.Load(configPath)
	if err != nil {
		t.Fatalf("reload config: %v", err)
	}
	if reloaded.Transport.HTTP.InitPassword != rotated {
		t.Fatalf("config initPassword not updated, got=%s want=%s", reloaded.Transport.HTTP.InitPassword, rotated)
	}
}
