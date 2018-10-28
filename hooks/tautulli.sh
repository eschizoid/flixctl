#!/usr/bin/env bash

set -e
set -x
set -o pipefail

set +u
source /home/webhook/.bashrc
set -u

case $# in
   0)
      echo "Usage: $0 {forward}"
      exit 1
      ;;
   1)
      case $1 in
         forward)
            /home/webhook/go/bin/flixctl tautulli forward --message "${TAUTULLI_MESSAGE}" --slack-notification-channel "${SLACK_TORRENT_INCOMING_HOOK_URL}"
            ;;
         *)
            echo "'$1' is not a valid torrent command."
            echo "Usage: $0 {forward}"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {forward}"
      exit 3
      ;;
esac
