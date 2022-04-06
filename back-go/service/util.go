package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

func nowFormat() string {
	return time.Now().Format(timeFormat)
}

func makeFilePart(name, part string) string {
	return fmt.Sprintf("%s.part%s", name, part)
}

// 文件 md5 值
func fileMD5(filename string) (string, error) {
	if info, err := os.Stat(filename); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", errors.New(fmt.Sprintf("%s is a dir", filename))
	}

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func splitPath(dir string) []string {
	paths := strings.Split(dir, "/")
	l := []string{}
	for _, v := range paths {
		if v != "" {
			l = append(l, v)
		}
	}
	return l
}

func Must(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return i
}

func WriteFile(filename string, reader io.Reader) (written int64, err error) {
	newFile, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer newFile.Close()

	return io.Copy(newFile, reader)
}

func CopyFile(src, dest string) (written int64, err error) {
	srcF, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcF.Close()

	return WriteFile(dest, srcF)
}

var (
	KB = uint64(math.Pow(2, 10))
	MB = uint64(math.Pow(2, 20))
	GB = uint64(math.Pow(2, 30))
	TB = uint64(math.Pow(2, 40))
)

func CelsiusToFahrenheit(c int) int {
	return c*9/5 + 32
}

func BytesToKB(b uint64) float64 {
	return float64(b) / float64(KB)
}

func BytesToMB(b uint64) float64 {
	return float64(b) / float64(MB)
}

func BytesToGB(b uint64) float64 {
	return float64(b) / float64(GB)
}

func BytesToTB(b uint64) float64 {
	return float64(b) / float64(TB)
}

func ConvertBytes(b uint64) (float64, string) {
	switch {
	case b < KB:
		return float64(b), "B"
	case b < MB:
		return BytesToKB(b), "KB"
	case b < GB:
		return BytesToMB(b), "MB"
	case b < TB:
		return BytesToGB(b), "GB"
	default:
		return BytesToTB(b), "TB"
	}
}

func ConvertBytesString(b uint64) string {
	cf, s := ConvertBytes(b)
	return fmt.Sprintf("%.1f%s", cf, s)
}
