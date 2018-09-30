package ebs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	. "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Attach(sess *Session, instanceId string, volumeId string) {
	svc := ec2.New(sess)
	input := &ec2.AttachVolumeInput{
		Device:     aws.String("/dev/sdf"),
		InstanceId: aws.String(instanceId),
		VolumeId:   aws.String(volumeId),
	}
	_, err := svc.AttachVolume(input)
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

func Detach(sess *Session, volumeId string) {
	svc := ec2.New(sess)
	input := &ec2.DetachVolumeInput{
		VolumeId: aws.String(volumeId),
	}
	_, err := svc.DetachVolume(input)
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

func Create(sess *Session, snapshotId string) {
	svc := ec2.New(sess)
	tagList := &ec2.TagSpecification{
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("plex"),
			},
		},
		ResourceType: aws.String(ec2.ResourceTypeVolume),
	}
	input := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String("us-east-1a"),
		SnapshotId:       aws.String(snapshotId),
		VolumeType:       aws.String("sc1"),
		TagSpecifications: []*ec2.TagSpecification{tagList},
	}
	_, err := svc.CreateVolume(input)
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

func Delete(sess *Session, volumeId string) {
	svc := ec2.New(sess)
	input := &ec2.DeleteVolumeInput{
		VolumeId: aws.String(volumeId),
	}
	_, err := svc.DeleteVolume(input)
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

func FetchVolumeId(sess *Session, name string) string {
	var volumeId string
	svc := ec2.New(sess)
	params := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(name)},
			},
		},
	}
	result, err := svc.DescribeVolumes(params)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, volumes := range result.Volumes {
		volumeId = aws.StringValue(volumes.VolumeId)
	}
	return volumeId
}
