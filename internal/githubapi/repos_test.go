package githubapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/githubapi"
)

func Test_ListRepositories_success(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantUser := "Test_ListRepositories_success_user"
	wantUrl := fmt.Sprintf("/users/%s/repos", wantUser)
	wantAcceptHeaders := []string{"application/vnd.github.mercy-preview+json", "application/vnd.github.nebula-preview+json"}
	mux.HandleFunc(wantUrl, func(wtr http.ResponseWriter, req *http.Request) {
		v := new(github.RepositoryListOptions)
		json.NewDecoder(req.Body).Decode(v)
		testMethod(t, req, "GET")
		testHeader(t, req, "Accept", strings.Join(wantAcceptHeaders, ", "))
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

	repos, err := githubapi.ListRepositories(ctx, client, wantUser, opt)
	if err != nil {
		t.Errorf("githubapi.ListRepositories returned error: %v", err)
	}

	want := []*github.Repository{{ID: github.Int64(1)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("githubapi.ListRepositories returned %+v, want %+v", repos, want)
	}
}

func Test_ListRepositories_error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantUser := "Test_ListRepositories_error_user"
	wantUrl := fmt.Sprintf("/users/%s/repos", wantUser)
	wantAcceptHeaders := []string{"application/vnd.github.mercy-preview+json", "application/vnd.github.nebula-preview+json"}
	mux.HandleFunc(wantUrl, func(wtr http.ResponseWriter, req *http.Request) {
		v := new(github.RepositoryListOptions)
		json.NewDecoder(req.Body).Decode(v)
		testMethod(t, req, "GET")
		testHeader(t, req, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, req, values{
			"visibility":  "public",
			"affiliation": "owner,collaborator",
			"sort":        "created",
			"direction":   "asc",
			"page":        "2",
		})

		wtr.WriteHeader(422)
	})

	opt := &github.RepositoryListOptions{
		Visibility:  "public",
		Affiliation: "owner,collaborator",
		Sort:        "created",
		Direction:   "asc",
		ListOptions: github.ListOptions{Page: 2},
	}

	ctx := context.Background()

	_, err := githubapi.ListRepositories(ctx, client, wantUser, opt)
	if err == nil {
		t.Errorf("githubapi.ListRepositories should have returned an error")
	}
}

func Test_ListForks_success(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantOwner := "wantOwner"
	wantRepo := "wantRepo"
	wantUrl := fmt.Sprintf("/repos/%s/%s/forks", wantOwner, wantRepo)
	mux.HandleFunc(wantUrl, func(wtr http.ResponseWriter, req *http.Request) {
		v := new(github.RepositoryListForksOptions)
		json.NewDecoder(req.Body).Decode(v)
		testMethod(t, req, "GET")
		testHeader(t, req, "Accept", "application/vnd.github.mercy-preview+json")
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

	repos, err := githubapi.ListForks(ctx, client, wantOwner, wantRepo, opt)
	if err != nil {
		t.Errorf("githubapi.ListForks returned error: %v", err)
	}

	want := []*github.Repository{{ID: github.Int64(1)}, {ID: github.Int64(2)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("githubapi.ListForks returned %+v, want %+v", repos, want)
	}
}

func Test_ListForks_error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantOwner := "wantOwner"
	wantRepo := "wantRepo"
	wantUrl := fmt.Sprintf("/repos/%v/%v/forks", wantOwner, wantRepo)
	mux.HandleFunc(wantUrl, func(wtr http.ResponseWriter, req *http.Request) {
		v := new(github.RepositoryListForksOptions)
		json.NewDecoder(req.Body).Decode(v)
		testMethod(t, req, "GET")
		testHeader(t, req, "Accept", "application/vnd.github.mercy-preview+json")
		testFormValues(t, req, values{
			"sort": "newest",
			"page": "1",
		})

		wtr.WriteHeader(422)
	})

	opt := &github.RepositoryListForksOptions{
		Sort:        "newest",
		ListOptions: github.ListOptions{Page: 1},
	}

	ctx := context.Background()

	_, err := githubapi.ListForks(ctx, client, wantOwner, wantRepo, opt)
	if err == nil {
		t.Errorf("githubapi.ListForks should have returned an error")
	}
}

func Test_MergeUpstream_success(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantOwner := "Test_MergeUpstream_success_owner"
	wantRepo := "Test_MergeUpstream_success_repo"
	wantBranch := "Test_MergeUpstream_success_branch"

	input := &github.RepoMergeUpstreamRequest{
		Branch: github.String(wantBranch),
	}

	wantUrl := fmt.Sprintf("/repos/%s/%s/merge-upstream", wantOwner, wantRepo)

	mux.HandleFunc(wantUrl, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur)
		testMethod(t, req, "POST")
		if !cmp.Equal(rmur, input) {
			t.Errorf("Request body = %+v, want %+v", rmur, input)
		}

		fmt.Fprint(wtr, `{"merge_type":"m"}`)
	})

	ctx := context.Background()

	result, err := githubapi.MergeUpstream(ctx, client, wantOwner, wantRepo, wantBranch)
	if err != nil {
		t.Errorf("githubapi.MergeUpstream returned error: %v", err)
	}

	want := &github.RepoMergeUpstreamResult{MergeType: github.String("m")}
	if !cmp.Equal(result, want) {
		t.Errorf("githubapi.MergeUpstream returned %+v, want %+v", result, want)
	}
}

func Test_MergeUpstream_error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantOwner := "Test_MergeUpstream_error_owner"
	wantRepo := "Test_MergeUpstream_error_repo"
	wantBranch := "Test_MergeUpstream_error_branch"

	input := &github.RepoMergeUpstreamRequest{
		Branch: github.String(wantBranch),
	}

	wantUrl := fmt.Sprintf("/repos/%s/%s/merge-upstream", wantOwner, wantRepo)

	mux.HandleFunc(wantUrl, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur)
		testMethod(t, req, "POST")
		if !cmp.Equal(rmur, input) {
			t.Errorf("Request body = %+v, want %+v", rmur, input)
		}

		wtr.WriteHeader(422)
	})

	ctx := context.Background()

	_, err := githubapi.MergeUpstream(ctx, client, wantOwner, wantRepo, wantBranch)
	if err == nil {
		t.Errorf("githubapi.MergeUpstream should have returned an error")
	}
}

