#!/usr/bin/env bash

set -e
set -u
set -o pipefail

source ~/.bash_profile

flixctl torrent download --magnet-link ${MAGNET_LINK}
