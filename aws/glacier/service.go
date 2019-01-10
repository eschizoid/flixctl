package glacier

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/glacier"
)

const (
	maxFileChunkSize      = 1024 * 1024 * 4 // 4MB
	maxTreeHashChunkSize  = 1024 * 1024     // 1MB
	inventoryRetrievalJob = "inventory-retrieval"
	archiveRetrievalJob   = "archive-retrieval"
)

func InitiateJob(svc *glacier.Glacier, archiveID string) *glacier.InitiateJobOutput {
	jobParameters := getJobParameters(archiveID)
	input := &glacier.InitiateJobInput{
		AccountId:     aws.String("-"),
		JobParameters: jobParameters,
		VaultName:     aws.String("plex"),
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
			fmt.Println(err.Error())
		}
		return nil
	}
	return result
}

func DescribeJob(svc *glacier.Glacier, jobID string) *glacier.JobDescription {
	input := &glacier.DescribeJobInput{
		AccountId: aws.String("-"),
		JobId:     aws.String(jobID),
		VaultName: aws.String("plex"),
	}
	result, err := svc.DescribeJob(input)
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
			fmt.Println(err.Error())
		}
		return nil
	}
	return result
}

func InitiateMultipartUploadInput(svc *glacier.Glacier, fileDescription string) *glacier.InitiateMultipartUploadOutput {
	input := &glacier.InitiateMultipartUploadInput{
		AccountId:          aws.String("-"),
		ArchiveDescription: aws.String(fileDescription),
		PartSize:           aws.String(strconv.Itoa(maxFileChunkSize)),
		VaultName:          aws.String("plex"),
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
			fmt.Println(err.Error())
		}
		return nil
	}
	return result
}

func UploadMultipartPartInput(svc *glacier.Glacier, uploadID string, fileChunkNames []string) []glacier.UploadMultipartPartOutput {
	var results = make([]glacier.UploadMultipartPartOutput, len(fileChunkNames))
	for i, name := range fileChunkNames {
		file, err := os.Open(name)
		showError(err)
		fs, _ := file.Stat()
		bytesRange := fmt.Sprintf("bytes %d-%d/*", i*maxFileChunkSize, (i*maxFileChunkSize)+(int(fs.Size())-1))
		fmt.Println("\n" + bytesRange)
		input := &glacier.UploadMultipartPartInput{
			AccountId: aws.String("-"),
			Body:      file,
			Checksum:  aws.String(ComputeTreeHash(name)),
			Range:     aws.String(bytesRange),
			UploadId:  aws.String(uploadID),
			VaultName: aws.String("plex"),
		}
		result, err := svc.UploadMultipartPart(input)
		if err != nil {
			abortMultipartUpload(svc, uploadID)
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case glacier.ErrCodeResourceNotFoundException:
					fmt.Println(glacier.ErrCodeResourceNotFoundException, aerr.Error())
				case glacier.ErrCodeInvalidParameterValueException:
					fmt.Println(glacier.ErrCodeInvalidParameterValueException, aerr.Error())
				case glacier.ErrCodeMissingParameterValueException:
					fmt.Println(glacier.ErrCodeMissingParameterValueException, aerr.Error())
				case glacier.ErrCodeRequestTimeoutException:
					fmt.Println(glacier.ErrCodeRequestTimeoutException, aerr.Error())
				case glacier.ErrCodeServiceUnavailableException:
					fmt.Println(glacier.ErrCodeServiceUnavailableException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return nil
		}
		results = append(results, *result)
	}
	return results
}

func CompleteMultipartUpload(svc *glacier.Glacier, uploadID string, fileName string) *glacier.ArchiveCreationOutput {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	stats, _ := file.Stat()
	input := &glacier.CompleteMultipartUploadInput{
		AccountId:   aws.String("-"),
		ArchiveSize: aws.String(strconv.FormatInt(stats.Size(), 10)),
		Checksum:    aws.String(ComputeTreeHash(fileName)),
		UploadId:    aws.String(uploadID),
		VaultName:   aws.String("plex"),
	}
	result, err := svc.CompleteMultipartUpload(input)
	if err != nil {
		abortMultipartUpload(svc, uploadID)
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
			fmt.Println(err.Error())
		}
		return nil
	}
	return result
}

func ListJobs(svc *glacier.Glacier) *glacier.ListJobsOutput {
	input := &glacier.ListJobsInput{
		AccountId: aws.String("-"),
		VaultName: aws.String("plex"),
	}
	result, err := svc.ListJobs(input)
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
			fmt.Println(err.Error())
		}
		return nil
	}
	return result
}

func GetJobOutput(svc *glacier.Glacier, jobID string) *glacier.GetJobOutputOutput {
	input := &glacier.GetJobOutputInput{
		AccountId: aws.String("-"),
		VaultName: aws.String("plex"),
		JobId:     aws.String(jobID),
	}
	result, err := svc.GetJobOutput(input)
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
			fmt.Println(err.Error())
		}
		return nil
	}
	return result
}

func abortMultipartUpload(svc *glacier.Glacier, uploadID string) {
	input := &glacier.AbortMultipartUploadInput{
		AccountId: aws.String("-"),
		UploadId:  aws.String(uploadID),
		VaultName: aws.String("plex"),
	}

	result, err := svc.AbortMultipartUpload(input)
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
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}

func getJobParameters(archiveID string) *glacier.JobParameters {
	var jobParameters *glacier.JobParameters
	if archiveID == "" {
		jobParameters = &glacier.JobParameters{
			Type: aws.String(inventoryRetrievalJob),
		}
	} else {
		jobParameters = &glacier.JobParameters{
			ArchiveId: aws.String(archiveID),
			Tier:      aws.String("Expedited"),
			Type:      aws.String(archiveRetrievalJob),
		}
	}
	return jobParameters
}
