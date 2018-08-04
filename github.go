package main

import (
	"context"
		"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"fmt"
)

// GitHubClient is a clint to interact with Github API
type GitHubClient struct {
	Owner, Repo string
	Client *github.Client
}

// ComparedCommit represents one commit and mainly used for formatting purpose
type ComparedCommit struct {
	Author, Message, HTMLURL string
}

// ComparedCommits represents a series of commits
type ComparedCommits struct {
	Commits []*ComparedCommit
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

// UpdateVersion updates a version.rb file with a given content
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

	opt := &github.RepositoryContentFileOptions{Message: &message, Content: content, SHA: &sha, Branch: &branch}

	_, res, err := c.Client.Repositories.UpdateFile(context.TODO(), c.Owner, c.Repo, path, opt)

	if err != nil {
		return errors.Wrap(err, "failed to update version file")
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("create version: invalid status: %s", res.Status)
	}

	return nil
}

// TODO Enable to add custom labels to PR
// CreatePullRequest creates a new pull request
func (c *GitHubClient) CreatePullRequest(title, head, base, body string) (*github.PullRequest, error) {
	if len(title) == 0 {
		return nil, errors.New("missing Github Pull Request title")
	}

	if len(head) == 0 {
		return nil, errors.New("missing Github Pull Request head branch")
	}

	if len(base) == 0 {
		return nil, errors.New("missing Github Pull Request base branch")
	}

	if len(body) == 0 {
		return nil, errors.New("missing Github Pull Request body")
	}

	opt := &github.NewPullRequest{Title: &title, Head: &head, Base: &base, Body: &body}

	pr, res, err := c.Client.PullRequests.Create(context.TODO(), c.Owner, c.Repo, opt)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new pull request")
	}

	if res.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("create pull request: invalid status: %s", res.Status)
	}

	return pr, nil
}

// ClosePullRequest closes a Pull Request with a give Pull Request number
func (c *GitHubClient) ClosePullRequest(number int) error {
	opt := &github.PullRequest{State: github.String("close")}

	_, res, err := c.Client.PullRequests.Edit(context.TODO(), c.Owner, c.Repo, number, opt)

	if err != nil {
		return errors.Wrap(err, "failed to close a pull request")
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("create pull request: invalid status: %s", res.Status)
	}

	return nil
}

// CreateRelease creates a new release
func (c *GitHubClient) CreateRelease(tagName, targetCommitish, name, body string) (*github.RepositoryRelease, error) {
	if len(tagName) == 0 {
		return nil, errors.New("missing Github Release Tag Name")
	}

	if len(targetCommitish) == 0 {
		return nil, errors.New("missing Github Release Target Commitish")
	}

	if len(name) == 0 {
		return nil, errors.New("missing Github Release Name")
	}

	if len(body) == 0 {
		return nil, errors.New("missing Github Release body")
	}

	opt := &github.RepositoryRelease{
		TagName: &tagName,
		TargetCommitish: &targetCommitish,
		Name: &name,
		Body: &body,
		Draft: github.Bool(true),
		}

	rr, res, err := c.Client.Repositories.CreateRelease(context.TODO(), c.Owner, c.Repo, opt)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new release")
	}

	if res.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("create release: invalid status: %s", res.Status)
	}

	return rr, nil
}

// DeleteRelease deletes a release
func (c *GitHubClient) DeleteRelease(id int64) (error) {
	res, err := c.Client.Repositories.DeleteRelease(context.TODO(), c.Owner, c.Repo, id)

	if err != nil {
		return errors.Wrap(err, "failed to create a new release")
	}

	if res.StatusCode != http.StatusNoContent {
		return errors.Errorf("create release: invalid status: %s", res.Status)
	}

	return nil
}

// CompareCommits compares and gets diffs between two commits
func (c *GitHubClient) CompareCommits(base, head string) (*ComparedCommits, error) {
	if len(base) == 0 {
		return nil, errors.New("missing GitHub base commit")
	}

	if len(head) == 0 {
		return nil, errors.New("missing GitHub head commit")
	}

	cc, res, err := c.Client.Repositories.CompareCommits(context.TODO(), c.Owner, c.Repo, base, head)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("compare commits: invalid status: %s", res.Status)
	}

	var ccs []*ComparedCommit

	for _, c := range cc.Commits {
		ccs = append(ccs, &ComparedCommit{Author: *c.Author.Login, Message: *c.Commit.Message, HTMLURL: *c.HTMLURL})
	}

	return &ComparedCommits{Commits: ccs}, nil
}

// DeleteLatestRef deletes the latest Ref of the given branch, intended to be used for rollbacks
func (c *GitHubClient) DeleteLatestRef(branch string) error {
	if len(branch) == 0 {
		return errors.New("missing Github branch name")
	}

	res, err := c.Client.Git.DeleteRef(context.TODO(), c.Owner, c.Repo, "heads/" + branch)

	if err != nil {
		return errors.Wrapf(err, "failed to delete the latest ref of a branch name %s: %s", branch, err)
	}

	if res.StatusCode != http.StatusNoContent {
		return errors.Errorf("delete latest ref: invalid status: %s", res.Status)
	}

	return nil
}

func (cc *ComparedCommit) String() string {
	return fmt.Sprintf("@%s [%s](%s)", cc.Author, cc.Message, cc.HTMLURL)
}

func (ccs *ComparedCommits) String() string {
	var commits []byte
	last := len(ccs.Commits) - 1

	for i, c := range ccs.Commits {
		if i == last {
			commits = append(commits, c.String()...)
			continue
		}

		commits = append(commits, c.String() + "\n"...)
	}

	return string(commits)
}