var Test_SyncForks_success_no_update_getRepositories_HasFired = false

func Test_SyncForks_success_no_update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	owner := "Test_owner"
	repo := "Test_repo"
	branch := "Test_branch"
	fullRepo := fmt.Sprintf("%s:%s", owner, repo)
	fullBranch := fmt.Sprintf("%s:%s", repo, branch)

	userJson := `{"login":"` + owner + `","id":666,"name":"My Test User"}`
	reposJson := `[{"id":123,"owner":` + userJson + `,"name":"` + repo + `","full_name":"` + fullRepo + `","fork":true,"default_branch":"` + branch + `"},{"id":321,"owner":` + userJson + `,"name":"not-a-fork","full_name":"` + owner + `:not-a-fork","fork":false,"default_branch":"` + branch + `"}]`

	// setup for client.Users.Get(ctx, userName)
	userUrl := "/user"
	mux.HandleFunc(userUrl, func(wtr http.ResponseWriter, req *http.Request) {
		testMethod(t, req, "GET")
		fmt.Fprint(wtr, userJson)
	})

	Test_SyncForks_success_no_update_getRepositories_HasFired = false

	// setup for client.Repositories.List(ctx, user, opts)
	userRepoListUrl := fmt.Sprintf("/users/%s/repos", owner)
	mux.HandleFunc(userRepoListUrl, func(wtr http.ResponseWriter, req *http.Request) {
		rlo := new(github.RepositoryListOptions)
		json.NewDecoder(req.Body).Decode(rlo)
		testMethod(t, req, "GET")

		if Test_SyncForks_success_no_update_getRepositories_HasFired {
			fmt.Fprint(wtr, `[]`)
		} else {
			fmt.Fprint(wtr, reposJson)
			Test_SyncForks_success_no_update_getRepositories_HasFired = true
		}
	})

	// setup for client.Repositories.MergeUpstream(ctx, owner, repo, &req)
	repoMergeUrl := fmt.Sprintf("/repos/%s/%s/merge-upstream", owner, repo)
	mux.HandleFunc(repoMergeUrl, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur)
		testMethod(t, req, "POST")

		fmt.Fprint(wtr, `{"message":"This branch is not behind the upstream `+fullBranch+`.","merge_type":"none","base_branch":"`+fullRepo+`"}`)
	})

	ctx := context.Background()

	err := githubapi.SyncForks(ctx, client, "", true, false)
	if err != nil {
		t.Errorf("githubapi.SyncForks returned error: %v", err)
	}
}

