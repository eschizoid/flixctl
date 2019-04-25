package constants

var RenewCertsCommands = []string{
	`openssl pkcs12 \
        -export \
        -password env:PLEX_PASSWORD \
        -out /var/lib/plexmediaserver/ssl/marianoflix.duckdns.org.pfx \
        -inkey /opt/ssl/marianoflix.duckdns.org/privkey.pem \
        -in /opt/ssl/marianoflix.duckdns.org/cert.pem \
        -certfile /opt/ssl/marianoflix.duckdns.org/fullchain.pem`,
	`sudo chown plex:plex /var/lib/plexmediaserver/ssl/marianoflix.duckdns.org.pfx`,
	`/opt/dehydrated/dehydrated -c -o /opt/ssl`,
}

var RestartServicesCommands = []string{
	`for plex_service in httpd \
        jackett \
        nzbget \
        ombi \
        plexmediaserver \
        radarr \
        sonarr \
        s3fs \
        tautulli \
        transmission-daemon; do \
        sudo systemctl restart ${plex_service} \
    done`,
}

var SlackCleanerCommands = []string{
	`for channel in monitoring; do \
        /bin/slack-cleaner --perform \
            --quiet \
            --token "${SLACK_LEGACY_TOKEN}" \
            --message \
            --group ${channel} \
            --bot \
    done`,
	`sleep 5`,
	`for channel in monitoring; do \
        /bin/slack-cleaner --perform \
            --quiet \
            --token "${SLACK_LEGACY_TOKEN}" \
            --message \
            --group ${channel} \
            --user "*" \
    done`,
	`sleep 5`,
	`for channel in new-releases requests travis; do \
        /bin/slack-cleaner --perform \
            --quiet \
            --token "${SLACK_LEGACY_TOKEN}" \
            --message \
            --channel ${channel} \
        --bot \
    done`,
	`sleep 5`,
	`for channel in new-releases requests travis; do \
        /bin/slack-cleaner --perform \
            --quiet \
            --token "${SLACK_LEGACY_TOKEN}" \
            --message \
            --channel ${channel} \
            --user "*" \
        done`,
	`sleep 5`,
}
