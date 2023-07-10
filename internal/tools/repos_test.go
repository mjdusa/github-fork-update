package tools_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/tools"
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

	repos, err := tools.ListRepositories(ctx, client, wantUser, opt)
	if err != nil {
		t.Errorf("tools.ListRepositories returned error: %v", err)
	}

	want := []*github.Repository{{ID: github.Int64(1)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("tools.ListRepositories returned %+v, want %+v", repos, want)
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

	_, err := tools.ListRepositories(ctx, client, wantUser, opt)
	if err == nil {
		t.Errorf("tools.ListRepositories should have returned an error")
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

	repos, err := tools.ListForks(ctx, client, wantOwner, wantRepo, opt)
	if err != nil {
		t.Errorf("tools.ListForks returned error: %v", err)
	}

	want := []*github.Repository{{ID: github.Int64(1)}, {ID: github.Int64(2)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("tools.ListForks returned %+v, want %+v", repos, want)
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

	_, err := tools.ListForks(ctx, client, wantOwner, wantRepo, opt)
	if err == nil {
		t.Errorf("tools.ListForks should have returned an error")
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

	result, err := tools.MergeUpstream(ctx, client, wantOwner, wantRepo, wantBranch)
	if err != nil {
		t.Errorf("tools.MergeUpstream returned error: %v", err)
	}

	want := &github.RepoMergeUpstreamResult{MergeType: github.String("m")}
	if !cmp.Equal(result, want) {
		t.Errorf("tools.MergeUpstream returned %+v, want %+v", result, want)
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

	_, err := tools.MergeUpstream(ctx, client, wantOwner, wantRepo, wantBranch)
	if err == nil {
		t.Errorf("tools.MergeUpstream should have returned an error")
	}
}
