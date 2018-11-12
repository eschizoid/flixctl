package glacier

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/service/glacier"
)

func Chunk(fileName string) []string {
	var files []string
	file, err := os.Open(fileName)
	ShowError(err)
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
		if err != nil {
			fmt.Println(err)
		}
		// write to disk
		fileName := "part-" + strconv.FormatUint(i, 10)
		_, err = os.Create(fileName)
		if err != nil {
			fmt.Println(err)
		}
		// write/save buffer to disk
		err = ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)
		if err != nil {
			fmt.Println(err)
		}
		files = append(files, fileName)
	}
	return files
}

func ComputeTreeHash(fileName string) string {
	file, err := os.Open(fileName)
	ShowError(err)
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

func ShowError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
