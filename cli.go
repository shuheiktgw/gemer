package main

import (
	"io"
	"flag"
	"os"
	"fmt"
	"strings"
)

const EnvGitHubToken = "GITHUB_TOKEN"

const (
	ExitCodeOK    = iota
	ExitCodeError
	ExitCodeParseFlagsError
	ExitCodeInvalidFlagError
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
		dryRun bool
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

	flags.BoolVar(&dryRun, "dry-run", false, "a long option for dry run")
	flags.BoolVar(&dryRun, "d", false, "a short option for dry run")

	flags.BoolVar(&major, "major", false, "an option to increment major version")
	flags.BoolVar(&minor, "minor", false, "an option to increment minor version")
	flags.BoolVar(&patch, "patch", true, "an option to increment patch version")

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
		return ExitCodeInvalidFlagError
	}

	if len(repo) == 0 {
		fmt.Fprintf(cli.errStream, "Failed to set up gemer: GitHub repository nane is missing\n" +
			"Please set it via `-r` option\n\n")
		return ExitCodeInvalidFlagError
	}

	if len(path) == 0 {
		path = fmt.Sprintf("lib/%s/version.rb", strings.ToLower(repo))
	}

	if len(token) == 0 {
		fmt.Fprintf(cli.errStream, "Failed to set up gemer: GitHub Personal Access Token is missing\n" +
			"Please set it via `%s` environment variable or `-t` option\n\n" +
			"To create GitHub Personal Access token, see https://bit.ly/2rvbeT1\n",
			EnvGitHubToken)
		return ExitCodeInvalidFlagError
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

	gemer := Gemer{GitHubClient: client, outStream: cli.outStream}

	if dryRun {
		err := gemer.DryUpdateVersion(branch, path, ver)
		if err != nil {
			fmt.Fprintf(cli.errStream, "Failed to update version with dry-run option: %s\n", err)
			return ExitCodeError
		}

		return ExitCodeOK
	}

	result, err := gemer.UpdateVersion(branch, path, ver)
	if err != nil {
		fmt.Fprintf(cli.errStream, "Failed to update version: %s\n", err)
		return ExitCodeError
	}

	fmt.Fprintf(cli.outStream, "Now, your gem is ready to release! Remaining tasks are ...\n\n" +
		"1. Access %s and merge the PR\n" +
		"2. Access %s and publish the release\n", result.PrURL, result.ReleaseURL)

	return ExitCodeOK
}
