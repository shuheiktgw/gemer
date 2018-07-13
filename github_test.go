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

	for _, tc := range cases {
		if _, err := NewGitHubClient(tc.owner, tc.repo, tc.token); err == nil {
			t.Fatal("error is not supposed to be nil")
		}
	}
}

func TestNewGitHubClientSuccess(t *testing.T) {
	if _, err := NewGitHubClient("testOwner", "testRepo", "testToken"); err != nil {
		t.Fatal("unexpected error occured: want: error is nil, got: ", err)
	}
}

func TestGetVersion(t *testing.T) {
	c := testGitHubClient(t)

	_, err := c.GetVersion("master", fmt.Sprintf("lib/%s/version.rb", TestRepo))

	if err != nil {
		t.Fatal("GetVersion failed: ", err)
	}
}