#!/bin/bash -eu

go get github.com/Songmu/goxz/cmd/goxz
go get github.com/tcnksm/ghr

# Extract value of Version const from version.go
VERSION=`grep 'Version =' version.go | sed -E 's/.*"(.+)"$$/\1/'`

# Path to built files
FILES=./pkg/dist/v${VERSION}

goxz -pv=v${VERSION} -arch=386,amd64 -d=${FILES}
ghr -t ${GHR_GITHUB_TOKEN} --replace ${VERSION} ${FILES}

