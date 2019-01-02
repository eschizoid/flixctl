#!/usr/bin/env bash

set -e
set -x
set -o pipefail

set +u
source /home/webhook/.bashrc
set -u

case $# in
   0)
      echo "Usage: $0 {jobs|archive|download|initiate|inventory|jobs}"
      exit 1
      ;;
   1)
      case $1 in
         archive)
            /home/webhook/go/bin/flixctl library archive \
                --file "${FILE}"
            ;;
         donwload)
            /home/webhook/go/bin/flixctl library download \
                --job-id "${JOB_ID}" \
                --file "/plex/movies/movie-$(date +%Y-%m-%d.%H:%M:%S).zip"
            ;;
         initiate)
            /home/webhook/go/bin/flixctl library initiate
            ;;
         inventory)
            /home/webhook/go/bin/flixctl library inventory \
                --enable-sync "${ENABLE_LIBRARY_SYNC}" \
                --job-id "${JOB_ID}"
            ;;
         jobs)
            /home/webhook/go/bin/flixctl library jobs \
                --filter "${FILTER}"
            ;;
         *)
            echo "'$1' is not a valid library command."
            echo "Usage: $0 {jobs|archive|download|initiate|inventory|jobs}"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {jobs|archive|download|initiate|inventory|jobs}"
      exit 3
      ;;
esac
