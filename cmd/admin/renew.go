package admin

import (
	"github.com/spf13/cobra"
)

var RemoteRenewCerts = []string{
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

var RenewCertsCmd = &cobra.Command{
	Use:   "renew-certs",
	Short: "To Renew Certs",
	Long:  "to renew tls certificates all plex related services",
	Run: func(cmd *cobra.Command, args []string) {
		RenewCerts()
	},
}

func RenewCerts() {
	conn := GetSSHConnection()
	defer conn.Close()
	for _, command := range RemoteRenewCerts {
		RunCommand(command, conn)
	}
}
