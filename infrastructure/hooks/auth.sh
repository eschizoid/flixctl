#!/usr/bin/env bash

set -e
set -o pipefail

set +u
source /home/webhook/.bashrc
set -u

case $# in
   0)
      echo "Usage: $0 {token}"
      exit 1
      ;;
   1)
      case $1 in
         token)
            /home/webhook/go/bin/flixctl auth token \
                --slack-client-id="${SLACK_CLIENT_ID}" \
                --slack-client-secret="${SLACK_CLIENT_SECRET}" \
                --slack-code="${SLACK_CODE}" \
                --slack-redrect-uri="${SLACK_REDIRECT_URI}"
            ;;
         *)
            echo "'$1' is not a valid admin command."
            echo "Usage: $0 {token}"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {token}"
      exit 3
      ;;
esac
