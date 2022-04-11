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

	fmt.Println()
	// 根目录移动到下层
	copy2(filePtr.FileInfo.FileInfos["222"], filePtr.FileInfo.FileInfos["111"], "222")
	fmt.Println()
	// 下层目录移动到根目录
	copy2(filePtr.FileInfo.FileInfos["111"].FileInfos["333"], filePtr.FileInfo, "333")
	//copy2(filePtr.FileInfo.FileInfos["test.txt"], filePtr.FileInfo, "test2.txt")
	fmt.Println()
	// 同层移动 改名
	copy2(filePtr.FileInfo.FileInfos["666"], filePtr.FileInfo, "777")

	fmt.Println()
	// 下层目录移动到另一下层目录
	copy2(filePtr.FileInfo.FileInfos["111"].FileInfos["333"], filePtr.FileInfo.FileInfos["222"], "333")

	fmt.Println()
	walk(filePtr.FileInfo, func(file *fileInfo) error {
		fmt.Println(file)
		return nil
	})

}
