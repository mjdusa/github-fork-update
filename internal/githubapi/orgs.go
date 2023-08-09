package githubapi

import (
	"context"

	"github.com/google/go-github/v53/github"
)

func ListOrganizations(ctx context.Context, client *github.Client, username string,
	opts *github.ListOptions) ([]*github.Organization, error) {
	orgs, _, err := client.Organizations.List(ctx, username, opts)
	if err != nil {
		err = WrapError("client.Organizations.List error:", err)
	}

	return orgs, err
}
