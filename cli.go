package main

import (
	"io"
)

const EnvGitHubToken = "GITHUB_TOKEN"

const (
	ExitCodeOK    = iota
	ExitCodeError
)

type CLI struct {
	outStream, errStream io.Writer
}

func (*CLI)Run(args []string) int {
	return ExitCodeOK
}