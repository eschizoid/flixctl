#!/usr/bin/env bash

set -e
#set -x
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
                --download-dir "${DOWNLOAD_DIR}" \
                --magnet-link "$(echo  "${MAGNET_LINK}" | base64 --decode)" \
                --slack-notification "${SLACK_NOTIFICATION}" \
                --slack-notification-channel "${SLACK_TORRENT_INCOMING_HOOK_URL}"
            ;;
         search)
            /home/webhook/go/bin/flixctl torrent \
                search \
                --download-dir "${DOWNLOAD_DIR}" \
                --keywords "${KEYWORDS}" \
                --minimum-quality "1080" \
                --slack-notification "${SLACK_NOTIFICATION}" \
                --slack-notification-channel "${SLACK_TORRENT_INCOMING_HOOK_URL}"
            ;;
         status)
            /home/webhook/go/bin/flixctl torrent \
                status \
                --slack-notification-channel \
                --slack-notification "${SLACK_NOTIFICATION}" \
                --slack-notification-channel "${SLACK_TORRENT_INCOMING_HOOK_URL}"
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
