package httptest_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/mjdusa/github-fork-update/internal/http/httptest"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPTestServerRoot(t *testing.T) {
	path := "/"
	log := bytes.NewBufferString("")
	server, err := httptest.NewHTTPTestServer(path, log)
	assert.NoError(t, err, "NewHTTPTestServer returned error: %v", err)
	assert.NotNil(t, server, "Expected HTTPTestServer, got nil")
	defer server.Close()
	assert.Equal(t, path, server.Path, "Expected server URL to be '/', got '%s'", server.Path)
	assert.NotNil(t, server.Mux, "Expected Mux to be initialized, got nil")
	assert.NotNil(t, server.APIHandler, "Expected APIHandler to be initialized, got nil")
	assert.NotNil(t, server.Server, "Expected Server to be initialized, got nil")
}

func TestNewHTTPTestServerBase(t *testing.T) {
	path := "/base"
	log := bytes.NewBufferString("")
	server, err := httptest.NewHTTPTestServer(path, log)
	assert.NoError(t, err, "NewHTTPTestServer returned error: %v", err)
	assert.NotNil(t, server, "Expected HTTPTestServer, got nil")
	defer server.Close()
	assert.Equal(t, path, server.Path, "Expected server URL to be '/', got '%s'", server.Path)
	assert.NotNil(t, server.Mux, "Expected Mux to be initialized, got nil")
	assert.NotNil(t, server.APIHandler, "Expected APIHandler to be initialized, got nil")
	assert.NotNil(t, server.Server, "Expected Server to be initialized, got nil")
}

func TestHTTPTestServerClose(t *testing.T) {
	path := "/base"
	log := bytes.NewBufferString("")
	server, err := httptest.NewHTTPTestServer(path, log)
	assert.NoError(t, err, "NewHTTPTestServer returned error: %v", err)
	assert.NotNil(t, server, "Expected HTTPTestServer, got nil")
	server.Close()

	req, _ := http.NewRequest(http.MethodGet, server.Server.URL, nil)
	resp, err := http.DefaultClient.Do(req)
	assert.Error(t, err, "Expected error after server close, got nil")
	assert.Nil(t, resp, "Expected no response after server close, got response")
}
