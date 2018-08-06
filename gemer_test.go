package main

import (
	"testing"
	"fmt"
	"io/ioutil"
)

func testGemmer(t *testing.T) *Gemer {
	c := testGitHubClient(t)
	return &Gemer{GitHubClient: c, outStream: ioutil.Discard}
}

func TestGemerUpdateVersionSuccess(t *testing.T) {
	cases := []struct {
		branch, path string
	}{
		{branch: "develop", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
	}

	for i, tc := range cases {
		g := testGemmer(t)

		result, err := g.UpdateVersion(tc.branch, tc.path, PatchVersion)

		if err != nil {
			t.Fatalf("#%d error occurred while updating version: %s", i, err)
		} else {
			if e := g.rollbackUpdateVersion(nil, result); e != nil {
				t.Errorf("#%d error occurred while rolling back: %s", i, e)
			}
		}
	}
}

func TestGemerUpdateVersionFail(t *testing.T) {
	cases := []struct {
		branch, path string
	}{
		{branch: "unknown", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
		{branch: "develop", path: "unknown/version.rb"},
	}

	for i, tc := range cases {
		g := testGemmer(t)

		_, err := g.UpdateVersion(tc.branch, tc.path, PatchVersion)

		if err == nil {
			t.Fatalf("#%d error is not supposed to be nil", i)
		}
	}
}

func TestGemerDryUpdateVersionSuccess(t *testing.T) {
	cases := []struct {
		version int
	}{
		{version: PatchVersion},
		{version: MinorVersion},
		{version: MajorVersion},
	}

	for i, tc := range cases {
		g := testGemmer(t)

		err := g.DryUpdateVersion("develop", fmt.Sprintf("lib/%s/version.rb", TestRepo), tc.version)

		if err != nil {
			t.Fatalf("#%d error occurred while dry updating version: %s", i, err)
		}
	}
}

func TestGemerDryUpdateVersionFail(t *testing.T) {
	cases := []struct {
		branch, path string
	}{
		{branch: "unknown", path: fmt.Sprintf("lib/%s/version.rb", TestRepo)},
		{branch: "develop", path: "unknown/version.rb"},
	}

	for i, tc := range cases {
		g := testGemmer(t)

		err := g.DryUpdateVersion(tc.branch, tc.path, PatchVersion)

		if err == nil {
			t.Fatalf("#%d error is not supposed to be nil", i)
		}
	}
}