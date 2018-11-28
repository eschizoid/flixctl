package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func dispatch() {
}

func main() {
	lambda.Start(dispatch)
}
