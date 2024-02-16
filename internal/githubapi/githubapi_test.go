package githubapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/githubapi"
	"github.com/mjdusa/github-fork-update/internal/http/httptest"
)

func NewGitHubAPITokenClient(ctx context.Context, token string, apiUrl string) (*github.Client, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("empty token error")
	}

	url, err := url.Parse(apiUrl + githubapi.GitHubAPIBaseURLPath + "/")
	if err != nil {
		return nil, fmt.Errorf("url.Parse returned error: %w", err)
	}

	client := github.NewTokenClient(ctx, token)
	client.BaseURL = url
	client.UploadURL = url

	return client, nil
}

func NewTestGitHubAPI(ctx context.Context, auth string, url string) (*githubapi.GitHubAPI, error) {
	if len(auth) == 0 {
		return nil, fmt.Errorf("empty token error")
	}

	if len(url) == 0 {
		return nil, fmt.Errorf("empty url error")
	}

	client, err := NewGitHubAPITokenClient(ctx, auth, url)
	if err != nil {
		return nil, fmt.Errorf("NewGitHubAPITokenClient returned error: %w", err)
	}

	if client == nil {
		return nil, fmt.Errorf("NewGitHubAPITokenClient returned nil")
	}

	api := githubapi.GitHubAPI{
		Client: client,
	}

	return &api, nil
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	t.Helper()
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	r.ParseForm()
	if got := r.Form; !cmp.Equal(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestNewGitHubAPITokenClient(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer("path", os.Stderr)
	if serr != nil {
		t.Errorf("NewHTTPTestServer returned error: %v", serr)
	}
	defer srvr.Close()

	ctx := context.Background()
	token := "testToken"

	client, err := NewGitHubAPITokenClient(ctx, token, srvr.Server.URL)
	if err != nil {
		t.Errorf("NewGitHubAPITokenClient returned error: %v", err)
	}

	if client == nil {
		t.Error("Expected GitHub client, got nil")
	}
}

func TestNewGitHubAPI(t *testing.T) {
	ctx := context.Background()
	api, err := githubapi.NewGitHubAPI(ctx, "testAuth")
	if err != nil {
		t.Errorf("NewGitHubAPI returned error: %v", err)
	}

	if api == nil {
		t.Errorf("NewGitHubAPI() api = nil")
	}
}

func TestNewGitHubAPIerror(t *testing.T) {
	ctx := context.Background()
	api, err := githubapi.NewGitHubAPI(ctx, "")
	if err == nil {
		t.Errorf("Expected error, but got none")
	}
	if api != nil {
		t.Errorf("Expected api to be nil, but got %v", api)
	}
}
