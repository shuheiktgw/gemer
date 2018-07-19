package main

import (
	"testing"
	"os"
	"fmt"
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
			err = c.DeleteLatestRef(tc.new)

			if err != nil {
				t.Fatalf("#d DeleteLatestRef failed: %s", err)
			}

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

		err = c.DeleteLatestRef(tc.new)

		if err != nil {
			t.Fatalf("#d DeleteLatestRef failed: %s", err)
		}
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

		_, _ , err := c.GetVersion(tc.branch, tc.path)

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

		_, _, err := c.GetVersion(tc.branch, tc.path)

		if err != nil {
			t.Fatalf("#%d GetVersion failed: %s", i, err)
		}
	}
}

// TODO: Move this test to integration tests folder and add error patterns tests
func TestUpdateVersionSuccess(t *testing.T) {
	cases := []struct {
		path, message, sha, branch string
		content []byte
	}{
		{path: fmt.Sprintf("lib/%s/version.rb", TestRepo), message: "Bumps up to 0.1.1", content: []byte("module GithubAPITest\n  VERSION = '0.1.1'\nend")},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		if err := c.CreateNewBranch("develop", "test"); err != nil {
			t.Fatalf("#%d CreateNewBranch failed: %s", i, err)
		}

		_, sha, err := c.GetVersion("test", tc.path)

		if err != nil {
			t.Fatalf("#%d GetVersion failed: %s", i, err)
		}

		if err := c.UpdateVersion(tc.path, tc.message, *sha, "test", tc.content); err != nil {
			t.Fatalf("#%d UpdateVersion failed: %s", i, err)
		}

		err = c.DeleteLatestRef("test")

		if err != nil {
			t.Fatalf("#d DeleteLatestRef failed: %s", err)
		}
	}
}

func TestCreatePullRequestFail(t *testing.T) {
	cases := []struct {
		title, head, base, body string
	}{
		{title: "", head: "pr-test", base: "master", body: "PR!"},
		{title: "test pr", head: "", base: "master", body: "PR!"},
		{title: "test pr", head: "pr-test", base: "", body: "PR!"},
		{title: "test pr", head: "pr-test", base: "master", body: ""},
		{title: "test pr", head: "unknown", base: "master", body: "PR!"},
		{title: "test pr", head: "pr-test", base: "unknown", body: "PR!"},
		{title: "test pr", head: "pr-test", base: "pr-test", body: "PR!"},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		if n, err := c.CreatePullRequest(tc.title, tc.head, tc.base, tc.body); err == nil {
			if e := c.ClosePullRequest(n); e != nil {
				t.Errorf("%d ClosePullRequest failed: might need to close a PR manually: %s", i, e)
			}
			t.Fatalf("#%d CreatePullRequest is supposed to fail", i)
		}
	}
}

// TODO: Move this test to integration tests folder
func TestCreatePullRequestSuccess(t *testing.T) {
	cases := []struct {
		title, head, base, body string
	}{
		{title: "Test PR from develop to master", head: "pr-test", base: "master", body: "This is a test!"},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		n, err := c.CreatePullRequest(tc.title, tc.head, tc.base, tc.body)

		if err != nil {
			t.Fatalf("#%d CreatePullRequest failed: %s", i, err)
		}

		if e := c.ClosePullRequest(n); e != nil {
			t.Errorf("%d ClosePullRequest failed: might need to close a PR manually: %s", i, e)
		}
	}
}

func TestClosePullRequestFail(t *testing.T) {
	c := testGitHubClient(t)

	if err := c.ClosePullRequest(0); err == nil {
		t.Fatalf("ClosePullRequest is supposed to fail")
	}
}