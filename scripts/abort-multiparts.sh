#!/usr/bin/env bash
set -e
set -x
set -o pipefail

aws glacier list-multipart-uploads \
    --account-id - \
    --vault-name plex \
    | jq '.UploadsList[]' \
    | jq -sc 'sort_by(.CreationDate) | reverse' \
    | jq '.[].MultipartUploadId' -r \
    | while read id; \
        do aws glacier abort-multipart-upload \
            --account-id - \
            --vault-name plex \
            --upload-id "${id}";
        done
