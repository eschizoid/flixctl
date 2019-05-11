package admin

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/eschizoid/flixctl/aws/s3"
	"github.com/eschizoid/flixctl/worker"
	"golang.org/x/crypto/ssh"
)

func GetSSHConnection() *ssh.Client {
	config := &ssh.ClientConfig{
		User: "centos",
		Auth: []ssh.AuthMethod{
			DownloadPublicKey(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", os.Getenv("FLIXCTL_HOST")+":22", config)
	if err != nil {
		panic(err)
	}
	return conn
}

func DownloadPublicKey() ssh.AuthMethod {
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

func RunCommand(cmd string, conn *ssh.Client) {
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

func AsyncCommandExecution(tasks []worker.TaskFunction) {
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
