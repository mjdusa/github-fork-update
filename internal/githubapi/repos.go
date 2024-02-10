package githubapi

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
)

// ListRepositories list the repositories of the specified user.
func ListRepositories(ctx context.Context, client *github.Client, user string,
	opts *github.RepositoryListOptions) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.List(ctx, user, opts)
	if err != nil {
		return nil, fmt.Errorf("client.Repositories.List error: %w", err)
	}

	return repos, nil
}

// ListForks lists the forks of the specified repository.
func ListForks(ctx context.Context, client *github.Client, owner string, repo string,
	opts *github.RepositoryListForksOptions) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.ListForks(ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("client.Repositories.ListForks error: %w", err)
	}

	return repos, nil
}

// MergeUpstream merges the upstream repository into the fork for the specified branch.
func MergeUpstream(ctx context.Context, client *github.Client, owner string, repo string,
	branch string) (*github.RepoMergeUpstreamResult, error) {
	req := github.RepoMergeUpstreamRequest{
		Branch: &branch,
	}

	result, _, err := client.Repositories.MergeUpstream(ctx, owner, repo, &req)
	if err != nil {
		return nil, fmt.Errorf("client.Repositories.MergeUpstream error: %w", err)
	}

	return result, nil
}

func MergeUpstreamFork(ctx context.Context, client *github.Client, repoOwner string,
	repoName string, repoBranch string, verbose bool) error {
	res, err := MergeUpstream(ctx, client, repoOwner, repoName, repoBranch)
	if err != nil {
		return fmt.Errorf("client.Repositories.MergeUpstreamFork error: %w", err)
	}

	if res.MergeType == nil || *res.MergeType == "none" {
		if verbose {
			fmt.Printf("-> Repo '%s/%s %s' %s\n", repoOwner, repoName, repoBranch, *res.Message)
		}
	} else {
		fmt.Printf("-> Repo '%s/%s %s' %s\n", repoOwner, repoName, repoBranch, *res.Message)
	}

	return nil
}

func SyncForks(ctx context.Context, client *github.Client, userName string, verboseFlag bool, debugFlag bool) error {
	if client == nil {
		return fmt.Errorf("SyncForks error: client is nil")
	}

	user, _, err := client.Users.Get(ctx, userName)
	if err != nil {
		return fmt.Errorf("client.Users.Get error: %w", err)
	}

	page := 1
	perPage := 30

	for {
		//nolint:exhaustruct // defaults are desired except for paging
		opts := github.RepositoryListOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		}

		repos, gerr := ListRepositories(ctx, client, *user.Login, &opts)
		if gerr != nil {
			return fmt.Errorf("ListRepositories error: %w", gerr)
		}

		if len(repos) == 0 {
			break
		}

		for _, repo := range repos {
			if *repo.Fork {
				merr := MergeUpstreamFork(ctx, client, *repo.Owner.Login, *repo.Name, *repo.DefaultBranch, verboseFlag)
				if merr != nil {
					return fmt.Errorf("MergeUpstreamFork error: %w", merr)
				}
			} else if verboseFlag || debugFlag {
				fmt.Printf("-> Repo '%s/%s %s' is not a fork, skipping...\n", *repo.Owner.Login, *repo.Name, *repo.DefaultBranch)
			}
		}

		page++
	}

	return nil
}
