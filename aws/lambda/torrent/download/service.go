package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
)

func queueDownloadCommand(evt json.RawMessage) {
	// search dynamo
	// search s3
	// post to SNS magnet link
}

func main() {
	lambda.Start(queueDownloadCommand)
}
