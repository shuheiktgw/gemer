package main

import (
	"fmt"
	"context"

	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"errors"
)

// GitHubClient is a clint to interact with Github API
type GitHubClient struct {
	Repo string
	Client *github.Client
}

// NewGitHubClient creates and initializes a new GitHubClient
func NewGitHubClient(repo, token string) (*GitHubClient, error) {
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
		Repo: repo,
		Client: client,
		}, nil
}

func (c *GitHubClient) BumpUp() {
	repos, _, err := c.githubClient.Repositories.List(c.ctx, "", nil)
	fmt.Println(repos)
	fmt.Println(err)
}