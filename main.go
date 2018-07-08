package main

import "os"

const (
	ExitCodeOK int = iota
	ExitCodeError
)

func main() {
	os.Exit(Run(os.Args))
}

func Run(args []string) int {
	return ExitCodeOK
}
