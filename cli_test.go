package main

import (
	"testing"
	"bytes"
	"strings"
	"fmt"
)

func testCli() (*CLI, *bytes.Buffer, *bytes.Buffer) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)

	return &CLI{outStream: outStream, errStream: errStream}, outStream, errStream
}

func TestCliRunFail(t *testing.T) {
	cases := []struct {
		command string
		expectedErrorCode int
	}{
		{command: "gemer -repository testRepo -branch testBranch -path test/path", expectedErrorCode: ExitCodeInvalidFlagError},
		{command: "gemer -username testUser -branch testBranch -path test/path", expectedErrorCode: ExitCodeInvalidFlagError},
	}

	for i, tc := range cases {
		cli, _, _ := testCli()
		args := strings.Split(tc.command, " ")

		if got := cli.Run(args); got != tc.expectedErrorCode {
			t.Fatalf("#%d %q exits with %d, want %d", i, tc.command, got, tc.expectedErrorCode)
		}
	}
}

func TestCliRun_dryRunFlag(t *testing.T) {
	cases := []struct {
		command string
		expectedErrorCode int
	}{
		{command: fmt.Sprintf("gemer -username %s -repository %s -branch %s -d", TestOwner, TestRepo, "master"), expectedErrorCode: ExitCodeOK},
		{command: fmt.Sprintf("gemer -username %s -repository %s -branch %s -dry-run", TestOwner, TestRepo, "master"), expectedErrorCode: ExitCodeOK},
	}

	for i, tc := range cases {
		cli, _, _ := testCli()
		args := strings.Split(tc.command, " ")

		if got := cli.Run(args); got != tc.expectedErrorCode {
			t.Fatalf("#%d %q exits with %d, want %d", i, tc.command, got, tc.expectedErrorCode)
		}
	}
}

func TestCliRun_versionFlag(t *testing.T) {
	command := "gemer -version"

	cli, outStream, _ := testCli()
	args := strings.Split(command, " ")

	if got := cli.Run(args); got != ExitCodeOK {
		t.Fatalf("%q exits with %d, want %d", command, got, ExitCodeOK)
	}

	want := fmt.Sprintf("%s current version v%s\n", Name, Version)
	if got := outStream.String(); got != want {
		t.Fatalf("%q outputs %s, want %s", command, got, want)
	}
}