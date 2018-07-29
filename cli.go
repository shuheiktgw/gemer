package main

import (
	"io"
	"flag"
	"os"
	"fmt"
)

const EnvGitHubToken = "GITHUB_TOKEN"

const (
	ExitCodeOK    = iota
	ExitCodeError
	ExitCodeParseFlagsError
)

type CLI struct {
	outStream, errStream io.Writer
}

func (cli *CLI)Run(args []string) int {
	var (
		owner string
		repo string
		branch string
		path string
		token string
		version bool
		patch bool
		minor bool
		major bool
	)

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&owner, "username", "", "a long option for a GitHub username of your gem")
	flags.StringVar(&owner, "u", "", "a short option for a GitHub username of your gem")

	flags.StringVar(&repo, "repository", "", "a long option for a GitHub repository of your gem")
	flags.StringVar(&repo, "r", "", "a short option for a GitHub repository of your gem")

	flags.StringVar(&branch, "branch", "master", "a long option for a GitHub branch your release is based on")
	flags.StringVar(&branch, "b", "master", "a long option for a GitHub branch your release is based on")

	flags.StringVar(&path, "path", "", "a long option for a path to version.rb from the root of your gem")
	flags.StringVar(&path, "p", "", "a short option for a path to version.rb from the root of your gem")

	flags.StringVar(&token, "token", os.Getenv(EnvGitHubToken), "a long option for a GitHub token")
	flags.StringVar(&token, "t", os.Getenv(EnvGitHubToken), "a short option for a GitHub token")

	flags.BoolVar(&version, "version", false, "a long option to show the current version of gemer")
	flags.BoolVar(&version, "v", false, "a short option to show the current version of gemer")

	flags.BoolVar(&major, "major", false, "an option to increment major version")
	flags.BoolVar(&minor, "minor", false, "an option to increment minor version")
	flags.BoolVar(&patch, "patch", false, "an option to increment patch version")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagsError
	}

	if version {
		fmt.Fprintf(cli.outStream, OutputVersion())
		return ExitCodeOK
	}

	if len(owner) == 0 {
		fmt.Fprintf(cli.errStream, "Failed to set up gemer: GitHub username is missing\n" +
			"Please set it via `-u` option\n\n")
	}

	if len(repo) == 0 {
		fmt.Fprintf(cli.errStream, "Failed to set up gemer: GitHub repository nane is missing\n" +
			"Please set it via `-r` option\n\n")
	}

	if len(branch) == 0 {
		fmt.Fprintf(cli.errStream, "Failed to set up gemer: GitHub branch nane is missing\n" +
			"Please set it via `-b` option\n\n")
	}

	if len(token) == 0 {
		fmt.Fprintf(cli.errStream, "Failed to set up gemer: GitHub Personal Access Token is missing\n" +
			"Please set it via `%s` environment variable or `-t` option\n\n" +
			"To create GitHub Personal Access token, see https://bit.ly/2rvbeT1\n",
			EnvGitHubToken)
	}

	// The default is PatchVersion
	ver := PatchVersion

	if major {
		ver = MajorVersion
	}

	if minor {
		ver = MinorVersion
	}

	client, err := NewGitHubClient(owner, repo, token)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Failed to create a GitHub client: %s\n", err)
		return ExitCodeError
	}

	gemer := Gemer{GitHubClient: client}

	_, prNum, releaseID, err := gemer.UpdateVersion(branch, path, ver)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Failed to update version: %s\n", err)
		return ExitCodeError
	}

	prURL := fmt.Sprintf("https://github.com/%s/%s/pull/%d", owner, repo, prNum)
	releaseURL := fmt.Sprintf("https://github.com/%s/%s/releases/%d", owner, repo, releaseID)

	fmt.Fprintf(cli.outStream, "Now, your gem is ready to release! Remaining tasks are ...\n\n" +
		"1. Access %s and merge the PR\n" +
		"2. Access %s and publish the release\n", prURL, releaseURL)

	return ExitCodeOK
}
