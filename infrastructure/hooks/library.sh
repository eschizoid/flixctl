#!/usr/bin/env bash

set -e
#set -x
set -o pipefail

set +u
source /home/webhook/.bashrc
set -u

case $# in
   0)
      echo "Usage: $0 {archive|download|initiate|inventory|jobs}"
      exit 1
      ;;
   1)
      case $1 in
         archive)
            /home/webhook/go/bin/flixctl library archive \
                --file "${FILE}"
            ;;
         download)
            /home/webhook/go/bin/flixctl library download \
                --job-id "${JOB_ID}" \
                --target-file "/plex/movies/glacier/movie-$(date +%Y-%m-%d.%H:%M:%S).zip"
            ;;
         initiate)
            /home/webhook/go/bin/flixctl library initiate
            ;;
         inventory)
            /home/webhook/go/bin/flixctl library inventory \
                --enable-sync "${ENABLE_LIBRARY_SYNC}" \
                --job-id "${JOB_ID}" \
                --slack-notification "${SLACK_NOTIFICATION}" \
                --slack-notification-channel "${SLACK_LIBRARY_INCOMING_HOOK_URL}"
            ;;
         jobs)
            /home/webhook/go/bin/flixctl library jobs \
                --filter "${FILTER}" \
                --slack-notification "${SLACK_NOTIFICATION}" \
                --slack-notification-channel "${SLACK_LIBRARY_INCOMING_HOOK_URL}"
            ;;
         *)
            echo "'$1' is not a valid library command."
            echo "Usage: $0 {archive|download|initiate|inventory|jobs}"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {archive|download|initiate|inventory|jobs}"
      exit 3
      ;;
esac
