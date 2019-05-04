package main

import (
	"context"
	"encoding/json"

	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/eschizoid/flixctl/aws/lambda/admin/constants"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/aws/s3"
	auth "github.com/eschizoid/flixctl/cmd/auth"
	"github.com/eschizoid/flixctl/worker"
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
	var tasks []worker.TaskFunction
	switch input.Argument {
	case "oauth-token":
		fmt.Printf("Executing %s command \n", input.Argument)
		message := fmt.Sprintf("Succesfully executed oauth-tokens")
		commandTask := func() interface{} {
			auth.GetOauthToke()
			return message
		}
		asyncCommandExecution(append(tasks, commandTask))
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "purge-slack":
		fmt.Printf("Executing %s command \n", input.Argument)
		slackChannels := []string{"monitoring", "new-releases", "requests", "travis"}
		for _, channel := range slackChannels {
			command := fmt.Sprintf(constants.SlackCleanerCommand, os.Getenv("SLACK_LEGACY_TOKEN"), channel)
			message := fmt.Sprintf("Succesfully purged slack channel %s", channel)
			commandTask := func() interface{} {
				runCommand(command, conn)
				return message
			}
			tasks = append(tasks, commandTask)
		}
		asyncCommandExecution(tasks)
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "renew-certs":
		fmt.Printf("Executing %s command \n", input.Argument)
		for _, command := range constants.RenewCertsCommands {
			runCommand(command, conn)
		}
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "restart-services":
		fmt.Printf("Executing %s command \n", input.Argument)
		services := []string{"httpd", "jackett", "nzbget", "ombi", "plexmediaserver", "radarr", "sonarr", "s3fs", "tautulli", "transmission-daemon"}
		for _, service := range services {
			command := fmt.Sprintf(constants.RestartServicesCommand, service)
			message := fmt.Sprintf("Succesfully restarted service %s", service)
			commandTask := func() interface{} {
				runCommand(command, conn)
				return message
			}
			tasks = append(tasks, commandTask)
		}
		asyncCommandExecution(tasks)
		fmt.Printf("Succesfully executed %s \n", input.Argument)
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

func asyncCommandExecution(tasks []worker.TaskFunction) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resultChannel := worker.PerformTasks(ctx, tasks)
	for result := range resultChannel {
		switch v := result.(type) {
		case error:
			fmt.Println(v)
		case string:
			fmt.Println(v)
		default:
			fmt.Println("Some unknown type ")
		}
	}
}

func main() {
	lambda.Start(executeAdminCommand)
}
