package server

import (
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func Launch() {

	loadFilePath(config.FilePath)
	port := strings.Split(config.WebAddr, ":")[1]

	webAddr := fmt.Sprintf("0.0.0.0:%s", port)
	hServer := dhttp.NewHttpServer(webAddr)

	_, root := path.Split(config.FilePath)
	addrStr := fmt.Sprintf(`var httpAddr = "http://%s/file/";
var root = "%s";
var sliceSize = %d*1024*1024;`, config.WebAddr, root, config.SliceSize)
	err := ioutil.WriteFile(path.Join(config.WebIndex, "js/addr.js"), []byte(addrStr), os.ModePerm)
	if err != nil {
		panic(err)
	}

	//跨域
	header := http.Header{}
	header.Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	header.Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	header.Set("content-type", "application/json")             //返回数据格式是json
	hServer.SetResponseWriterHeader(&header)

	hServer.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(config.WebIndex))))

	hServer.HandleFuncUrlParam("/file/list", fileList)
	hServer.HandleFuncUrlParam("/file/delete", fileDelete)
	hServer.HandleFuncUrlParam("/file/mkdir", fileMkdir)
	hServer.HandleFuncJson("/file/check", &fileCheckReq{}, fileCheck)
	hServer.HandleFunc("/file/upload", fileUpload)
	hServer.HandleFuncUrlParam("/file/download", fileDownload)

	if err := hServer.Listen(); err != nil {
		logger.Errorf(err.Error())
	}
}
