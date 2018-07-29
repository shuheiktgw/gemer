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

		branchName, prNum, releaseID, err := g.UpdateVersion(tc.branch, tc.path, PatchVersion)

		if err != nil {
			t.Fatalf("#%d error occurred while updating version: %s", i, err)
		} else {
			if e := g.rollbackUpdateVersion(nil, branchName, prNum, releaseID); e != nil {
				t.Errorf("#%d error occurred while rolling back: %s", i, e)
			}
		}
	}
}

func TestGemerUpdateVersionFail(t *testing.T) {
	cases := []struct {
		branch, path string
	}{
		{branch: "", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
		{branch: "develop", path: ""},
		{branch: "unknown", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
		{branch: "develop", path: "unknown/version.rb"},
	}

	for i, tc := range cases {
		g := testGemmer(t)

		_, _, _, err := g.UpdateVersion(tc.branch, tc.path, PatchVersion)

		if err == nil {
			t.Fatalf("#%d error is not supposed to be nil", i)
		}
	}
}