#!/usr/bin/env bash

HOOKS="$(cat /opt/webhook-linux-amd64/hooks.json | grep -o 'id.*' | cut -f2- -d:)"

HOOKS="$(echo ${HOOKS::-1})"

echo "{\"ids\":[${HOOKS}]}"
