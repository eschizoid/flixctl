package ec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const DryRunOperation = "DryRunOperation"

func Start(svc *ec2.EC2, id string) {
	instanceID := aws.StringSlice([]string{id})
	input := &ec2.StartInstancesInput{
		InstanceIds: instanceID,
		DryRun:      aws.Bool(true),
	}
	_, err := svc.StartInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == DryRunOperation {
		input.DryRun = aws.Bool(false)
		_, err = svc.StartInstances(input)
		if err != nil {
			fmt.Println("Error", err)
		} else {
			describeInstancesInput := &ec2.DescribeInstancesInput{
				InstanceIds: instanceID,
			}
			if aerr := svc.WaitUntilInstanceRunning(describeInstancesInput); aerr != nil {
				panic(aerr)
			}
		}
	} else {
		fmt.Println("Error", err)
	}
}

func Stop(svc *ec2.EC2, id string) {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
		DryRun: aws.Bool(true),
	}
	_, err := svc.StopInstances(input)
	awsErr, ok := err.(awserr.Error)
	if ok && awsErr.Code() == DryRunOperation {
		input.DryRun = aws.Bool(false)
		_, err = svc.StopInstances(input)
		if err != nil {
			fmt.Println("Error", err)
		} else {
			describeInstancesInput := &ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice([]string{id}),
			}
			if aerr := svc.WaitUntilInstanceStopped(describeInstancesInput); aerr != nil {
				panic(aerr)
			}
		}
	} else {
		fmt.Println("Error", err)
	}
}

func Status(svc *ec2.EC2, id string) string {
	var status string
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
			}
		}
	}
	return status
}

func FetchInstanceID(svc *ec2.EC2, name string) string {
	var instanceID string
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
			instanceID = *instance.InstanceId
			break
		}
	}
	return instanceID
}
