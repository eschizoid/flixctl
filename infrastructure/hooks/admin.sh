#!/usr/bin/env bash

set -e
#set -x
set -o pipefail

set +u
source /home/webhook/.bashrc
set -u

case $# in
   0)
      echo "Usage: $0 {endpoints|metrics|purge-slack|renew-certs|upgrade}"
      exit 1
      ;;
   1)
      case $1 in
         endpoints)
            HOOKS="$(cat /opt/webhook-linux-amd64/hooks.json | grep -o 'id.*' | cut -f2- -d: | sort)"
            HOOKS="$(echo ${HOOKS::-1})"
            echo "[${HOOKS}]"
            ;;
         metrics)
            ;;
         purge-slack)
            for channel in monitoring; do
                /bin/slack-cleaner --perform \
                    --quiet \
                    --token "${SLACK_LEGACY_TOKEN}" \
                    --message \
                    --group ${channel} \
                    --bot
                sleep 5
                /bin/slack-cleaner --perform \
                    --quiet \
                    --token "${SLACK_LEGACY_TOKEN}" \
                    --message \
                    --group ${channel} \
                    --user "*"
                sleep 5
            done
            for channel in new-releases requests travis; do
                /bin/slack-cleaner --perform \
                    --quiet \
                    --token "${SLACK_LEGACY_TOKEN}" \
                    --message \
                    --channel ${channel} \
                    --bot
                sleep 5
                /bin/slack-cleaner --perform \
                    --quiet \
                    --token "${SLACK_LEGACY_TOKEN}" \
                    --message \
                    --channel ${channel} \
                    --user "*"
                sleep 5
            done
            echo "{\"slack_purged\": \"true\"}"
            ;;
         renew-certs)
            openssl pkcs12 -export \
                -password env:PLEX_PASSWORD \
                -out /var/lib/plexmediaserver/ssl/marianoflix.duckdns.org.pfx \
                -inkey /opt/webhook-linux-amd64/privkey.pem \
                -in /opt/webhook-linux-amd64/cert.pem \
                -certfile /opt/webhook-linux-amd64/fullchain.pem
            /opt/dehydrated/dehydrated -c -o /opt/ssl
            echo "{\"certificates_updated\": \"true\"}"
            ;;
         upgrade)
            rm -rf /home/webhook/go/src/github.com/eschizoid/flixctl
            /usr/local/go/bin/go get -u github.com/eschizoid/flixctl
            cd /home/webhook/go/src/github.com/eschizoid/flixctl
            cp -r infrastructure/hooks/*.sh /opt/webhook-linux-amd64/
            /bin/make install
            /home/webhook/go/bin/flixctl version
            ;;
         *)
            echo "'$1' is not a valid admin command."
            echo "Usage: $0 {endpoints|metrics|purge-slack|renew-certs|upgrade}"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {endpoints|metrics|purge-slack|renew-certs|upgrade}"
      exit 3
      ;;
esac
