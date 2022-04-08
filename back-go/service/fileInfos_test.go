package service

import (
	"fmt"
	"testing"
)

func TestFileInfo(t *testing.T) {
	loadFilePath("./../cmd/cloud")
	fmt.Println()
	walk(filePtr.FileInfo, func(file *fileInfo) error {
		fmt.Println(file)
		return nil
	})

	copy2(filePtr.FileInfo.FileInfos["111"], filePtr.FileInfo, "333")
	copy2(filePtr.FileInfo.FileInfos["test.txt"], filePtr.FileInfo, "test2.txt")

	fmt.Println()
	walk(filePtr.FileInfo, func(file *fileInfo) error {
		fmt.Println(file)
		return nil
	})

}
