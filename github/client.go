package github

import (
	"context"

	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"fmt"
)

type client struct {
	githubClient *github.Client
	ctx context.Context
}

func New(token string) (*client) {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &client{githubClient: github.NewClient(tc), ctx: ctx}
}

func (c *client) BumpUp() {
	repos, _, err := c.githubClient.Repositories.List(c.ctx, "", nil)
	fmt.Println(repos)
	fmt.Println(err)
}