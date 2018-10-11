#!/usr/bin/env bash

set -e
set -u
set -o pipefail

source ~/.bashrc

flixctl torrent download --magnet-link ${MAGNET_LINK}
