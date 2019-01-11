#!/usr/bin/env bash
set -e
set -x
set -o pipefail

flixctl library inventory \
    --slack-notification false \
    | jq '.[].ArchiveID' -r \
    | awk 'NF' \
    | while read id; \
        do aws glacier delete-archive \
            --account-id='-' \
            --vault-name='plex' \
            --archive-id=${id};
        done