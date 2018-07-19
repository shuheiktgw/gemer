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
// TODO: Might want to return struct instead of simple strings
func (g *Gemer) UpdateVersion(branch, path string) (string, int, error) {
	if len(branch) == 0 {
		return "", 0, errors.New("missing Github branch name")
	}

	if len(path) == 0 {
		return "", 0, errors.New("missing Github version.rb path")
	}

	content, sha, err := g.GitHubClient.GetVersion(branch, path)

	if err != nil {
		return "", 0, err
	}

	currentV := extractVersion(content)

	if len(currentV) == 0 {
		return "", 0, errors.Errorf("failed to extract version from version.rb: version.rb content: %s", content)
	}

	nextV, err := convertToNext(currentV)

	if err != nil {
		return "", 0, err
	}

	newBranchName := "bumps_up_to_" + nextV

	err = g.GitHubClient.CreateNewBranch(branch, newBranchName)

	if err != nil {
		return "", 0, err
	}

	newContent := strings.Replace(content, currentV, nextV, 1)
	message := "Bumps up to " + nextV

	err = g.GitHubClient.UpdateVersion(path, message, sha, newBranchName, []byte(newContent))

	if err != nil {
		return newBranchName, 0, err
	}

	prNum, err := g.GitHubClient.CreatePullRequest(message, newBranchName, branch, message)

	return newBranchName, prNum, err
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