func Test_SyncForks_nil_client(t *testing.T) {
	want := fmt.Errorf("SyncForks error: client is nil")
	if want == nil {
		t.Errorf("Test_SyncForks_nil_client want shouldn't be nil")
	}

	ctx := context.Background()

	err := githubapi.SyncForks(ctx, nil, "", true, false)
	if err == nil {
		t.Errorf("githubapi.SyncForks should have returned an error")
	} else if strings.Compare(err.Error(), want.Error()) != 0 {
		t.Errorf("githubapi.SyncForks returned %+v, want %+v", err, want)
	}
}

func Test_SyncForks_bad_userName(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	want := fmt.Errorf("client.Users.Get error:: parse \"users/%%\": invalid URL escape \"%%\"")
	if want == nil {
		t.Errorf("Test_SyncForks_bad_userName want shouldn't be nil")
	}

	ctx := context.Background()

	err := githubapi.SyncForks(ctx, client, "%", true, false)
	if err == nil {
		t.Errorf("githubapi.SyncForks should have returned an error")
	} else if strings.Compare(err.Error(), want.Error()) != 0 {
		t.Errorf("githubapi.SyncForks returned %+v, want %+v", err, want)
	}
}

func Test_SyncForks_bad_get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	login := "Test_user"

	userJson := `{"login":"` + login + `","id":666,"name":"My Test User"}`

	// setup for client.Users.Get(ctx, userName)
	userUrl := "/user"
	mux.HandleFunc(userUrl, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, userJson)
	})

	// setup for client.Repositories.List(ctx, user, opts)
	userRepoListUrl := fmt.Sprintf("/users/%s/repos", login)
	mux.HandleFunc(userRepoListUrl, func(wtr http.ResponseWriter, req *http.Request) {
		rlo := new(github.RepositoryListOptions)
		json.NewDecoder(req.Body).Decode(rlo)
		testMethod(t, req, "GET")

		wtr.WriteHeader(422)
	})

	ctx := context.Background()

	err := githubapi.SyncForks(ctx, client, "", true, false)
	if err == nil {
		t.Errorf("githubapi.SyncForks should have returned an error")
	}
}

var Test_SyncForks_bad_merge_getRepositories_HasFired = false

