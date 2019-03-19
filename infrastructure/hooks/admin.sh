#!/usr/bin/env bash

set -e
set -o pipefail

set +u
source /home/webhook/.bashrc
set -u

case $# in
   0)
      echo "Usage: $0 {endpoints|metrics|purge-slack|maybe-shutdown|renew-certs|restart|upgrade"
      exit 1
      ;;
   1)
      case $1 in
         endpoints)
            HOOKS="$(cat /opt/webhook-linux-amd64/hooks.json | grep -o 'id.*' | cut -f2- -d: | sort)"
            HOOKS="$(echo ${HOOKS::-1})"
            echo "[${HOOKS}]"
            ;;
         maybe-shutdown)
            /home/webhook/go/bin/flixctl plex monitor \
                --slack-notification "${SLACK_NOTIFICATION}" \
                --max-inactive-time 30
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
                -inkey /opt/ssl/marianoflix.duckdns.org/privkey.pem \
                -in /opt/ssl/marianoflix.duckdns.org/cert.pem \
                -certfile /opt/ssl/marianoflix.duckdns.org/fullchain.pem
            sudo chown plex:plex /var/lib/plexmediaserver/ssl/marianoflix.duckdns.org.pfx
            /opt/dehydrated/dehydrated -c -o /opt/ssl
            echo "{\"certificates_updated\": \"true\"}"
            ;;
         restart)
            for plex_service in httpd \
                jackett \
                nzbget \
                ombi \
                plexmediaserver \
                radarr \
                sonarr \
                s3fs \
                tautulli \
                transmission-daemon; do
                sudo systemctl restart ${plex_service}
            done
            echo "{\"services_restarted\": \"true\"}"
            ;;
         upgrade)
            rm -rf /home/webhook/go/src/github.com/eschizoid/flixctl
            /usr/local/go/bin/go get -u github.com/eschizoid/flixctl
            cd /home/webhook/go/src/github.com/eschizoid/flixctl
            cp -r infrastructure/hooks/{*.sh,*.json} /opt/webhook-linux-amd64/
            /bin/make install
            /home/webhook/go/bin/flixctl version
            ;;
         *)
            echo "'$1' is not a valid admin command."
            echo "Usage: $0 {endpoints|metrics|purge-slack|maybe-shutdown|renew-certs|restart|upgrade"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {endpoints|metrics|purge-slack|maybe-shutdown|renew-certs|restart|upgrade"
      exit 3
      ;;
esac
