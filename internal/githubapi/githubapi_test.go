package githubapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/githubapi"
	"github.com/mjdusa/github-fork-update/internal/http/httptest"
	"github.com/stretchr/testify/assert"
)

const HTTPHeaderKeyAccept = "Accept"

func NewGitHubAPITokenClient(ctx context.Context, token string, apiURL string) (*github.Client, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("empty token error")
	}

	url, err := url.Parse(apiURL + githubapi.GitHubAPIBaseURLPath + "/")
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

func testFormValues(t *testing.T, req *http.Request, values values) {
	t.Helper()
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	err := req.ParseForm()
	assert.NoError(t, err, "ParseForm() returned error: %w", err)

	if got := req.Form; !cmp.Equal(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, req *http.Request, want string) {
	t.Helper()
	if got := req.Header.Get(HTTPHeaderKeyAccept); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", HTTPHeaderKeyAccept, got, want)
	}
}

func testMethod(t *testing.T, req *http.Request, want string) {
	t.Helper()
	if got := req.Method; got != want {
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

func TestListOrganizationsSuccess(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	wantUser := "Test_ListOrganizations_success_user"
	wantURL := fmt.Sprintf("/users/%s/orgs", wantUser)
	srvr.Mux.HandleFunc(wantURL, func(wtr http.ResponseWriter, req *http.Request) {
		lo := new(github.ListOptions)
		json.NewDecoder(req.Body).Decode(lo) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodGet)
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

func TestListOrganizationsError(t *testing.T) {
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

func TestListRepositoriesSuccess(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	wantUser := "Test_ListRepositories_success_user"
	wantURL := fmt.Sprintf("/users/%s/repos", wantUser)
	wantAcceptHeaders := []string{"application/vnd.github.mercy-preview+json",
		"application/vnd.github.nebula-preview+json"}
	srvr.Mux.HandleFunc(wantURL, func(wtr http.ResponseWriter, req *http.Request) {
		rlo := new(github.RepositoryListOptions)
		json.NewDecoder(req.Body).Decode(rlo) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodGet)
		testHeader(t, req, strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, req, values{
			"visibility":  "public",
			"affiliation": "owner,collaborator",
			"sort":        "created",
			"direction":   "asc",
			"page":        "2",
		})

		fmt.Fprint(wtr, `[{"id":1}]`)
	})

	opt := &github.RepositoryListOptions{
		Visibility:  "public",
		Affiliation: "owner,collaborator",
		Sort:        "created",
		Direction:   "asc",
		ListOptions: github.ListOptions{Page: 2},
	}

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	repos, err := gha.ListRepositories(ctx, wantUser, opt)
	if err != nil {
		t.Errorf("githubapi.ListRepositories returned error: %v", err)
	}

	want := []*github.Repository{{ID: github.Int64(1)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("githubapi.ListRepositories returned %+v, want %+v", repos, want)
	}
}

func TestListRepositoriesError(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	wantUser := "Test_ListRepositories_error_user"
	wantURL := fmt.Sprintf("/users/%s/repos", wantUser)
	wantAcceptHeaders := []string{"application/vnd.github.mercy-preview+json",
		"application/vnd.github.nebula-preview+json"}
	srvr.Mux.HandleFunc(wantURL, func(wtr http.ResponseWriter, req *http.Request) {
		rlo := new(github.RepositoryListOptions)
		json.NewDecoder(req.Body).Decode(rlo) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodGet)
		testHeader(t, req, strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, req, values{
			"visibility":  "public",
			"affiliation": "owner,collaborator",
			"sort":        "created",
			"direction":   "asc",
			"page":        "2",
		})

		wtr.WriteHeader(http.StatusUnprocessableEntity)
	})

	opt := &github.RepositoryListOptions{
		Visibility:  "public",
		Affiliation: "owner,collaborator",
		Sort:        "created",
		Direction:   "asc",
		ListOptions: github.ListOptions{Page: 2},
	}

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	_, err := gha.ListRepositories(ctx, wantUser, opt)
	if err == nil {
		t.Errorf("githubapi.ListRepositories should have returned an error")
	}
}

func TestListForksSuccess(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	wantOwner := "wantOwner"
	wantRepo := "wantRepo"
	wantURL := fmt.Sprintf("/repos/%s/%s/forks", wantOwner, wantRepo)
	srvr.Mux.HandleFunc(wantURL, func(wtr http.ResponseWriter, req *http.Request) {
		rlfo := new(github.RepositoryListForksOptions)
		json.NewDecoder(req.Body).Decode(rlfo) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodGet)
		testHeader(t, req, "application/vnd.github.mercy-preview+json")
		testFormValues(t, req, values{
			"sort": "newest",
			"page": "1",
		})

		fmt.Fprint(wtr, `[{"id":1},{"id":2}]`)
	})

	opt := &github.RepositoryListForksOptions{
		Sort:        "newest",
		ListOptions: github.ListOptions{Page: 1},
	}

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	repos, err := gha.ListForks(ctx, wantOwner, wantRepo, opt)
	if err != nil {
		t.Errorf("githubapi.ListForks returned error: %v", err)
	}

	want := []*github.Repository{{ID: github.Int64(1)}, {ID: github.Int64(2)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("githubapi.ListForks returned %+v, want %+v", repos, want)
	}
}

func TestListForksError(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	wantOwner := "wantOwner"
	wantRepo := "wantRepo"
	wantURL := fmt.Sprintf("/repos/%v/%v/forks", wantOwner, wantRepo)
	srvr.Mux.HandleFunc(wantURL, func(wtr http.ResponseWriter, req *http.Request) {
		rlfo := new(github.RepositoryListForksOptions)
		json.NewDecoder(req.Body).Decode(rlfo) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodGet)
		testHeader(t, req, "application/vnd.github.mercy-preview+json")
		testFormValues(t, req, values{
			"sort": "newest",
			"page": "1",
		})

		wtr.WriteHeader(http.StatusUnprocessableEntity)
	})

	opt := &github.RepositoryListForksOptions{
		Sort:        "newest",
		ListOptions: github.ListOptions{Page: 1},
	}

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	_, err := gha.ListForks(ctx, wantOwner, wantRepo, opt)
	if err == nil {
		t.Errorf("githubapi.ListForks should have returned an error")
	}
}

func TestMergeUpstreamSuccess(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	wantOwner := "Test_MergeUpstream_success_owner"
	wantRepo := "Test_MergeUpstream_success_repo"
	wantBranch := "Test_MergeUpstream_success_branch"

	input := &github.RepoMergeUpstreamRequest{
		Branch: github.String(wantBranch),
	}

	wantURL := fmt.Sprintf("/repos/%s/%s/merge-upstream", wantOwner, wantRepo)

	srvr.Mux.HandleFunc(wantURL, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodPost)
		if !cmp.Equal(rmur, input) {
			t.Errorf("Request body = %+v, want %+v", rmur, input)
		}

		fmt.Fprint(wtr, `{"merge_type":"m"}`)
	})

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	result, err := gha.MergeUpstream(ctx, wantOwner, wantRepo, wantBranch)
	if err != nil {
		t.Errorf("githubapi.MergeUpstream returned error: %v", err)
	}

	want := &github.RepoMergeUpstreamResult{MergeType: github.String("m")}
	if !cmp.Equal(result, want) {
		t.Errorf("githubapi.MergeUpstream returned %+v, want %+v", result, want)
	}
}

func TestMergeUpstreamError(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	wantOwner := "Test_MergeUpstream_error_owner"
	wantRepo := "Test_MergeUpstream_error_repo"
	wantBranch := "Test_MergeUpstream_error_branch"

	input := &github.RepoMergeUpstreamRequest{
		Branch: github.String(wantBranch),
	}

	wantURL := fmt.Sprintf("/repos/%s/%s/merge-upstream", wantOwner, wantRepo)

	srvr.Mux.HandleFunc(wantURL, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodPost)
		if !cmp.Equal(rmur, input) {
			t.Errorf("Request body = %+v, want %+v", rmur, input)
		}

		wtr.WriteHeader(http.StatusUnprocessableEntity)
	})

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	_, err := gha.MergeUpstream(ctx, wantOwner, wantRepo, wantBranch)
	if err == nil {
		t.Errorf("githubapi.MergeUpstream should have returned an error")
	}
}

