package main

import (
	"context"

	"golang.org/x/oauth2"
	"github.com/pkg/errors"
	"github.com/google/go-github/github"
	"net/http"
	"fmt"
	"strings"
)

// GitHubClient is a clint to interact with Github API
type GitHubClient struct {
	Owner, Repo string
	Client *github.Client
}

// NewGitHubClient creates and initializes a new GitHubClient
func NewGitHubClient(owner, repo, token string) (*GitHubClient, error) {
	if len(owner) == 0 {
		return nil, errors.New("missing Github owner name")
	}

	if len(repo) == 0 {
		return nil, errors.New("missing Github repository name")
	}

	if len(token) == 0 {
		return nil, errors.New("missing Github API token")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
		})
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	return &GitHubClient{
		Owner: owner,
		Repo: repo,
		Client: client,
		}, nil
}

// GetVersion gets the latest version.rb file
func (c *GitHubClient) GetVersion(branch, path string) (*github.RepositoryContent, error) {
	if len(branch) == 0 {
		return nil, errors.New("missing Github branch name")
	}

	if len(path) == 0 {
		return nil, errors.New("missing Github version.rb path")
	}

	if !strings.HasSuffix(path, "version.rb") {
		return nil, errors.Errorf("invalid version file path: version file path must ends with version.rb: invalid path: %s", path)
	}

	file, _, res, err := c.Client.Repositories.GetContents(context.TODO(), c.Owner, c.Repo, path, nil)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get version file")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("get version: invalid status: %s", res.Status)
	}

	return file, nil
}