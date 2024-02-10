package githubapi

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
)

func ListOrganizations(ctx context.Context, client *github.Client, username string,
	opts *github.ListOptions) ([]*github.Organization, error) {
	orgs, _, err := client.Organizations.List(ctx, username, opts)
	if err != nil {
		return nil, fmt.Errorf("client.Organizations.List error: %w", err)
	}

	return orgs, nil
}
