package httptest_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/mjdusa/github-fork-update/internal/http/httptest"
)

func TestNewHTTPTestServerRoot(t *testing.T) {
	log := bytes.NewBufferString("")
	server, err := httptest.NewHTTPTestServer("/", log)
	if err != nil {
		t.Errorf("NewHTTPTestServer returned error: %v", err)
	}

	if server == nil {
		t.Error("Expected HTTPTestServer, got nil")
	}
	defer server.Close()

	if server.Path != "/" {
		t.Errorf("Expected server URL to be '/', got '%s'", server.Path)
	}

	if server.Mux == nil {
		t.Error("Expected Mux to be initialized, got nil")
	}

	if server.APIHandler == nil {
		t.Error("Expected APIHandler to be initialized, got nil")
	}

	if server.Server == nil {
		t.Error("Expected Server to be initialized, got nil")
	}
}

func TestNewHTTPTestServerBase(t *testing.T) {
	log := bytes.NewBufferString("")
	server, err := httptest.NewHTTPTestServer("/base", log)
	if err != nil {
		t.Errorf("NewHTTPTestServer returned error: %v", err)
	}

	if server == nil {
		t.Error("Expected HTTPTestServer, got nil")
	}
	defer server.Close()

	if server.Path != "/base" {
		t.Errorf("Expected server URL to be '/base', got '%s'", server.Path)
	}

	if server.Mux == nil {
		t.Error("Expected Mux to be initialized, got nil")
	}

	if server.APIHandler == nil {
		t.Error("Expected APIHandler to be initialized, got nil")
	}

	if server.Server == nil {
		t.Error("Expected Server to be initialized, got nil")
	}
}

func TestHTTPTestServerClose(t *testing.T) {
	log := bytes.NewBufferString("")
	server, err := httptest.NewHTTPTestServer("/base", log)
	if err != nil {
		t.Errorf("NewHTTPTestServer returned error: %v", err)
	}

	server.Close()

	req, _ := http.NewRequest("GET", server.Server.URL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		t.Error("Expected error after server close, got nil")
	}

	if resp != nil {
		t.Error("Expected no response after server close, got response")
	}
}
