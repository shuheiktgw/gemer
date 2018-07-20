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

		err := g.UpdateVersion(tc.branch, tc.path)

		if err != nil {
			t.Fatalf("#%d error occurred while updating version: %s", i, err)
		}
	}
}