package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

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

	config := &ssh.ClientConfig{
		User: "username",
		Auth: []ssh.AuthMethod{
			downloadPublicKey(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, _ := ssh.Dial("tcp", os.Getenv("FLIXCTL_HOST")+":22", config)
	defer conn.Close()

	switch input.Command {
	case "renew-certs":
		for _, command := range constants.RenewCertsCommands {
			runCommand(command, conn)
		}
	case "restart-services":
		for _, command := range constants.RestartServicesCommands {
			runCommand(command, conn)
		}
	case "purge-slack":
		for _, command := range constants.SlackCleanerCommands {
			runCommand(command, conn)
		}
	}
}

func downloadPublicKey() ssh.AuthMethod {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	downloader := s3manager.NewDownloader(awsSession)
	sshKey := s3.DownloadItem(downloader, "marianoflix", "certicates/marianoflix.pem")
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
