#!/usr/bin/env bash

set -e
set -x
set -o pipefail

set +u
source /home/webhook/.bashrc
set -u

/home/webhook/go/bin/flixctl plex "${PLEX_COMMAND}"
