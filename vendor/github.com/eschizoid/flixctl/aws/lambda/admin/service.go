package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
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
			publicKey(os.Getenv("FLIXCTL_PUBLIC_KEY_PATH")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, _ := ssh.Dial("tcp", os.Getenv("FLIXCTL_HOST"), config)
	defer conn.Close()

	switch input.Command {
	case "renew-certs":
		runCommand("", conn)
	case "restart-services":
		runCommand("", conn)
	case "purge-slack":
		runCommand("", conn)
	}
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

func publicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func main() {
	lambda.Start(executeAdminCommand)
}
