package s3

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func DownloadItem(downloader *s3manager.Downloader, bucket string, key string, destination string) *os.File {
	file, err := os.Create(destination)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", destination, err)
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return file
}

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
