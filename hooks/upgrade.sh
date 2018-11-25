#!/usr/bin/env bash

set -e
set -x
set -o pipefail

rm -rf ~/go/src/github.com/eschizoid/flixctl
go get github.com/eschizoid/flixctl
