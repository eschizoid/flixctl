#!/usr/bin/env bash

set -e
set -x
set -o pipefail

set +u
source /home/webhook/.bashrc
set -u

case $# in
   0)
      echo "Usage: $0 {download|search|status}"
      exit 1
      ;;
   1)
      case $1 in
         download)
            /home/webhook/go/bin/flixctl torrent \
                download \
                --magnet-link "${MAGNET_LINK}" \
                --slack-notification-channel "${SLACK_TORRENT_INCOMING_HOOK_URL}"
            ;;
         search)
            /home/webhook/go/bin/flixctl torrent \
                search \
                --keywords "${KEYWORDS}" \
                --minimum-quality "1080" \
                --slack-notification-channel "${SLACK_TORRENT_INCOMING_HOOK_URL}"
            ;;
         status)
            /home/webhook/go/bin/flixctl torrent \
                status \
                --slack-notification-channel \
                "${SLACK_TORRENT_INCOMING_HOOK_URL}"
            ;;
         *)
            echo "'$1' is not a valid torrent command."
            echo "Usage: $0 {download|search|status}"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {download|search|status}"
      exit 3
      ;;
esac