var TestSyncForksSuccessNoUpdategetRepositoriesHasFired = false //nolint:gochecknoglobals,nolintlint,lll  // This is a test variable

func TestSyncForksSuccessNoUpdate(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	owner := "Test_owner"
	repo := "Test_repo"
	branch := "Test_branch"
	fullRepo := fmt.Sprintf("%s:%s", owner, repo)
	fullBranch := fmt.Sprintf("%s:%s", repo, branch)

	userJSON := `{"login":"` + owner + `","id":666,"name":"My Test User"}`
	reposJSON := `[{"id":123,"owner":` + userJSON + `,"name":"` + repo +
		`","full_name":"` + fullRepo + `","fork":true,"default_branch":"` + branch +
		`"},{"id":321,"owner":` + userJSON + `,"name":"not-a-fork","full_name":"` + owner +
		`:not-a-fork","fork":false,"default_branch":"` + branch + `"}]`

	// setup for client.Users.Get(ctx, userName)
	userURL := "/user"
	srvr.Mux.HandleFunc(userURL, func(wtr http.ResponseWriter, req *http.Request) {
		testMethod(t, req, http.MethodGet)
		fmt.Fprint(wtr, userJSON)
	})

	TestSyncForksSuccessNoUpdategetRepositoriesHasFired = false //nolint:gochecknoglobals,nolintlint,lll  // This is a test variable

	// setup for client.Repositories.List(ctx, user, opts)
	userRepoListURL := fmt.Sprintf("/users/%s/repos", owner)
	srvr.Mux.HandleFunc(userRepoListURL, func(wtr http.ResponseWriter, req *http.Request) {
		rlo := new(github.RepositoryListOptions)
		json.NewDecoder(req.Body).Decode(rlo) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodGet)

		if TestSyncForksSuccessNoUpdategetRepositoriesHasFired {
			fmt.Fprint(wtr, `[]`)
		} else {
			fmt.Fprint(wtr, reposJSON)
			TestSyncForksSuccessNoUpdategetRepositoriesHasFired = true
		}
	})

	// setup for client.Repositories.MergeUpstream(ctx, owner, repo, &req)
	repoMergeURL := fmt.Sprintf("/repos/%s/%s/merge-upstream", owner, repo)
	srvr.Mux.HandleFunc(repoMergeURL, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodPost)

		fmt.Fprint(wtr, `{"message":"This branch is not behind the upstream `+
			fullBranch+`.","merge_type":"none","base_branch":"`+fullRepo+`"}`)
	})

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	err := gha.SyncForks(ctx, "", true, false)
	if err != nil {
		t.Errorf("githubapi.SyncForks returned error: %v", err)
	}
}

