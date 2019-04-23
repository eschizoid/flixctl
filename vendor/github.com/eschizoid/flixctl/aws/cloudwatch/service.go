package cloudwatch

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

func DisableRule(svc *cloudwatchevents.CloudWatchEvents, ruleName string) {
	input := &cloudwatchevents.DisableRuleInput{
		Name: aws.String(ruleName),
	}
	_, err := svc.DisableRule(input)
	if err != nil {
		panic(err)
	}
}

func EnableRule(svc *cloudwatchevents.CloudWatchEvents, ruleName string) {
	input := &cloudwatchevents.EnableRuleInput{
		Name: aws.String(ruleName),
	}
	_, err := svc.EnableRule(input)
	if err != nil {
		panic(err)
	}
}
