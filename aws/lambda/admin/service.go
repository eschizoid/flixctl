package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/eschizoid/flixctl/aws/lambda/admin/constants"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/aws/s3"
	"golang.org/x/crypto/ssh"
)

func executeAdminCommand(evt json.RawMessage) {
	var input models.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	fmt.Printf("Exectuing λ with payload: %+v\n", input)

	config := &ssh.ClientConfig{
		User: "centos",
		Auth: []ssh.AuthMethod{
			downloadPublicKey(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", os.Getenv("FLIXCTL_HOST")+":22", config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	switch input.Argument {
	case "renew-certs":
		fmt.Printf("Executing %s command", input.Argument)
		for _, command := range constants.RenewCertsCommands {
			runCommand(command, conn)
		}
		fmt.Printf("Succesfully executed %s", input.Argument)
	case "restart-services":
		fmt.Printf("Executing %s command", input.Argument)
		services := []string{"httpd", "jackett", "nzbget", "ombi", "plexmediaserver", "radarr", "sonarr", "s3fs", "tautulli", "transmission-daemon"}
		for _, service := range services {
			fmt.Printf("Exectuing λ with payload: %+v\n", input)
			runCommand(fmt.Sprintf(constants.RestartServicesCommand, service), conn)
			fmt.Printf("Succesfully restarted %s service", service)
		}
		fmt.Printf("Succesfully executed %s", input.Argument)
	case "purge-slack":
		fmt.Printf("Executing %s command", input.Argument)
		slackChannels := []string{"monitoring", "new-releases", "requests", "travis"}
		for _, channel := range slackChannels {
			runCommand(fmt.Sprintf(constants.SlackCleanerCommands[0], os.Getenv("SLACK_LEGACY_TOKEN"), channel), conn)
			time.Sleep(10 * time.Second)
			runCommand(fmt.Sprintf(constants.SlackCleanerCommands[1], os.Getenv("SLACK_LEGACY_TOKEN"), channel), conn)
			time.Sleep(10 * time.Second)
		}
		fmt.Printf("Succesfully executed %s", input.Argument)
	}
}

func downloadPublicKey() ssh.AuthMethod {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	downloader := s3manager.NewDownloader(awsSession)
	sshKey := s3.DownloadItem(downloader, "marianoflix", "certicates/marianoflix.pem", "/tmp/marianoflix.pem")
	key, err := ioutil.ReadFile(sshKey.Name())
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func runCommand(cmd string, conn *ssh.Client) {
	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close() //nolint:errcheck
	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, sessStdOut) //nolint:errcheck
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, sessStderr) //nolint:errcheck
	err = sess.Run(cmd)
	if err != nil {
		panic(err)
	}
}

func main() {
	lambda.Start(executeAdminCommand)
}
