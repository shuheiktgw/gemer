package main

import (
	"context"

	"golang.org/x/oauth2"
	"github.com/pkg/errors"
	"github.com/google/go-github/github"
	"net/http"
)

// GitHubClient is a clint to interact with Github API
type GitHubClient struct {
	Owner, Repo string
	Client *github.Client
}

// NewGitHubClient creates and initializes a new GitHubClient
func NewGitHubClient(owner, repo, token string) (*GitHubClient, error) {
	if len(repo) == 0 {
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

// GetLatestRef gets the latest ref of a given branch
func (c *GitHubClient) GetLatestRef(branch string) (*github.Reference, error) {
	if len(branch) == 0 {
		return nil, errors.New("missing Github branch name")
	}

	ref, res, err := c.Client.Git.GetRef(context.TODO(), c.Owner, c.Repo, "heads/" + branch)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get ref")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("get ref: invalid status: %s", res.Status)
	}
	
	return ref, nil
}