package main

import (
	"os"
	"io"
)

const (
	ExitCodeOK    = iota
	ExitCodeError
)

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}

type CLI struct {
	outStream, errStream io.Writer
}

func (*CLI)Run(args []string) int {
	return ExitCodeOK
}
