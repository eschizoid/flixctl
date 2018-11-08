#!/usr/bin/env bash
set -e
set -x
set -o pipefail

aws glacier list-jobs \
    --account-id - \
    --vault-name plex \
    | jq '.JobList[] | select(.ArchiveId != null) | .ArchiveId' -r \
    | while read id; \
        do aws glacier delete-archive \
            --account-id - \
            --vault-name plex \
            --archive-id "${id}";
        done