package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	configpkg "github.com/aid297/aid/simpleDB/config"
	"github.com/aid297/aid/simpleDB/transport"
)

func TestNewHTTPServerFromConfig(t *testing.T) {
	config := configpkg.Default()
	config.Database.Path = filepath.Join(t.TempDir(), "service_db")
	config.Transport.HTTP.Address = ":19091"
	config.Transport.HTTP.TokenTTL = "30m"
	config.Transport.HTTP.Route.Deactivate = "/auth/deactivate-custom"
	config.Transport.HTTP.Route.SQLExecute = "/sql/custom-execute"

	server, err := newHTTPServerFromConfig(config, filepath.Join(t.TempDir(), "config.yaml"))
	if err != nil {
		t.Fatalf("new server from config: %v", err)
	}
	if server == nil || server.Engine() == nil {
		t.Fatal("expected initialized HTTP server")
	}
	registerServiceRoutes(server, config)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, config.Transport.HTTP.Route.Health, nil)
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", rec.Code, http.StatusOK)
	}

	loginRec := httptest.NewRecorder()
	loginReq := httptest.NewRequest(http.MethodPost, config.Transport.HTTP.Route.Login, http.NoBody)
	server.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusBadRequest {
		t.Fatalf("login status = %d, want %d", loginRec.Code, http.StatusBadRequest)
	}

	if server.LoginPath != transport.DefaultLoginPath && server.LoginPath != config.Transport.HTTP.Route.Login {
		t.Fatalf("unexpected login path = %s", server.LoginPath)
	}
	if server.DeactivatePath != config.Transport.HTTP.Route.Deactivate {
		t.Fatalf("unexpected deactivate path = %s", server.DeactivatePath)
	}
	if server.SQLExecutePath != config.Transport.HTTP.Route.SQLExecute {
		t.Fatalf("unexpected sql execute path = %s", server.SQLExecutePath)
	}
}
