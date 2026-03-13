package transport

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestChooseResponseFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		accept string
		want   httpBodyFormat
	}{
		{name: "default", accept: "", want: httpBodyFormatJSON},
		{name: "wildcard", accept: "*/*", want: httpBodyFormatJSON},
		{name: "xml", accept: "application/xml", want: httpBodyFormatXML},
		{name: "yaml", accept: "application/yaml", want: httpBodyFormatYAML},
		{name: "toml", accept: "application/toml", want: httpBodyFormatTOML},
		{name: "q_value", accept: "application/xml;q=0.1, application/json;q=0.9", want: httpBodyFormatJSON},
		{name: "suffix_json", accept: "application/vnd.api+json", want: httpBodyFormatJSON},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := chooseResponseFormat(tt.accept); got != tt.want {
				t.Fatalf("chooseResponseFormat(%q) = %q, want %q", tt.accept, got, tt.want)
			}
		})
	}
}

func TestBindRequestBody_DefaultsToJSONWhenMissingContentType(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	type Req struct {
		SQL string `json:"sql"`
	}

	r.POST("/echo", func(c *gin.Context) {
		var req Req
		if err := bindRequestBody(c, &req); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(c, http.StatusOK, req)
	})

	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(`{"sql":"select 1"}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want json", ct)
	}
	if !strings.Contains(w.Body.String(), `"sql":"select 1"`) {
		t.Fatalf("body = %q", w.Body.String())
	}
}

func TestHTTPResponse_NegotiatesByAccept(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	type Resp struct {
		SQL string `json:"sql"`
	}

	r.GET("/resp", func(c *gin.Context) {
		writeJSON(c, http.StatusOK, Resp{SQL: "select 1"})
	})

	tests := []struct {
		name           string
		accept         string
		wantCTContains string
	}{
		{name: "json_default", accept: "", wantCTContains: "application/json"},
		{name: "xml", accept: "application/xml", wantCTContains: "application/xml"},
		{name: "yaml", accept: "application/yaml", wantCTContains: "yaml"},
		{name: "toml", accept: "application/toml", wantCTContains: "application/toml"},
		{name: "q_json_wins", accept: "application/xml;q=0.1, application/json;q=0.9", wantCTContains: "application/json"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/resp", nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
			ct := w.Header().Get("Content-Type")
			if !strings.Contains(strings.ToLower(ct), strings.ToLower(tt.wantCTContains)) {
				t.Fatalf("Content-Type = %q, want contains %q", ct, tt.wantCTContains)
			}
			if w.Header().Get("Vary") != "Accept" {
				t.Fatalf("Vary = %q, want %q", w.Header().Get("Vary"), "Accept")
			}
		})
	}
}

func TestBindRequestBody_ByContentType(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	type Req struct {
		SQL string `json:"sql"`
	}

	r.POST("/bind", func(c *gin.Context) {
		var req Req
		if err := bindRequestBody(c, &req); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(c, http.StatusOK, req)
	})

	tests := []struct {
		name        string
		contentType string
		body        string
	}{
		{name: "json", contentType: "application/json", body: `{"sql":"select 1"}`},
		{name: "yaml", contentType: "application/yaml", body: "sql: select 1\n"},
		{name: "toml", contentType: "application/toml", body: `sql = "select 1"` + "\n"},
		{name: "xml", contentType: "application/xml", body: `<Req><SQL>select 1</SQL></Req>`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(tt.body))
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), `"sql":"select 1"`) {
				t.Fatalf("body = %q", w.Body.String())
			}
		})
	}
}
