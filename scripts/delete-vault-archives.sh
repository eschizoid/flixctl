#!/usr/bin/env bash
set -e
set -x
set -o pipefail

flixctl library retrieve \
    --type "InventoryRetrieval" \
    | jq '.[].ArchiveID' -r \
    | while read id; \
        do aws glacier delete-archive \
            --account-id='-' \
            --vault-name='plex' \
            --archive-id=${id};
        done