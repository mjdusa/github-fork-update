package githubapi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/githubapi"
)

type rateLimitCategory uint8
type requestContext uint8
type values map[string]string

const (
	coreCategory         rateLimitCategory = iota
	bypassRateLimitCheck requestContext    = iota
	baseURLPath                            = "/api-v3"
	Version                                = "v53.1.0"
	defaultAPIVersion                      = "2022-11-28"
	defaultBaseURL                         = "https://api.github.com/"
	defaultUserAgent                       = "go-github" + "/" + Version
	uploadBaseURL                          = "https://uploads.github.com/"
)

func setup() (client *github.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is
	// configured to use test server.
	client = github.NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url

	return client, mux, server.URL, server.Close
}

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

func Test_WrapError_nilError(t *testing.T) {
	err := githubapi.WrapError("Test_WrapError_nilError message", nil)
	if err != nil {
		t.Errorf("githubapi.WrapError should be return nil")
	}
}

func Test_WrapError_GoodError(t *testing.T) {
	innerErrMsg := "Test_WrapError_GoodError__inner"
	innerErr := fmt.Errorf(innerErrMsg)
	wrapMsg := "Test_WrapError_GoodError message"
	err := githubapi.WrapError(wrapMsg, innerErr)
	if err == nil {
		t.Errorf("githubapi.WrapError should NOT return nil when error is not nil")
	}

	want := fmt.Errorf("%s: %w", wrapMsg, innerErr)

	if err == nil {
		t.Errorf("githubapi.WrapError returned nil")
	} else if want == nil {
		t.Errorf("githubapi.WrapError returned %+v, but want is nil", err)
	} else if err.Error() != want.Error() {
		t.Errorf("githubapi.WrapError returned %+v, want %+v", err, want)
	}
}
