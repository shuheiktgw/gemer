package main

import (
	"github.com/pkg/errors"
	"github.com/blang/semver"
	"regexp"
	"strings"
)

var versionRegex = regexp.MustCompile(`VERSION\s*=\s*['"](\d+\.\d+\.\d+)['"]`)

// Gemer wraps GithubClient and simplifies interactions with GitHub API
type Gemer struct {
	GitHubClient *GitHubClient
}

// TODO: Enable to specify version from command line
func (g *Gemer) UpdateVersion(branch, path string) error {
	if len(branch) == 0 {
		return errors.New("missing Github branch name")
	}

	if len(path) == 0 {
		return errors.New("missing Github version.rb path")
	}

	content, sha, err := g.GitHubClient.GetVersion(branch, path)

	if err != nil {
		return err
	}

	sc := string(content)

	currentV := extractVersion(sc)

	if len(currentV) == 0 {
		return errors.Errorf("failed to extract version from version.rb: version.rb content: %s", sc)
	}

	nextV, err := convertToNext(currentV)

	if err != nil {
		return err
	}

	newBranchName := "bumps_up_to_" + nextV

	err = g.GitHubClient.CreateNewBranch(branch, newBranchName)

	if err != nil {
		return err
	}

	newContent := strings.Replace(sc, currentV, nextV, 1)

	return g.GitHubClient.UpdateVersion(path, "Bumps up to " + nextV, *sha, newBranchName, []byte(newContent))
}

func extractVersion(c string) string {
	m := versionRegex.FindStringSubmatch(c)

	if m == nil {
		return ""
	}

	return m[1]
}

func convertToNext(current string) (string, error) {
	v, err := semver.New(current)

	if err != nil {
		return "", errors.Wrapf(err, "error occurred while parsing current version: current version: %s", current)
	}

	v.Patch += v.Patch + 1

	return v.String(), nil
}
