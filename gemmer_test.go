package main

import (
	"testing"
	"fmt"
)

func testGemmer(t *testing.T) *Gemer {
	c := testGitHubClient(t)
	return &Gemer{GitHubClient: c}
}

func TestGemerUpdateVersionSuccess(t *testing.T) {
	cases := []struct {
		branch, path string
	}{
		{branch: "develop", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
	}

	for i, tc := range cases {
		g := testGemmer(t)

		nb, prNum, releaseID, err := g.UpdateVersion(tc.branch, tc.path)

		if len(nb) != 0 {
			if e := g.GitHubClient.DeleteLatestRef(nb); e != nil {
				t.Errorf("%d error occurred while deleting newly created branch: branch name: %s, error: %s", i, nb, e)
			}
		}

		if prNum != 0 {
			if e := g.GitHubClient.ClosePullRequest(prNum); e != nil {
				t.Errorf("%d error occurred while closing newly created pr: PR number: %d, error: %s", i, prNum, e)
			}
		}

		if releaseID != 0 {
			if e := g.GitHubClient.DeleteRelease(releaseID); e != nil {
				t.Errorf("%d error occurred while deleting newly created release: Release IDr: %d, error: %s", i, releaseID, e)
			}
		}

		if err != nil {
			t.Fatalf("#%d error occurred while updating version: %s", i, err)
		}
	}
}