package tools_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v53/github"
	"github.com/mjdusa/github-fork-update/internal/tools"
)

func Test_ListOrganizations_success(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	wantUser := "Test_ListOrganizations_success_user"
	wantUrl := fmt.Sprintf("/users/%s/orgs", wantUser)
	mux.HandleFunc(wantUrl, func(wtr http.ResponseWriter, req *http.Request) {
		lo := new(github.ListOptions)
		json.NewDecoder(req.Body).Decode(lo)
		testMethod(t, req, "GET")
		testFormValues(t, req, values{
			"page": "2",
		})

		fmt.Fprint(wtr, `[{"id":1}]`)
	})

	ctx := context.Background()
	opts := &github.ListOptions{
		Page: 2,
	}

	repos, err := tools.ListOrganizations(ctx, client, wantUser, opts)
	if err != nil {
		t.Errorf("tools.ListOrganizations returned error: %v", err)
	}

	want := []*github.Organization{{ID: github.Int64(1)}}
	if !cmp.Equal(repos, want) {
		t.Errorf("tools.ListOrganizations returned %+v, want %+v", repos, want)
	}
}

func Test_ListOrganizations_error(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	wantUser := "%"
	opts := &github.ListOptions{
		Page: 2,
	}

	ctx := context.Background()
	_, err := tools.ListOrganizations(ctx, client, wantUser, opts)
	if err == nil {
		t.Errorf("tools.ListOrganizations should have returned error")
	}
}