func Test_SyncForks_bad_merge(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	login := "Test_user"
	owner := "Test_owner"
	repo := "Test_repo"
	branch := "Test_branch"
	fullRepo := fmt.Sprintf("%s:%s", owner, repo)

	userJson := `{"login":"` + login + `","id":666,"name":"My Test User"}`
	reposJson := `[{"id":123,"owner":` + userJson + `,"name":"` + repo + `","full_name":"` + fullRepo + `","fork":true,"default_branch":"` + branch + `"}]`

	// setup for client.Users.Get(ctx, userName)
	userUrl := "/user"
	mux.HandleFunc(userUrl, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, userJson)
	})

	Test_SyncForks_bad_merge_getRepositories_HasFired = false

	// setup for client.Repositories.List(ctx, user, opts)
	userRepoListUrl := fmt.Sprintf("/users/%s/repos", login)
	mux.HandleFunc(userRepoListUrl, func(wtr http.ResponseWriter, req *http.Request) {
		rlo := new(github.RepositoryListOptions)
		json.NewDecoder(req.Body).Decode(rlo)
		testMethod(t, req, "GET")

		if Test_SyncForks_bad_merge_getRepositories_HasFired {
			fmt.Fprint(wtr, `[]`)
		} else {
			fmt.Fprint(wtr, reposJson)
			Test_SyncForks_bad_merge_getRepositories_HasFired = true
		}
	})

	// setup for client.Repositories.MergeUpstream(ctx, owner, repo, &req)
	repoMergeUrl := fmt.Sprintf("/repos/%s/%s/merge-upstream", owner, repo)
	mux.HandleFunc(repoMergeUrl, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur)
		testMethod(t, req, "POST")

		wtr.WriteHeader(404)
	})

	ctx := context.Background()

	err := githubapi.SyncForks(ctx, client, "", true, false)
	if err == nil {
		t.Errorf("githubapi.SyncForks should have returned an error")
	}
}

func Test_MergeUpstreamFork_success_no_update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	owner := "Test_owner"
	repo := "Test_repo"
	branch := "Test_branch"
	fullRepo := fmt.Sprintf("%s:%s", owner, repo)
	fullBranch := fmt.Sprintf("%s:%s", repo, branch)

	repoMergeUrl := fmt.Sprintf("/repos/%s/%s/merge-upstream", owner, repo)
	mux.HandleFunc(repoMergeUrl, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur)
		testMethod(t, req, "POST")

		fmt.Fprint(wtr, `{"message":"This branch is not behind the upstream `+fullBranch+`.","merge_type":"none","base_branch":"`+fullRepo+`"}`)
	})

	ctx := context.Background()

	err := githubapi.MergeUpstreamFork(ctx, client, owner, repo, branch, true)
	if err != nil {
		t.Errorf("githubapi.SyncForks returned error: %v", err)
	}
}

func Test_MergeUpstreamFork_success_fast_forward(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	owner := "Test_owner"
	repo := "Test_repo"
	branch := "Test_branch"
	fullRepo := fmt.Sprintf("%s:%s", owner, repo)
	fullBranch := fmt.Sprintf("%s:%s", repo, branch)

	repoMergeUrl := fmt.Sprintf("/repos/%s/%s/merge-upstream", owner, repo)
	mux.HandleFunc(repoMergeUrl, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur)
		testMethod(t, req, "POST")

		fmt.Fprint(wtr, `{"message":"Successfully fetched and fast-forwarded from upstream `+fullBranch+`.","merge_type":"fast-forward","base_branch":"`+fullRepo+`"}`)
	})

	ctx := context.Background()

	err := githubapi.MergeUpstreamFork(ctx, client, owner, repo, branch, true)
	if err != nil {
		t.Errorf("githubapi.SyncForks returned error: %v", err)
	}
}

func Test_MergeUpstreamFork_success_merge(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	owner := "Test_owner"
	repo := "Test_repo"
	branch := "Test_branch"
	fullRepo := fmt.Sprintf("%s:%s", owner, repo)
	fullBranch := fmt.Sprintf("%s:%s", repo, branch)

	repoMergeUrl := fmt.Sprintf("/repos/%s/%s/merge-upstream", owner, repo)
	mux.HandleFunc(repoMergeUrl, func(wtr http.ResponseWriter, req *http.Request) {
		rmur := new(github.RepoMergeUpstreamRequest)
		json.NewDecoder(req.Body).Decode(rmur)
		testMethod(t, req, "POST")

		fmt.Fprint(wtr, `{"message":"Successfully merged from upstream `+fullBranch+`.","merge_type":"merge","base_branch":"`+fullRepo+`"}`)
	})

	ctx := context.Background()

	err := githubapi.MergeUpstreamFork(ctx, client, owner, repo, branch, true)
	if err != nil {
		t.Errorf("githubapi.SyncForks returned error: %v", err)
	}
}
