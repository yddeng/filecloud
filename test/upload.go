package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

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

func main() {
	filePath := os.Args[1]
	//uploadPath := os.Args[1]

	md5, err := fileMD5(filePath)
	if err != nil {
		panic(err)
	}

	uploadPath := "cloud"
	if err := fileCheck(uploadPath, filePath, md5); err != nil {
		panic(err)
	}

	if err := fileUpload(uploadPath, filePath, md5); err != nil {
		panic(err)
	}

}

func fileCheck(uploadPath, filePath, md5 string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	cnt := len(data) / (2 * 1024 * 1024)
	if len(data)%(2*1024*1024) != 0 {
		cnt++
	}

	_, filename := filepath.Split(filePath)
	elem := map[string]interface{}{
		"path":     uploadPath,
		"filename": filename,
		"md5":      md5,
		"total":    cnt,
		"size":     len(data),
	}

	req, _ := dhttp.PostJson("http://127.0.0.1:9987/file/check", elem)
	ret, err := req.ToString()
	fmt.Println(ret)
	return err
}

func fileUpload(uploadPath, filePath, md5 string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	cnt := len(data) / (2 * 1024 * 1024)
	if len(data)%(2*1024*1024) != 0 {
		cnt++
	}

	_, filename := filepath.Split(filePath)

	wg := sync.WaitGroup{}
	wg.Add(cnt)

	for i := 0; i < cnt; i++ {
		go func(i int) {
			req, _ := http.NewRequest("POST", "http://127.0.0.1:9987/file/upload", nil)

			buf := new(bytes.Buffer)
			writer := multipart.NewWriter(buf)
			fileWriter, err := writer.CreateFormFile("file", filename)
			if err != nil {
				return
			}

			start := i * (1024 * 1024 * 2)
			end := (i + 1) * (1024 * 1024 * 2)
			if end > len(data) {
				end = len(data)
			}
			_, err = fileWriter.Write(data[start:end])
			if err != nil {
				return
			}
			writer.WriteField("path", uploadPath)
			writer.WriteField("md5", md5)
			writer.WriteField("current", strconv.Itoa(i+1))
			writer.WriteField("filename", filename)
			contentType := writer.FormDataContentType()
			writer.Close()

			req.Header.Set("Content-Type", contentType)
			req.Body = ioutil.NopCloser(buf)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			wg.Done()
		}(i)
	}

	wg.Wait()
	return nil
}
