#!/usr/bin/env bash

set -e
set -x
set -o pipefail

set +u
source /home/webhook/.bashrc
set -u

case $# in
   0)
      echo "Usage: $0 {archive|initiate|jobs|retrieve}"
      exit 1
      ;;
   1)
      case $1 in
         archive)
            /home/webhook/go/bin/flixctl library archive \
                --file "${FILE}"
            ;;
         initiate)
            /home/webhook/go/bin/flixctl library initiate
            ;;
         jobs)
            /home/webhook/go/bin/flixctl library jobs \
                --filter "${FILTER}"
            ;;
         retrieve)
            if [[ "${RETRIEVAL_TYPE}" = "ArchiveRetrieval" ]]; then
                /home/webhook/go/bin/flixctl library retrieve \
                    --job-id "${JOB_ID}" \
                    --file "/plex/movies/movie-$(date +%Y-%m-%d.%H:%M:%S).zip" \
                    --type "${RETRIEVAL_TYPE}"
            elif [[ "${RETRIEVAL_TYPE}" = "InventoryRetrieval" ]]; then
                /home/webhook/go/bin/flixctl library retrieve \
                    --job-id "${JOB_ID}" \
                    --type "${RETRIEVAL_TYPE}"
            else
                "Unknown parameter"
            fi
            ;;
         *)
            echo "'$1' is not a valid library command."
            echo "Usage: $0 {jobs|archive|initiate|retrieve}"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {archive|initiate|jobs|retrieve}"
      exit 3
      ;;
esac
