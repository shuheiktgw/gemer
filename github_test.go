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