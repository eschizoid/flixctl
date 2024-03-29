#!/usr/bin/env bash
set -e
set -x
set -o pipefail

flixctl library inventory \
    --slack-notification false \
    | jq '.[].ArchiveID' -r \
    | awk 'NF' \
    | while read id; \
        do flixctl library delete \
            --archive-id ${id};
        done