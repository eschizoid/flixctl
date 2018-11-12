package snapshot

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Create(svc *ec2.EC2, volumeID string, name string) {
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
			fmt.Println(aerr.Error())
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

func FetchSnapshotID(svc *ec2.EC2, name string) string {
	var snapshotID string
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

func Delete(svc *ec2.EC2, snapshotID string) {
	input := &ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshotID),
	}
	_, err := svc.DeleteSnapshot(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return
	}
}
