package tools_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/tools"
)

func Test_WrapError_nilError(t *testing.T) {
	err := tools.WrapError("Test_WrapError_nilError message", nil)
	if err != nil {
		t.Errorf("tools.WrapError should be return nil")
	}
}

func Test_WrapError_GoodError(t *testing.T) {
	innerErrMsg := "Test_WrapError_GoodError__inner"
	innerErr := fmt.Errorf(innerErrMsg)
	wrapMsg := "Test_WrapError_GoodError message"
	err := tools.WrapError(wrapMsg, innerErr)
	if err == nil {
		t.Errorf("tools.WrapError should NOT return nil when error is not nil")
	}

	want := fmt.Errorf("%s: %w", wrapMsg, innerErr)

	if err.Error() != want.Error() {
		t.Errorf("tools.WrapError returned %+v, want %+v", err, want)
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

	err := tools.SyncForks(ctx, client, "", true)
	if err != nil {
		t.Errorf("tools.SyncForks returned error: %v", err)
	}
}

func Test_SyncForks_nil_client(t *testing.T) {
	want := fmt.Errorf("SyncForks error: client is nil")

	ctx := context.Background()

	err := tools.SyncForks(ctx, nil, "", true)
	if err == nil {
		t.Errorf("tools.SyncForks should have returned an error")
	}

	if strings.Compare(err.Error(), want.Error()) != 0 {
		t.Errorf("tools.SyncForks returned %+v, want %+v", err, want)
	}
}

func Test_SyncForks_bad_userName(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	want := fmt.Errorf("client.Users.Get error:: parse \"users/%%\": invalid URL escape \"%%\"")

	ctx := context.Background()

	err := tools.SyncForks(ctx, client, "%", true)
	if err == nil {
		t.Errorf("tools.SyncForks should have returned an error")
	}

	if strings.Compare(err.Error(), want.Error()) != 0 {
		t.Errorf("tools.SyncForks returned %+v, want %+v", err, want)
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

	err := tools.SyncForks(ctx, client, "", true)
	if err == nil {
		t.Errorf("tools.SyncForks should have returned an error")
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

	err := tools.SyncForks(ctx, client, "", true)
	if err == nil {
		t.Errorf("tools.SyncForks should have returned an error")
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

	err := tools.MergeUpstreamFork(ctx, client, owner, repo, branch, true)
	if err != nil {
		t.Errorf("tools.SyncForks returned error: %v", err)
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

	err := tools.MergeUpstreamFork(ctx, client, owner, repo, branch, true)
	if err != nil {
		t.Errorf("tools.SyncForks returned error: %v", err)
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

	err := tools.MergeUpstreamFork(ctx, client, owner, repo, branch, true)
	if err != nil {
		t.Errorf("tools.SyncForks returned error: %v", err)
	}
}
