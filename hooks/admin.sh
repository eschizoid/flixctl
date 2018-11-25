#!/usr/bin/env bash

set -e
set -x
set -o pipefail

HOOKS="$(cat /opt/webhook-linux-amd64/hooks.json | grep -o 'id.*' | cut -f2- -d:)"

HOOKS="$(echo ${HOOKS::-1})"

echo "{\"ids\":[${HOOKS}]}"
