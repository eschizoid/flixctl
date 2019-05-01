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

var RestartServicesCommand = "sudo systemctl restart %s"

var SlackCleanerCommands = []string{
	`/bin/slack-cleaner --perform \
        --quiet \
        --token %s \
        --message \
        --group %s \
        --bot`,
	`/bin/slack-cleaner --perform \
        --quiet \
        --token %s \
        --message \
        --group %s \
        --user "*"`,
}
