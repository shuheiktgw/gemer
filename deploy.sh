#!/bin/bash -eu

# Extract value of Version const from version.go
VERSION=`grep 'Version =' version.go | sed -E 's/.*"(.+)"$$/\1/'`

# Path to built files
FILES=./pkg/dist/v${VERSION}

goxz -pv=v${VERSION} -arch=386,amd64 -d=${FILES}
ghr -t ${GHR_GITHUB_TOKEN} --replace ${VERSION} ${FILES}

