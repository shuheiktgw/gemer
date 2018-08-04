package main

import (
	"regexp"
	"strings"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/blang/semver"
)

const (
	MajorVersion = iota
	MinorVersion
	PatchVersion
)

var versionRegex = regexp.MustCompile(`VERSION\s*=\s*['"](\d+\.\d+\.\d+)['"]`)

// Gemer wraps GithubClient and simplifies interactions with GitHub API
type Gemer struct {
	GitHubClient *GitHubClient
	outStream io.Writer
}

type UpdateVersionResult struct {
	Branch string
	PrNumber int
	ReleaseID int64
}

// TODO: Enable to specify version from command line
func (g *Gemer) UpdateVersion(branch, path string, version int) (*UpdateVersionResult, error) {
	if len(branch) == 0 {
		return nil, errors.New("missing Github branch name")
	}

	if len(path) == 0 {
		return nil, errors.New("missing Github version.rb path")
	}

	content, sha, err := g.GitHubClient.GetVersion(branch, path)

	if err != nil {
		return nil, err
	}

	currentV := extractVersion(content)

	if len(currentV) == 0 {
		return nil, errors.Errorf("failed to extract version from version.rb: version.rb content: %s", content)
	}

	nextV, err := convertToNext(currentV, version)

	if err != nil {
		return nil, err
	}

	newBranchName := "bumps_up_to_" + nextV
	fmt.Fprintln(g.outStream, "==> Create a new branch")
	err = g.GitHubClient.CreateNewBranch(branch, newBranchName)

	if err != nil {
		return nil, err
	}

	newContent := strings.Replace(content, currentV, nextV, 1)
	message := "Bumps up to " + nextV

	fmt.Fprintln(g.outStream, "==> Update version.rb")
	err = g.GitHubClient.UpdateVersion(path, message, sha, newBranchName, []byte(newContent))
	result := &UpdateVersionResult{Branch: newBranchName}

	if err != nil {
		return result, g.rollbackUpdateVersion(err, result)
	}

	fmt.Fprintln(g.outStream, "==> Create a new pull request")
	prNum, err := g.GitHubClient.CreatePullRequest(message, newBranchName, branch, message)
	result = &UpdateVersionResult{Branch: newBranchName, PrNumber: prNum}

	if err != nil {
		return result, g.rollbackUpdateVersion(err, result)
	}

	currentTag := "v" + currentV
	ccs, err := g.GitHubClient.CompareCommits(currentTag, newBranchName)

	if err != nil {
		return result, g.rollbackUpdateVersion(err, result)
	}

	nextTag := "v" + nextV
	releaseBody := nextTag + "includes commits below!\n" + ccs.String()
	fmt.Fprintln(g.outStream, "==> Create a release")
	releaseID, err := g.GitHubClient.CreateRelease(nextTag, branch, "Release " + nextTag, releaseBody)
	result = &UpdateVersionResult{Branch: newBranchName, PrNumber: prNum, ReleaseID: releaseID}

	if err != nil {
		return result, g.rollbackUpdateVersion(err, result)
	}

	return result, nil
}

func (g *Gemer) rollbackUpdateVersion(err error, ur *UpdateVersionResult) error {
	if len(ur.Branch) != 0 {
		if e := g.GitHubClient.DeleteLatestRef(ur.Branch); e != nil {
			return errors.Wrapf(e, "error occurred while rolling back from UpdateVersion results: original error: %s", err)
		}
	}

	if ur.PrNumber != 0 {
		if e := g.GitHubClient.ClosePullRequest(ur.PrNumber); e != nil {
			return errors.Wrapf(e, "error occurred while rolling back from UpdateVersion results: original error: %s", err)
		}
	}

	if ur.ReleaseID != 0 {
		if e := g.GitHubClient.DeleteRelease(ur.ReleaseID); e != nil {
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

func convertToNext(current string, version int) (string, error) {
	v, err := semver.New(current)

	if err != nil {
		return "", errors.Wrapf(err, "error occurred while parsing current version: current version: %s", current)
	}

	if version == MajorVersion {
		v.Major = v.Major + 1
	}

	if version == MinorVersion {
		v.Minor = v.Minor + 1
	}

	if version == PatchVersion {
		v.Patch = v.Patch + 1
	}

	return v.String(), nil
}
