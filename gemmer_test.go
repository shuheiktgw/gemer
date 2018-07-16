package main

import (
	"testing"
	"fmt"
)

func TestGemmerUpdateVersionSuccess(t *testing.T) {
	cases := []struct {
		branch, path string
	}{
		{branch: "develop", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
	}

	for i, tc := range cases {
		c := testGitHubClient(t)
		g := Gemer{GitHubClient: c}

		err := g.UpdateVersion(tc.branch, tc.path)

		if err != nil {
			t.Fatalf("#%d error occurred while updating version: %s", i, err)
		}
	}
}