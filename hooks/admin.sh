#!/usr/bin/env bash

case $# in
   0)
      echo "Usage: $0 {endpoints|renew-ssl-certificates|upgrade}"
      exit 1
      ;;
   1)
      case $1 in
         endpoints)
            HOOKS="$(cat /opt/webhook-linux-amd64/hooks.json | grep -o 'id.*' | cut -f2- -d:)"
            HOOKS="$(echo ${HOOKS::-1})"
            echo "[${HOOKS}]"
            ;;
         renew-ssl-certificates)
            openssl pkcs12 -export \
                -password env:PLEX_PASSWORD \
                -out /var/lib/plexmediaserver/ssl/marianoflix.duckdns.org.pfx \
                -inkey /opt/webhook-linux-amd64/privkey.pem \
                -in /opt/webhook-linux-amd64/cert.pem \
                -certfile /opt/webhook-linux-amd64/fullchain.pem
            /opt/dehydrated/dehydrated -c -o /opt/ssl
            echo "{\"certificates_upgraded\": \"true\"}"
            ;;
         upgrade)
            rm -rf /home/webhook/go/src/github.com/eschizoid/flixctl
            /usr/local/go/bin/go get -u github.com/eschizoid/flixctl
            cd /home/webhook/go/src/github.com/eschizoid/flixctl
            /bin/make install
            /home/webhook/go/bin/flixctl version
            ;;
         *)
            echo "'$1' is not a valid admin command."
            echo "Usage: $0 {endpoints|renew-ssl-certificates|upgrade}"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {endpoints|upgrade}"
      exit 3
      ;;
esac
