package service

import (
	"fmt"
	"testing"
)

func TestFileInfo(t *testing.T) {
	loadFilePath("../cmd/log")

	walk(filePtr.FileInfo, func(file *fileInfo) error {
		fmt.Println(file)
		return nil
	})

}
