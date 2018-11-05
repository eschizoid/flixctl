package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/glacier"
)

func CreateVault(svc *glacier.Glacier) {
	input := &glacier.CreateVaultInput{
		AccountId: aws.String("-"),
		VaultName: aws.String("plex"),
	}
	result, err := svc.CreateVault(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case glacier.ErrCodeInvalidParameterValueException:
				fmt.Println(glacier.ErrCodeInvalidParameterValueException, aerr.Error())
			case glacier.ErrCodeMissingParameterValueException:
				fmt.Println(glacier.ErrCodeMissingParameterValueException, aerr.Error())
			case glacier.ErrCodeServiceUnavailableException:
				fmt.Println(glacier.ErrCodeServiceUnavailableException, aerr.Error())
			case glacier.ErrCodeLimitExceededException:
				fmt.Println(glacier.ErrCodeLimitExceededException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}

func DownloadFile(svc *glacier.Glacier) {
	input := &glacier.InitiateJobInput{
		AccountId: aws.String("-"),
		JobParameters: &glacier.JobParameters{
			Description: aws.String("My inventory job"),
			Format:      aws.String("CSV"),
			SNSTopic:    aws.String("arn:aws:sns:us-west-2:111111111111:Glacier-InventoryRetrieval-topic-Example"),
			Type:        aws.String("inventory-retrieval"),
		},
		VaultName: aws.String("examplevault"),
	}

	result, err := svc.InitiateJob(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case glacier.ErrCodeResourceNotFoundException:
				fmt.Println(glacier.ErrCodeResourceNotFoundException, aerr.Error())
			case glacier.ErrCodePolicyEnforcedException:
				fmt.Println(glacier.ErrCodePolicyEnforcedException, aerr.Error())
			case glacier.ErrCodeInvalidParameterValueException:
				fmt.Println(glacier.ErrCodeInvalidParameterValueException, aerr.Error())
			case glacier.ErrCodeMissingParameterValueException:
				fmt.Println(glacier.ErrCodeMissingParameterValueException, aerr.Error())
			case glacier.ErrCodeInsufficientCapacityException:
				fmt.Println(glacier.ErrCodeInsufficientCapacityException, aerr.Error())
			case glacier.ErrCodeServiceUnavailableException:
				fmt.Println(glacier.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

// To initiate a multipart upload
//
// The example initiates a multipart upload to a vault named my-vault with a part size
// of 1 MiB (1024 x 1024 bytes) per file.
func UploadFile(svc *glacier.Glacier) {
	input := &glacier.InitiateMultipartUploadInput{
		AccountId: aws.String("-"),
		PartSize:  aws.String("1048576"),
		VaultName: aws.String("plex"),
	}
	result, err := svc.InitiateMultipartUpload(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case glacier.ErrCodeResourceNotFoundException:
				fmt.Println(glacier.ErrCodeResourceNotFoundException, aerr.Error())
			case glacier.ErrCodeInvalidParameterValueException:
				fmt.Println(glacier.ErrCodeInvalidParameterValueException, aerr.Error())
			case glacier.ErrCodeMissingParameterValueException:
				fmt.Println(glacier.ErrCodeMissingParameterValueException, aerr.Error())
			case glacier.ErrCodeServiceUnavailableException:
				fmt.Println(glacier.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}
