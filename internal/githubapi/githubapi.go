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
