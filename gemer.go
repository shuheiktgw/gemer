package main

import (
	"regexp"
	"strings"
	"fmt"
	"io"
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/blang/semver"
	"github.com/google/go-github/github"
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
	PrURL string
	ReleaseURL string
}

func (g *Gemer) UpdateVersion(branch, path string, version int) (*UpdateVersionResult, error) {
	rc, err := g.GitHubClient.GetVersion(branch, path)

	if err != nil {
		return nil, err
	}

	content, err := decodeContent(rc)

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
	err = g.GitHubClient.UpdateVersion(path, message, *rc.SHA, newBranchName, []byte(newContent))
	result := &UpdateVersionResult{Branch: newBranchName}

	if err != nil {
		return result, g.rollbackUpdateVersion(err, result)
	}

	fmt.Fprintln(g.outStream, "==> Create a new pull request")
	pr, err := g.GitHubClient.CreatePullRequest(message, newBranchName, branch, message)
	result = &UpdateVersionResult{Branch: newBranchName, PrNumber: *pr.Number}

	if err != nil {
		return result, g.rollbackUpdateVersion(err, result)
	}

	currentTag := "v" + currentV
	ccs, err := g.GitHubClient.CompareCommits(currentTag, newBranchName)

	if err != nil {
		return result, g.rollbackUpdateVersion(err, result)
	}

	nextTag := "v" + nextV
	releaseBody := nextTag + " will include commits below!\n" + ccs.String()
	fmt.Fprintln(g.outStream, "==> Create a release")
	release, err := g.GitHubClient.CreateRelease(nextTag, branch, "Release " + nextTag, releaseBody)
	result = &UpdateVersionResult{Branch: newBranchName, PrNumber: *pr.Number, ReleaseID: *release.ID, PrURL: *pr.HTMLURL, ReleaseURL: *release.HTMLURL}

	if err != nil {
		return result, g.rollbackUpdateVersion(err, result)
	}

	return result, nil
}

func(g *Gemer) DryUpdateVersion(branch, path string, version int) error {
	rc, err := g.GitHubClient.GetVersion(branch, path)

	if err != nil {
		return err
	}

	content, err := decodeContent(rc)

	if err != nil {
		return err
	}

	currentV := extractVersion(content)

	if len(currentV) == 0 {
		return errors.Errorf("failed to extract version from version.rb: version.rb content: %s", content)
	}

	nextV, err := convertToNext(currentV, version)

	if err != nil {
		return err
	}

	currentTag := "v" + currentV
	ccs, err := g.GitHubClient.CompareCommits(currentTag, branch)

	fmt.Fprintf(g.outStream, "==> Create a branch named `bumps_up_to_%s`\n", nextV)
	fmt.Fprintf(g.outStream, "==> Update `Version` constant of the branch from `%s` to `%s`\n", currentV, nextV)
	fmt.Fprintf(g.outStream, "==> Create a pull request from `bumps_up_to_%s` branch to `%s` branch\n", nextV, branch)
	fmt.Fprintf(g.outStream, "==> Draft a release which contains the following commits\n\n")
	fmt.Fprintln(g.outStream, ccs)

	return nil
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

func decodeContent(rc *github.RepositoryContent) (string, error) {
	if *rc.Encoding != "base64" {
		return "", errors.Errorf("unexpected encoding: %s", *rc.Encoding)
	}

	decoded, err := base64.StdEncoding.DecodeString(*rc.Content)

	if err != nil {
		return "", errors.Wrap(err, "error occurred while decoding version.rb file")
	}

	return string(decoded), nil
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