func TestSyncForksBadUserName(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	want := fmt.Errorf("api.client.Users.Get error: parse \"users/%%\": invalid URL escape \"%%\"")
	if want == nil {
		t.Errorf("Test_SyncForks_bad_userName want shouldn't be nil")
	}

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	err := gha.SyncForks(ctx, "%", true, false)
	if err == nil {
		t.Errorf("githubapi.SyncForks should have returned an error")
	} else if strings.Compare(err.Error(), want.Error()) != 0 {
		t.Errorf("githubapi.SyncForks returned %+v, want %+v", err, want)
	}
}

func TestSyncForksBadGet(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	login := "Test_user"

	userJSON := `{"login":"` + login + `","id":666,"name":"My Test User"}`

	// setup for client.Users.Get(ctx, userName)
	userURL := "/user"
	srvr.Mux.HandleFunc(userURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, userJSON)
	})

	// setup for client.Repositories.List(ctx, user, opts)
	userRepoListURL := fmt.Sprintf("/users/%s/repos", login)
	srvr.Mux.HandleFunc(userRepoListURL, func(wtr http.ResponseWriter, req *http.Request) {
		rlo := new(github.RepositoryListOptions)
		json.NewDecoder(req.Body).Decode(rlo) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodGet)

		wtr.WriteHeader(http.StatusUnprocessableEntity)
	})

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	err := gha.SyncForks(ctx, "", true, false)
	if err == nil {
		t.Errorf("githubapi.SyncForks should have returned an error")
	}
}

var TestSyncForksBadMergegetRepositoriesHasFired = false //nolint:gochecknoglobals,nolintlint,lll  // This is a test variable

