package githubapi

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
)

const (
	GitHubAPIBaseURLPath = "/api-v3"
	// GitHubAPIVersion     = "v53.1.0"
	// DefaultAPIVersion    = "2022-11-28"
	// GitHubAPIURL       = "https://api.github.com/"
	// GitHubAPIUserAgent = "go-github" + "/" + GitHubAPIVersion
	// GitHubAPIUploadURL = "https://uploads.github.com/"
)

type GitHubAPI struct {
	Client *github.Client
}

func NewGitHubAPI(ctx context.Context, auth string) (*GitHubAPI, error) {
	if len(auth) == 0 {
		return nil, fmt.Errorf("empty token error")
	}

	client := github.NewTokenClient(ctx, auth)

	if client == nil {
		return nil, fmt.Errorf("NewTokenClient returned nil")
	}

	api := GitHubAPI{
		Client: client,
	}

	return &api, nil
}

func (api *GitHubAPI) ListOrganizations(ctx context.Context, username string,
	opts *github.ListOptions) ([]*github.Organization, error) {
	orgs, _, err := api.Client.Organizations.List(ctx, username, opts)
	if err != nil {
		return nil, fmt.Errorf("client.Organizations.List error: %w", err)
	}

	return orgs, nil
}

// ListRepositories list the repositories of the specified user.
func (api *GitHubAPI) ListRepositories(ctx context.Context, user string,
	opts *github.RepositoryListOptions) ([]*github.Repository, error) {
	repos, _, err := api.Client.Repositories.List(ctx, user, opts)
	if err != nil {
		return nil, fmt.Errorf("api.client.Repositories.List error: %w", err)
	}

	return repos, nil
}

// ListForks lists the forks of the specified repository.
func (api *GitHubAPI) ListForks(ctx context.Context, owner string, repo string,
	opts *github.RepositoryListForksOptions) ([]*github.Repository, error) {
	repos, _, err := api.Client.Repositories.ListForks(ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("api.client.Repositories.ListForks error: %w", err)
	}

	return repos, nil
}

// MergeUpstream merges the upstream repository into the fork for the specified branch.
func (api *GitHubAPI) MergeUpstream(ctx context.Context, owner string, repo string,
	branch string) (*github.RepoMergeUpstreamResult, error) {
	req := github.RepoMergeUpstreamRequest{
		Branch: &branch,
	}

	result, _, err := api.Client.Repositories.MergeUpstream(ctx, owner, repo, &req)
	if err != nil {
		return nil, fmt.Errorf("api.client.Repositories.MergeUpstream error: %w", err)
	}

	return result, nil
}

func (api *GitHubAPI) MergeUpstreamFork(ctx context.Context, repoOwner string,
	repoName string, repoBranch string, verbose bool) error {
	res, err := api.MergeUpstream(ctx, repoOwner, repoName, repoBranch)
	if err != nil {
		return fmt.Errorf("api.client.Repositories.MergeUpstreamFork error: %w", err)
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

func (api *GitHubAPI) SyncForks(ctx context.Context, userName string, verboseFlag bool, debugFlag bool) error {
	user, _, err := api.Client.Users.Get(ctx, userName)
	if err != nil {
		return fmt.Errorf("api.client.Users.Get error: %w", err)
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

		repos, gerr := api.ListRepositories(ctx, *user.Login, &opts)
		if gerr != nil {
			return fmt.Errorf("ListRepositories error: %w", gerr)
		}

		if len(repos) == 0 {
			break
		}

		for _, repo := range repos {
			if *repo.Fork {
				merr := api.MergeUpstreamFork(ctx, *repo.Owner.Login, *repo.Name, *repo.DefaultBranch, verboseFlag)
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
