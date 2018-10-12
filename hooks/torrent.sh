#!/usr/bin/env bash

set -e
set -o pipefail

set +u
source ~/.bashrc
set -u

case $# in
   0)
      echo "Usage: $0 {download|search|status}"
      exit 1
      ;;
   1)
      case $1 in
         download)
            /home/webhook/go/bin/flixctl torrent download --magnet-link ${MAGNET_LINK}
            ;;
         search)
            /home/webhook/go/bin/flixctl torrent search --keywords ${KEYWORDS} --minimum-quality "1080"
            ;;
         status)
            /home/webhook/go/bin/flixctl torrent status
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
