package main

import (
	"testing"
	"os"
	"fmt"
)

const (
	TestOwner = "shuheiktgwtest"
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

		_ , err := c.GetVersion(tc.branch, tc.path)

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

		content, err := c.GetVersion("test", tc.path)

		if err != nil {
			t.Fatalf("#%d GetVersion failed: %s", i, err)
		}

		if err := c.UpdateVersion(tc.path, tc.message, *content.SHA, "test", tc.content); err != nil {
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

		if pr, err := c.CreatePullRequest(tc.title, tc.head, tc.base, tc.body); err == nil {
			if e := c.ClosePullRequest(*pr.Number); e != nil {
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

		pr, err := c.CreatePullRequest(tc.title, tc.head, tc.base, tc.body)

		if err != nil {
			t.Fatalf("#%d CreatePullRequest failed: %s", i, err)
		}

		if e := c.ClosePullRequest(*pr.Number); e != nil {
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

func TestCreateReleaseFail(t *testing.T) {
	cases := []struct {
		tagName, targetCommitish, name, body string
	}{
		{tagName: "", targetCommitish: "pr-test", name: "Release v0.0.1", body: "v0.0.1 released!"},
		{tagName: "v0.0.1", targetCommitish: "", name: "Release v0.0.1", body: "v0.0.1 released!"},
		{tagName: "v0.0.1", targetCommitish: "pr-test", name: "", body: "v0.0.1 released!"},
		{tagName: "v0.0.1", targetCommitish: "pr-test", name: "Release v0.0.1", body: ""},
		{tagName: "v0.0.1", targetCommitish: "unknown", name: "Release v0.0.1", body: "v0.0.1 released!"},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		if _, err := c.CreateRelease(tc.tagName, tc.targetCommitish, tc.name, tc.body); err == nil {
			t.Fatalf("#%d CreateRelease is supposed to fail", i)
		}
	}
}

func TestCreateReleaseSuccess(t *testing.T) {
	cases := []struct {
		tagName, targetCommitish, name, body string
	}{
		{tagName: "v0.0.1", targetCommitish: "pr-test", name: "Release v0.0.1", body: "v0.0.1 released!"},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		rr, err := c.CreateRelease(tc.tagName, tc.targetCommitish, tc.name, tc.body)

		if err != nil {
			t.Fatalf("#%d CreateRelease failed: %s", i, err)
		}

		if e := c.DeleteRelease(*rr.ID); e != nil {
			t.Errorf("%d DeleteRelease failed: might need to delete a Release manually: %s", i, e)
		}
	}
}

func TestCompareCommitsFail(t *testing.T) {
	cases := []struct {
		base, head string
	}{
		{head: "", base: "master"},
		{head: "master", base: ""},
		{head: "unknown", base: "master"},
		{head: "master", base: "unknown"},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		_, err := c.CompareCommits(tc.head, tc.base)

		if err == nil {
			t.Fatalf("#%d CompareCommits is supposed to fail", i)
		}
	}
}

func TestCompareCommitsSuccess(t *testing.T) {
	cases := []struct {
		base, head string
	}{
		{head: "develop", base: "master"},
		{head: "master", base: "develop"},
		{head: "pr-test", base: "master"},
		{head: "master", base: "pr-test"},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)

		_, err := c.CompareCommits(tc.head, tc.base)

		if err != nil {
			t.Fatalf("#%d CompareCommits failed: %s", i, err)
		}
	}
}

func TestComparedCommitString(t *testing.T) {
	cc := &ComparedCommit{Author: "shuheiktgw", Message: "The Best Commit Ever!", HTMLURL: "https://github.com/shuheiktgw/github-api-test/commit/d6ed804c9bbaefef1832702db562a3b1e98e1291"}
	want := "@shuheiktgw [The Best Commit Ever!](https://github.com/shuheiktgw/github-api-test/commit/d6ed804c9bbaefef1832702db562a3b1e98e1291)"

	if cc.String() != want {
		t.Fatalf("invalid string: want: %s got: %s", want, cc.String())
	}
}

func TestComparedCommitsString(t *testing.T) {
	cc1 := &ComparedCommit{Author: "shuheiktgw", Message: "The Best Commit Ever!", HTMLURL: "https://github.com/shuheiktgw/github-api-test/commit/d6ed804c9bbaefef1832702db562a3b1e98e1291"}
	cc2 := &ComparedCommit{Author: "shuheiktgw", Message: "The Second Best Commit Ever!", HTMLURL: "https://github.com/shuheiktgw/github-api-test/commit/d6ed804c9bbaefef1832702db562a3b1e98e1291"}
	cc3 := &ComparedCommit{Author: "shuheiktgw", Message: "The Third Best Commit Ever!", HTMLURL: "https://github.com/shuheiktgw/github-api-test/commit/d6ed804c9bbaefef1832702db562a3b1e98e1291"}

	ccs := ComparedCommits{Commits: []*ComparedCommit{cc1, cc2, cc3}}

	want := `@shuheiktgw [The Best Commit Ever!](https://github.com/shuheiktgw/github-api-test/commit/d6ed804c9bbaefef1832702db562a3b1e98e1291)
@shuheiktgw [The Second Best Commit Ever!](https://github.com/shuheiktgw/github-api-test/commit/d6ed804c9bbaefef1832702db562a3b1e98e1291)
@shuheiktgw [The Third Best Commit Ever!](https://github.com/shuheiktgw/github-api-test/commit/d6ed804c9bbaefef1832702db562a3b1e98e1291)`

	if ccs.String() != want {
		t.Fatalf("invalid string: want: %s got: %s", want, ccs.String())
	}
}