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

	currentV := extractVersion(content)

	if len(currentV) == 0 {
		return errors.Errorf("failed to extract version from version.rb: version.rb content: %s", content)
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

	newContent := strings.Replace(content, currentV, nextV, 1)
	message := "Bumps up to " + nextV

	err = g.GitHubClient.UpdateVersion(path, message, sha, newBranchName, []byte(newContent))

	if err != nil {
		return g.rollbackUpdateVersion(err, newBranchName, 0, 0)
	}

	prNum, err := g.GitHubClient.CreatePullRequest(message, newBranchName, branch, message)

	if err != nil {
		return g.rollbackUpdateVersion(err, newBranchName, prNum, 0)
	}

	nextTag := "v" + nextV
	releaseID, err := g.GitHubClient.CreateRelease(nextTag, branch, "Release " + nextTag, nextTag + " is released!")

	if err != nil {
		return g.rollbackUpdateVersion(err, newBranchName, prNum, releaseID)
	}

	return err
}

func (g *Gemer) rollbackUpdateVersion(err error, branchName string, prNum int, releaseID int64) error {
	if len(branchName) != 0 {
		if e := g.GitHubClient.DeleteLatestRef(branchName); e != nil {
			return errors.Wrapf(e, "error occurred while rolling back from UpdateVersion results: original error: %s", err)
		}
	}

	if prNum != 0 {
		if e := g.GitHubClient.ClosePullRequest(prNum); e != nil {
			return errors.Wrapf(e, "error occurred while rolling back from UpdateVersion results: original error: %s", err)
		}
	}

	if releaseID != 0 {
		if e := g.GitHubClient.DeleteRelease(releaseID); e != nil {
			return errors.Wrapf(e, "error occurred while rolling back from UpdateVersion results: original error: %s", err)
		}
	}

	return err
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
