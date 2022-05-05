package dhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
)

type HandlerFunc func(w http.ResponseWriter, msg interface{}) // 回调方法

type HttpServer struct {
	handlers   *http.ServeMux
	respHeader *http.Header
	listenAddr string
}

func NewHttpServer(addr string) *HttpServer {
	s := new(HttpServer)
	s.handlers = http.NewServeMux()
	s.listenAddr = addr

	return s
}

func (s *HttpServer) SetResponseWriterHeader(header *http.Header) {
	s.respHeader = header
}

func setHeader(w *http.ResponseWriter, header *http.Header) {
	if header != nil {
		for key, values := range *header {
			for _, v := range values {
				(*w).Header().Add(key, v)
			}
		}
	}
}

// post
// 注册json请求,将请求数据通过json转成对应的结构
// 路由，结构，方法
func (s *HttpServer) HandleFuncJson(route string, elem interface{}, fn HandlerFunc) {
	if elem == nil {
		panic("elem is nil")
	}
	elemT := reflect.TypeOf(elem)

	s.handlers.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w, s.respHeader)
		msg := reflect.New(elemT.Elem()).Interface()
		err := json.NewDecoder(r.Body).Decode(&msg)
		defer r.Body.Close()

		if err != nil {
			log.Println("httpServer decode err:", err.Error())
			HttpServeError(w, 404, err.Error())
			return
		}

		fn(w, msg)
	})
}

// get
// 解析url的地址参数，如果参数不够
// 路由，参数，方法
func (s *HttpServer) HandleFuncUrlParam(route string, fn HandlerFunc) {
	s.handlers.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		setHeader(&w, s.respHeader)
		_ = r.ParseForm()
		fn(w, r.Form)
	})
}

// 注册http默认方法
func (s *HttpServer) Handle(route string, fn http.Handler) {
	s.handlers.Handle(route, fn)
}

func (s *HttpServer) HandleFunc(route string, fn http.HandlerFunc) {
	s.handlers.HandleFunc(route, func(w http.ResponseWriter, request *http.Request) {
		setHeader(&w, s.respHeader)
		fn(w, request)
	})
}

func (s *HttpServer) Listen() error {
	return http.ListenAndServe(s.listenAddr, s.handlers)
}

func HttpServeError(w http.ResponseWriter, status int, txt string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintln(w, txt)
}

// 文件上传
func HandleUpload(w http.ResponseWriter, r *http.Request) {
	//文件上传只允许POST方法
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method not allowed"))
		return
	}

	//从表单中读取文件
	file, fileHeader, err := r.FormFile("uploadFile")
	if err != nil {
		_, _ = io.WriteString(w, "Read file error")
		return
	}
	//defer 结束时关闭文件
	defer file.Close()
	//fmt.Println("filename: " + fileHeader.Filename)

	//创建文件
	newFile, err := os.Create("./" + fileHeader.Filename)
	if err != nil {
		_, _ = io.WriteString(w, "Create file error")
		return
	}
	//defer 结束时关闭文件
	defer newFile.Close()

	//将文件写到本地
	_, err = io.Copy(newFile, file)
	if err != nil {
		_, _ = io.WriteString(w, "Write file error")
		return
	}
	_, _ = io.WriteString(w, "Upload success")

}

// 文件下载
func handleDownload(w http.ResponseWriter, r *http.Request) {
	//文件上传只允许GET方法
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method not allowed"))
		return
	}
	//文件名
	filename := r.FormValue("filename")
	if filename == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
	log.Println("filename: " + filename)
	//打开文件
	file, err := os.Open("./" + filename)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
	//结束后关闭文件
	defer file.Close()

	//设置响应的header头
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
	//将文件写至responseBody
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
}
