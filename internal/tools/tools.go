package tools

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
)

func WrapError(message string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}

	return nil
}

func SyncForks(ctx context.Context, client *github.Client, userName string, verbose bool) error {
	if client == nil {
		return fmt.Errorf("SyncForks error: client is nil")
	}

	user, _, err := client.Users.Get(ctx, userName)
	if err != nil {
		return WrapError("client.Users.Get error:", err)
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
			return WrapError("ListRepositories error:", gerr)
		}

		if len(repos) == 0 {
			break
		}

		for _, repo := range repos {
			if *repo.Fork {
				err = MergeUpstreamFork(ctx, client, *repo.Owner.Login, *repo.Name, *repo.DefaultBranch, verbose)
				if err != nil {
					return err
				}
			} else if verbose {
				fmt.Printf("-> Repo '%s/%s %s' is not a fork, skipping...\n", *repo.Owner.Login, *repo.Name, *repo.DefaultBranch)
			}
		}

		page++
	}

	return nil
}

func MergeUpstreamFork(ctx context.Context, client *github.Client, repoOwner string,
	repoName string, repoBranch string, verbose bool) error {
	res, err := MergeUpstream(ctx, client, repoOwner, repoName, repoBranch)
	if err != nil {
		return WrapError("MergeUpstream error:", err)
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
