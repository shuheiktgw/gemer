package main

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"github.com/pkg/errors"
	"github.com/google/go-github/github"
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

// CreateNewBranch creates a new branch from the heads of the origin
func (c *GitHubClient) CreateNewBranch(origin, new string) error {
	originRef, res, err := c.Client.Git.GetRef(context.TODO(), c.Owner, c.Repo, "heads/" + origin)

	if err != nil {
		return errors.Wrapf(err, "failed to get ref: branch name: %s", origin)
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("get ref: branch name: %s invalid: status: %s", res.Status)
	}

	newRef := &github.Reference{
		Ref: github.String("refs/heads/" + new),
		Object: &github.GitObject{
			SHA: originRef.Object.SHA,
		},
	}

	_, res, err = c.Client.Git.CreateRef(context.TODO(), c.Owner, c.Repo, newRef)

	if err != nil {
		return errors.Wrap(err, "failed to create a new branch")
	}

	if res.StatusCode != http.StatusCreated {
		return errors.Errorf("create ref: invalid status: %s", res.Status)
	}

	return nil
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

	opt := &github.RepositoryContentGetOptions{Ref: branch}

	file, _, res, err := c.Client.Repositories.GetContents(context.TODO(), c.Owner, c.Repo, path, opt)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get version file")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("get version: invalid status: %s", res.Status)
	}

	return file, nil
}

func (c *GitHubClient) UpdateVersion(path, message, sha, branch string, content []byte) error {
	if len(path) == 0 {
		return errors.New("missing Github version.rb path")
	}

	if len(message) == 0 {
		return errors.New("missing Github commit message")
	}

	if len(content) == 0 {
		return errors.New("missing Github content")
	}

	if len(sha) == 0 {
		return errors.New("missing Github file sha")
	}

	if len(branch) == 0 {
		return errors.New("missing Github branch name")
	}

	// ea6e7457c75fc0b2db6dc3b41edb704d57fc6a5d
	opt := &github.RepositoryContentFileOptions{Message: &message, Content: content, SHA: &sha, Branch: &branch}

	_, res, err := c.Client.Repositories.UpdateFile(context.TODO(), c.Owner, c.Repo, path, opt)

	if err != nil {
		return errors.Wrap(err, "failed to update version file")
	}

	if res.StatusCode != http.StatusCreated {
		return errors.Errorf("create version: invalid status: %s", res.Status)
	}

	return nil
}