func TestSyncForksBadMerge(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	login := "Test_user"
	owner := "Test_owner"
	repo := "Test_repo"
	branch := "Test_branch"
	fullRepo := fmt.Sprintf("%s:%s", owner, repo)

	userJSON := `{"login":"` + login + `","id":666,"name":"My Test User"}`
	reposJSON := `[{"id":123,"owner":` + userJSON + `,"name":"` + repo + `","full_name":"` + fullRepo +
		`","fork":true,"default_branch":"` + branch + `"}]`

	// setup for client.Users.Get(ctx, userName)
	userURL := "/user"
	srvr.Mux.HandleFunc(userURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, userJSON)
	})

	TestSyncForksBadMergegetRepositoriesHasFired = false

	// setup for client.Repositories.List(ctx, user, opts)
	userRepoListURL := fmt.Sprintf("/users/%s/repos", login)
	srvr.Mux.HandleFunc(userRepoListURL, func(wtr http.ResponseWriter, req *http.Request) {
		rlo := new(github.RepositoryListOptions)
		json.NewDecoder(req.Body).Decode(rlo) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodGet)

		if TestSyncForksBadMergegetRepositoriesHasFired {
			fmt.Fprint(wtr, `[]`)
		} else {
			fmt.Fprint(wtr, reposJSON)
			TestSyncForksBadMergegetRepositoriesHasFired = true
		}
	})

	// setup for client.Repositories.MergeUpstream(ctx, owner, repo, &req)
	repoMergeURL := fmt.Sprintf("/repos/%s/%s/merge-upstream", owner, repo)
	srvr.Mux.HandleFunc(repoMergeURL, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodPost)

		wtr.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	err := gha.SyncForks(ctx, "", true, false)
	if err == nil {
		t.Errorf("githubapi.SyncForks should have returned an error")
	}
}

func TestMergeUpstreamForkSuccessNoUpdate(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	owner := "Test_owner"
	repo := "Test_repo"
	branch := "Test_branch"
	fullRepo := fmt.Sprintf("%s:%s", owner, repo)
	fullBranch := fmt.Sprintf("%s:%s", repo, branch)

	repoMergeURL := fmt.Sprintf("/repos/%s/%s/merge-upstream", owner, repo)
	srvr.Mux.HandleFunc(repoMergeURL, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodPost)

		fmt.Fprint(wtr, `{"message":"This branch is not behind the upstream `+
			fullBranch+`.","merge_type":"none","base_branch":"`+fullRepo+`"}`)
	})

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	err := gha.MergeUpstreamFork(ctx, owner, repo, branch, true)
	if err != nil {
		t.Errorf("githubapi.SyncForks returned error: %v", err)
	}
}

func TestMergeUpstreamForkSuccessFastForward(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	owner := "Test_owner"
	repo := "Test_repo"
	branch := "Test_branch"
	fullRepo := fmt.Sprintf("%s:%s", owner, repo)
	fullBranch := fmt.Sprintf("%s:%s", repo, branch)

	repoMergeURL := fmt.Sprintf("/repos/%s/%s/merge-upstream", owner, repo)
	srvr.Mux.HandleFunc(repoMergeURL, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodPost)

		fmt.Fprint(wtr, `{"message":"Successfully fetched and fast-forwarded from upstream `+
			fullBranch+`.","merge_type":"fast-forward","base_branch":"`+fullRepo+`"}`)
	})

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	err := gha.MergeUpstreamFork(ctx, owner, repo, branch, true)
	if err != nil {
		t.Errorf("githubapi.SyncForks returned error: %v", err)
	}
}

func TestMergeUpstreamForkSuccessMerge(t *testing.T) {
	srvr, serr := httptest.NewHTTPTestServer(githubapi.GitHubAPIBaseURLPath, os.Stderr)
	if serr != nil {
		panic(serr)
	}
	defer srvr.Close()

	owner := "Test_owner"
	repo := "Test_repo"
	branch := "Test_branch"
	fullRepo := fmt.Sprintf("%s:%s", owner, repo)
	fullBranch := fmt.Sprintf("%s:%s", repo, branch)

	repoMergeURL := fmt.Sprintf("/repos/%s/%s/merge-upstream", owner, repo)
	srvr.Mux.HandleFunc(repoMergeURL, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur) //nolint:errcheck  // We don't care about the error here
		testMethod(t, req, http.MethodPost)

		fmt.Fprint(wtr, `{"message":"Successfully merged from upstream `+
			fullBranch+`.","merge_type":"merge","base_branch":"`+fullRepo+`"}`)
	})

	ctx := context.Background()
	gha, nerr := NewTestGitHubAPI(ctx, "auth", srvr.Server.URL)
	if nerr != nil {
		t.Errorf("githubapi.NewGitHubAPI error: %v", nerr)
	}

	err := gha.MergeUpstreamFork(ctx, owner, repo, branch, true)
	if err != nil {
		t.Errorf("githubapi.SyncForks returned error: %v", err)
	}
}
