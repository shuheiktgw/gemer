package main

import (
	"testing"
	"os"
	"fmt"
	"context"
	"net/http"
)

const (
	TestOwner = "shuheiktgw"
	TestRepo = "github-api-test"
)

func testGitHubClient(t *testing.T) *GitHubClient {
	token := os.Getenv(EnvGitHubToken)
	client, err := NewGitHubClient(TestOwner, TestRepo, token)
	if err != nil {
		t.Fatal("NewGitHubClient failed:", err)
	}
	return client
}

func TestNewGitHubClientFail(t *testing.T) {
	cases := []struct {
		owner, repo, token string
	}{
		{owner: "", repo: "testRepo", token: "testToken"},
		{owner: "testOwner", repo: "", token: "testToken"},
		{owner: "testOwner", repo: "testRepo", token: ""},
	}

	for i, tc := range cases {
		if _, err := NewGitHubClient(tc.owner, tc.repo, tc.token); err == nil {
			t.Fatalf("#%d NewGitHubClient: error is not supposed to be nil", i)
		}
	}
}

func TestNewGitHubClientSuccess(t *testing.T) {
	if _, err := NewGitHubClient("testOwner", "testRepo", "testToken"); err != nil {
		t.Fatal("unexpected error occured: want: error is nil, got: ", err)
	}
}

func TestGetVersionFail(t *testing.T) {
	cases := []struct {
		branch, path string
	}{
		{branch: "", path: ""},
		{branch: "", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
		{branch: "develop", path: ""},
		{branch: "unknown", path: "unknown"},
		{branch: "unknown", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
		{branch: "develop", path: "unknown"},
		{branch: "develop", path: "README.md"},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		_, err := c.GetVersion(tc.branch, tc.path)

		if err == nil {
			t.Fatalf("#%d GetVersion: error is not supposed to be nil", i)
		}
	}
}

func TestGetVersionSuccess(t *testing.T) {
	cases := []struct {
		branch, path string
	}{
		{branch: "master", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
		{branch: "develop", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		_, err := c.GetVersion(tc.branch, tc.path)

		if err != nil {
			t.Fatalf("#%d GetVersion failed: %s", i, err)
		}
	}
}

func TestUpdateVersionSuccess(t *testing.T) {
	cases := []struct {
		path, message, sha, branch string
		content []byte
	}{
		{path: fmt.Sprintf("lib/%s/version.rb", TestRepo), message: "test", sha: "ea6e7457c75fc0b2db6dc3b41edb704d57fc6a5d", branch: "new", content: []byte("dGVzdA==")},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		err := c.UpdateVersion(tc.path, tc.message, tc.sha, tc.branch, tc.content)

		if err != nil {
			t.Fatalf("#%d UpdateVersion failed: %s", i, err)
		}
	}
}

func TestCreateNewBranchFail(t *testing.T) {
	cases := []struct {
		origin, new string
	}{
		{origin: "", new: "new"},
		{origin: "unknown", new: "new"},
		{origin: "develop", new: ""},
		{origin: "develop", new: "develop"},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		err := c.CreateNewBranch(tc.origin, tc.new)

		if err == nil {
			deleteBranch(t, c, tc.new)
			t.Fatalf("#%d error is not supposed to be nil: %s", i, err)
		}
	}
}

func TestCreateNewBranchSuccess(t *testing.T) {
	cases := []struct {
		origin, new string
	}{
		{origin: "develop", new: "new"},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		err := c.CreateNewBranch(tc.origin, tc.new)

		if err != nil {
			t.Fatalf("#%d CreateNewBranch failed: %s", i, err)
		}

		deleteBranch(t, c, tc.new)
	}
}

func deleteBranch(t *testing.T, c *GitHubClient, branch string) {
	res, err := c.Client.Git.DeleteRef(context.TODO(), c.Owner, c.Repo, "heads/" + branch)

	if err != nil {
		t.Fatalf("error occurred while deleting a newly created branch: error: %s", err)
	}

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("invalid http status: %s", res.StatusCode)
	}
}