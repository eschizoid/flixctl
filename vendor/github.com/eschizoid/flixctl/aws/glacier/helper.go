package glacier

import (
	"compress/flate"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/mholt/archiver"
)

func Chunk(fileName string) []string {
	var files []string
	file, err := os.Open(fileName)
	showError(err)
	defer file.Close()
	fileInfo, _ := file.Stat()
	var fileSize = fileInfo.Size()
	// calculate total number of parts the file will be chunked into
	totalParts := uint64(math.Ceil(float64(fileSize) / float64(maxFileChunkSize)))
	fmt.Printf("Splitting to %d pieces.\n", totalParts)
	for i := uint64(0); i < totalParts; i++ {
		partSize := int(math.Min(maxFileChunkSize, float64(fileSize-int64(i*maxFileChunkSize))))
		fmt.Printf("Part size: %d\n", partSize)
		partBuffer := make([]byte, partSize)
		_, err = file.Read(partBuffer)
		showError(err)
		// write to disk
		partFileName := fmt.Sprintf("/tmp/part-%s", strconv.FormatUint(i, 10))
		_, err = os.Create(partFileName)
		showError(err)
		// write/save buffer to disk
		err = ioutil.WriteFile(partFileName, partBuffer, os.ModeAppend)
		showError(err)
		files = append(files, partFileName)
	}
	return files
}

func Cleanup(files []string) {
	for _, fileChunk := range files {
		err := os.Remove(fileChunk)
		showError(err)
	}
}

func ComputeTreeHash(fileName string) string {
	file, err := os.Open(fileName)
	showError(err)
	defer file.Close()
	buf := make([]byte, maxTreeHashChunkSize)
	var hashes [][]byte
	for {
		n, err := io.ReadAtLeast(file, buf, maxTreeHashChunkSize)
		if n == 0 {
			break
		}
		tmpHash := sha256.Sum256(buf[:n])
		hashes = append(hashes, tmpHash[:])
		if err != nil {
			break // last chunk
		}
	}
	treeHash := fmt.Sprintf("%x", glacier.ComputeTreeHash(hashes))
	fmt.Printf("TreeHash: %s\n", treeHash)
	return treeHash
}

func GetStats(fileName string) os.FileInfo {
	file, err := os.Open(fileName)
	showError(err)
	defer file.Close()
	stats, err := file.Stat()
	showError(err)
	return stats
}

func Unzip(source string) {
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,
	}
	err := z.Unarchive(source, "/tmp")
	showError(err)
}

func Zip(source string) os.File {
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,
	}
	file, err := ioutil.TempFile("/tmp", "movie.*.zip")
	showError(err)
	err = z.Archive([]string{source}, file.Name())
	showError(err)
	return *file
}

func showError(err error) {
	if err != nil {
		panic(err)
	}
}
