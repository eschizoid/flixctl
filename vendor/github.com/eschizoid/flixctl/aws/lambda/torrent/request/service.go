package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
)

func queueRequestItemCommand(evt json.RawMessage) {
	// search dynamo
	// search s3
	// post so SNS the ombi request
}

func main() {
	lambda.Start(queueRequestItemCommand)
}
