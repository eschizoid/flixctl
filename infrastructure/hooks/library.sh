#!/usr/bin/env bash

set -e
#set -x
set -o pipefail

set +u
source /home/webhook/.bashrc
set -u

case $# in
   0)
      echo "Usage: $0 {catalogue|delete|download|initiate|initiate-archive|inventory|jobs|upload}"
      exit 1
      ;;
   1)
      case $1 in
         catalogue)
            /home/webhook/go/bin/flixctl library catalogue \
                --filter "${FILTER}" \
                --slack-notification "${SLACK_NOTIFICATION}" \
                --slack-notification-channel "${SLACK_REQUESTS_HOOK_URL}"
            ;;
         delete)
            /home/webhook/go/bin/flixctl library delete \
                --archive-id "${ARCHIVE_ID}"
            ;;
         download)
            /home/webhook/go/bin/flixctl library download \
                --job-id "${JOB_ID}" \
                --target-file "/plex/glacier/downloads/movie-$(date +%Y-%m-%d.%H:%M:%S).zip"
            ;;
         initiate)
            /home/webhook/go/bin/flixctl library initiate
            ;;
         initiate-archive)
            /home/webhook/go/bin/flixctl library initiate \
                --archive-id "${ARCHIVE_ID}"
            ;;
         inventory)
            /home/webhook/go/bin/flixctl library inventory \
                --enable-sync "${ENABLE_LIBRARY_SYNC}" \
                --job-id "${JOB_ID}" \
                --slack-notification "${SLACK_NOTIFICATION}" \
                --slack-notification-channel "${SLACK_REQUESTS_HOOK_URL}"
            ;;
         jobs)
            /home/webhook/go/bin/flixctl library jobs \
                --filter "${FILTER}" \
                --slack-notification "${SLACK_NOTIFICATION}" \
                --slack-notification-channel "${SLACK_REQUESTS_HOOK_URL}"
            ;;
         upload)
            /home/webhook/go/bin/flixctl library upload \
                --enable-batch-mode "false" \
                --source-file "${SOURCE_FILE}"
            ;;
         *)
            echo "'$1' is not a valid library command."
            echo "Usage: $0 {catalogue|delete|download|initiate|initiate-archive|inventory|jobs|upload}"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {catalogue|delete|download|initiate|initiate-archive|inventory|jobs|upload}"
      exit 3
      ;;
esac
