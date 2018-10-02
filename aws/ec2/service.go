package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	. "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

func Start(sess *Session, id string) {
	svc := ec2.New(sess)
	instanceId := aws.StringSlice([]string{id})
	input := &ec2.StartInstancesInput{
		InstanceIds: instanceId,
		DryRun: aws.Bool(true),
	}
	_, err := svc.StartInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		_, err = svc.StartInstances(input)
		if err != nil {
			fmt.Println("Error", err)
		} else {
			describeInstancesInput := &ec2.DescribeInstancesInput{
				InstanceIds: instanceId,
			}
			if err := svc.WaitUntilInstanceRunning(describeInstancesInput); err != nil {
				panic(err)
			}
		}
	} else {
		fmt.Println("Error", err)
	}
}

func Stop(sess *Session, id string) {
	svc := ec2.New(sess)
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
		DryRun: aws.Bool(true),
	}
	_, err := svc.StopInstances(input)
	awsErr, ok := err.(awserr.Error)
	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		_, err = svc.StopInstances(input)
		if err != nil {
			fmt.Println("Error", err)
		} else {
			describeInstancesInput := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice([]string{id}),
			}
			if err := svc.WaitUntilInstanceStopped(describeInstancesInput); err != nil {
				panic(err)
			}
		}
	} else {
		fmt.Println("Error", err)
	}
}

func Status(sess *Session, id string) (string) {
	var status string
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
		DryRun: aws.Bool(false),
	}
	result, err := svc.DescribeInstances(input)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				status = strings.Title(*instance.State.Name)
				fmt.Println("Plex current status: " + status)
			}
		}
	}
	return status
}

func FetchInstanceId(sess *Session, name string) string {
	var instanceId string
	svc := ec2.New(sess)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(name)},
			},
		},
	}
	result, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceId = *instance.InstanceId
			break
		}
	}
	return instanceId
}
