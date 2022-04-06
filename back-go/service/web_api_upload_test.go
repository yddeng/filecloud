package service

import (
	"bytes"
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"sync"
	"testing"
)

func TestUpload(t *testing.T) {
	filePath := "./config.go"

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

	_, filename := path.Split(filePath)

	elem := map[string]interface{}{
		"path":     uploadPath,
		"filename": filename,
		"md5":      md5,
		"total":    cnt,
		"size":     len(data),
	}

	req, _ := dhttp.PostJson("http://127.0.0.1:9987/upload/check", elem)
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

	_, filename := path.Split(filePath)

	wg := sync.WaitGroup{}
	wg.Add(cnt)

	for i := 0; i < cnt; i++ {
		go func(i int) {
			req, _ := http.NewRequest("POST", "http://127.0.0.1:9987/upload/upload", nil)

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
