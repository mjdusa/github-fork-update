package githubapi

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
)

func (api *GitHubAPI) ListOrganizations(ctx context.Context, username string,
	opts *github.ListOptions) ([]*github.Organization, error) {
	orgs, _, err := api.Client.Organizations.List(ctx, username, opts)
	if err != nil {
		return nil, fmt.Errorf("client.Organizations.List error: %w", err)
	}

	return orgs, nil
}
