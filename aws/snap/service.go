package ebs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	. "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Create(sess *Session, volumeId string) {
	svc := ec2.New(sess)
	tagList := &ec2.TagSpecification{
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("plex"),
			},
		},
		ResourceType: aws.String(ec2.ResourceTypeSnapshot),
	}
	input := &ec2.CreateSnapshotInput{
		Description:       aws.String("Plex Snapshot"),
		VolumeId:          aws.String(volumeId),
		TagSpecifications: []*ec2.TagSpecification{tagList},
	}

	_, err := svc.CreateSnapshot(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}
}

func Find(sess *Session, name string) string {
	var snapshotId string
	svc := ec2.New(sess)
	params := &ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(name)},
			},
		},
	}
	result, err := svc	.DescribeSnapshots(params)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, snapshots := range result.Snapshots {
		snapshotId = aws.StringValue(snapshots.SnapshotId)
	}
	return snapshotId
}

func Delete(sess *Session, snapshotId string) {
	svc := ec2.New(sess)
	input := &ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshotId),
	}
	_, err := svc.DeleteSnapshot(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}
}
