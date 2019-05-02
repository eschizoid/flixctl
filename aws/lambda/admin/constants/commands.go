package constants

var RenewCertsCommands = []string{
	`sudo openssl pkcs12 \
        -export \
        -password env:PLEX_PASSWORD \
        -out /var/lib/plexmediaserver/ssl/marianoflix.duckdns.org.pfx \
        -inkey /opt/ssl/marianoflix.duckdns.org/privkey.pem \
        -in /opt/ssl/marianoflix.duckdns.org/cert.pem \
        -certfile /opt/ssl/marianoflix.duckdns.org/fullchain.pem`,
	`sudo chown plex:plex /var/lib/plexmediaserver/ssl/marianoflix.duckdns.org.pfx`,
	`sudo /opt/dehydrated/dehydrated -c -o /opt/ssl`,
}

var RestartServicesCommand = "sudo systemctl restart %s"

var SlackCleanerCommand = `sudo /bin/slack-cleaner --perform \
    --quiet \
    --token %s \
    --rate 2 \
    --message \
    --channel %s \
    --bot \
    --user "*"`
