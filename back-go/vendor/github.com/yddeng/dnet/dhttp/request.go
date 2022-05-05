package dhttp

/*
import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	urlpkg "net/url"
	)

type Request struct {
	*http.Request
}

func NewEnd(req *http.Request) *Request {
	return &Request{req}
}

func (rq *Request) WriteBody(data interface{}) *Request {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		rq.Body = ioutil.NopCloser(bf)
		rq.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		rq.Body = ioutil.NopCloser(bf)
		rq.ContentLength = int64(len(t))
	}
	return rq
}

// Param adds query param in to request.
func (rq *Request) WriteParam(values urlpkg.Values) *Request {
	rq.WriteBody(values.Encode())
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return rq
}

// XMLBody adds request raw body encoding by XML.
func (rq *Request) WriteXML(obj interface{}) (*Request, error) {
	if rq.Body == nil && obj != nil {
		byts, err := xml.Marshal(obj)
		if err != nil {
			return rq, err
		}
		rq.Body = ioutil.NopCloser(bytes.NewReader(byts))
		rq.ContentLength = int64(len(byts))
		rq.Header.Set("Content-Type", "application/xml")
	}
	return rq, nil
}

// JSONBody adds request raw body encoding by JSON.
func (rq *Request) WriteJSON(obj interface{}) (*Request, error) {
	if rq.Body == nil && obj != nil {
		byts, err := json.Marshal(obj)
		if err != nil {
			return rq, err
		}
		rq.Body = ioutil.NopCloser(bytes.NewReader(byts))
		rq.ContentLength = int64(len(byts))
		rq.Header.Set("Content-Type", "application/json")
	}
	return rq, nil
}

func s()  {
	q := NewRequest(nil)
	q.
}

*/
