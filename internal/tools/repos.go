package tools

import (
	"context"

	"github.com/google/go-github/v53/github"
)

// ListRepositories list the repositories of the specified user.
func ListRepositories(ctx context.Context, client *github.Client, user string,
	opts *github.RepositoryListOptions) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.List(ctx, user, opts)
	if err != nil {
		err = WrapError("client.Repositories.List error:", err)
	}

	return repos, err
}

// ListForks lists the forks of the specified repository.
func ListForks(ctx context.Context, client *github.Client, owner string, repo string,
	opts *github.RepositoryListForksOptions) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.ListForks(ctx, owner, repo, opts)
	if err != nil {
		err = WrapError("client.Repositories.ListForks error:", err)
	}

	return repos, err
}

// MergeUpstream merges the upstream repository into the fork for the specified branch.
func MergeUpstream(ctx context.Context, client *github.Client, owner string, repo string,
	branch string) (*github.RepoMergeUpstreamResult, error) {
	req := github.RepoMergeUpstreamRequest{
		Branch: &branch,
	}

	result, _, err := client.Repositories.MergeUpstream(ctx, owner, repo, &req)
	if err != nil {
		err = WrapError("client.Repositories.MergeUpstream", err)
	}

	return result, err
}
