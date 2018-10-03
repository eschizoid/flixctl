package snapshot

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Create(sess *sess.Session, volumeID string, name string) {
	svc := ec2.New(sess, sess.Config)
	tagList := &ec2.TagSpecification{
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
		},
		ResourceType: aws.String(ec2.ResourceTypeSnapshot),
	}
	input := &ec2.CreateSnapshotInput{
		Description:       aws.String("Plex Snapshot"),
		VolumeId:          aws.String(volumeID),
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
	describeSnapshotsInput := &ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(name)},
			},
		},
	}
	if err := svc.WaitUntilSnapshotCompleted(describeSnapshotsInput); err != nil {
		panic(err)
	}
}

func FetchSnapshotID(sess *sess.Session, name string) string {
	var snapshotID string
	svc := ec2.New(sess, sess.Config)
	params := &ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(name)},
			},
		},
	}
	result, err := svc.DescribeSnapshots(params)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, snapshots := range result.Snapshots {
		snapshotID = aws.StringValue(snapshots.SnapshotId)
	}
	return snapshotID
}

func Delete(sess *sess.Session, snapshotID string) {
	svc := ec2.New(sess, sess.Config)
	input := &ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshotID),
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
