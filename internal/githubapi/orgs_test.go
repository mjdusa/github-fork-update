package githubapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/githubapi"
	"github.com/mjdusa/github-fork-update/internal/http/httptest"
)

func Test_ListOrganizations_success(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	wantUser := "Test_ListOrganizations_success_user"
	wantUrl := fmt.Sprintf("/users/%s/orgs", wantUser)
	srvr.Mux.HandleFunc(wantUrl, func(wtr http.ResponseWriter, req *http.Request) {
		lo := new(github.ListOptions)
		json.NewDecoder(req.Body).Decode(lo)
		testMethod(t, req, "GET")
		testFormValues(t, req, values{
			"page": "2",
		})

		fmt.Fprint(wtr, `[{"id":1}]`)
	})

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	opts := &github.ListOptions{
		Page:    2,
		PerPage: 0,
	}

	repos, lerr := gha.ListOrganizations(ctx, wantUser, opts)
	if lerr != nil {
		t.Errorf("githubapi.ListOrganizations returned error: %v", lerr)
	}

	want := []*github.Organization{{ID: github.Int64(1)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("githubapi.ListOrganizations returned %+v, want %+v", repos, want)
	}
}

func Test_ListOrganizations_error(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	wantUser := "%"
	opts := &github.ListOptions{
		Page: 2,
	}

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	_, lerr := gha.ListOrganizations(ctx, wantUser, opts)
	if lerr == nil {
		t.Errorf("githubapi.ListOrganizations should have returned error")
	}
}
