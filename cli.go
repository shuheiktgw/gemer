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

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagsError
	}

	if version {
		fmt.Fprintf(cli.outStream, OutputVersion())
		return ExitCodeOK
	}

	return ExitCodeOK
}
