package dhttp

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	urlpkg "net/url"
	"os"
	"path"
)

type Request struct {
	Client *http.Client
	req    *http.Request
	resp   *http.Response
	body   []byte // response body
}

func (rq *Request) Do() (resp *http.Response, err error) {
	if rq.resp != nil && rq.resp.StatusCode != 0 {
		return rq.resp, nil
	}

	if rq.Client == nil {
		rq.Client = &http.Client{}
	}

	resp, err = rq.Client.Do(rq.req)
	if err != nil {
		return nil, err
	}
	rq.resp = resp
	return
}

func (rq *Request) DoEnd() {
	resp, err := rq.Do()
	if err != nil {
		return
	}
	resp.Body.Close()
}

// PostFile add a post file to the request
func (rq *Request) PostFile(filename, filePath string) (*Request, error) {
	srcFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	fileWriter, err := writer.CreateFormFile("uploadFile", filename)
	if err != nil {
		return nil, err
	}

	//iocopy
	_, err = io.Copy(fileWriter, srcFile)
	if err != nil {
		return nil, err
	}
	contentType := writer.FormDataContentType()
	writer.Close()
	rq.req.Header.Set("Content-Type", contentType)
	rq.req.Body = ioutil.NopCloser(buf)

	return rq, nil
}

func (rq *Request) SetHeader(name, value string) *Request {
	rq.req.Header.Set(name, value)
	return rq
}

func (rq *Request) HttpRequest() *http.Request {
	return rq.req
}

func (rq *Request) HttpResponse() *http.Response {
	return rq.resp
}

func (rq *Request) WriteBody(data interface{}) *Request {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		rq.req.Body = ioutil.NopCloser(bf)
		rq.req.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		rq.req.Body = ioutil.NopCloser(bf)
		rq.req.ContentLength = int64(len(t))
	}
	return rq
}

// Param adds query param in to request.
func (rq *Request) WriteParam(values urlpkg.Values) *Request {
	rq.WriteBody(values.Encode())
	rq.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return rq
}

// XMLBody adds request raw body encoding by XML.
func (rq *Request) WriteXML(obj interface{}) (*Request, error) {
	if rq.req.Body == nil && obj != nil {
		byts, err := xml.Marshal(obj)
		if err != nil {
			return rq, err
		}
		rq.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		rq.req.ContentLength = int64(len(byts))
		rq.req.Header.Set("Content-Type", "application/xml")
	}
	return rq, nil
}

// JSONBody adds request raw body encoding by JSON.
func (rq *Request) WriteJSON(obj interface{}) (*Request, error) {
	if rq.req.Body == nil && obj != nil {
		byts, err := json.Marshal(obj)
		if err != nil {
			return rq, err
		}
		rq.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		rq.req.ContentLength = int64(len(byts))
		rq.req.Header.Set("Content-Type", "application/json")
	}
	return rq, nil
}

// String returns the body string in response.
// it calls Response inner.
func (rq *Request) ToString() (string, error) {
	data, err := rq.ToBytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (rq *Request) ToBytes() ([]byte, error) {
	if rq.body != nil {
		return rq.body, nil
	}
	resp, err := rq.Do()
	if err != nil {
		return nil, err
	}
	reader := resp.Body
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") == "gzip" && resp.Header.Get("Accept-Encoding") != "" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	}
	rq.body, err = ioutil.ReadAll(reader)
	return rq.body, err
}

// ToFile saves the body data in response to one file.
// it calls Response inner.
func (rq *Request) ToFile(filename string) error {
	resp, err := rq.Do()
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	err = pathExistAndMkdir(filename)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

//Check that the file directory exists, there is no automatically created
func pathExistAndMkdir(filename string) (err error) {
	filename = path.Dir(filename)
	_, err = os.Stat(filename)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(filename, os.ModePerm)
		if err == nil {
			return nil
		}
	}
	return err
}

// ToJSON returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (rq *Request) ToJSON(v interface{}) error {
	data, err := rq.ToBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// ToXML returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (rq *Request) ToXML(v interface{}) error {
	data, err := rq.ToBytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

func Get(url string) (*Request, error) {
	return NewRequest(url, "GET")
}

func GetBytes(url string) ([]byte, error) {
	req, err := Get(url)
	if err != nil {
		return nil, err
	}
	return req.ToBytes()
}

func PostJson(url string, obj interface{}) (*Request, error) {
	req, err := NewRequest(url, "POST")
	if err != nil {
		return nil, err
	}

	req, err = req.WriteJSON(obj)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// new
func NewRequest(url, method string) (rq *Request, err error) {
	rq = new(Request)
	rq.req, err = http.NewRequest(method, url, nil)
	return
}

// build URL params
// set params ?a=b&b=c
func BuildURLParams(url string, urlValue urlpkg.Values) string {
	return url + "?" + urlValue.Encode()
}
