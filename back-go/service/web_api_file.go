package service

import (
	"io"
	"net/http"
	"os"
)

type fileListData struct {
	Total int     `json:"total"`
	Items []*item `json:"items"`
}

type item struct {
	Filename string `json:"filename"`
	IsDir    bool   `json:"is_dir"`
	Size     int64  `json:"size"`
	Date     string `json:"date"`
}

type fileHandler struct {
}

func (*fileHandler) list(wait *WaitConn, req struct {
	Path string `json:"path"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	info, err := filePtr.findPath(req.Path, false)
	if err != nil {
		wait.SetResult("路径不存在", nil)
		return
	}

	items := make([]*item, 0, len(info.FileInfos))
	for _, info := range info.FileInfos {
		// 正在上传中的文件不同步
		if info.IsDir || info.FileOk {
			items = append(items, &item{
				Filename: info.Name,
				IsDir:    info.IsDir,
				Size:     info.FileSize,
				Date:     info.FileDate,
			})
		}
	}
	wait.SetResult("", &fileListData{
		Total: 0,
		Items: items,
	})
}

func (*fileHandler) delete(wait *WaitConn, req struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	if req.Path == "" || req.Filename == "" {
		wait.SetResult("请求参数错误!", nil)
		return
	}

	info, err := filePtr.findPath(req.Path, false)
	if err != nil {
		wait.SetResult("路径不存在", nil)
		return
	}

	if err = filePtr.remove(info, req.Filename); err != nil {
		wait.SetResult("文件不存在!", nil)
		return

	}
}

func (*fileHandler) download(wait *WaitConn, req struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	w := wait.Context().Writer

	info, err := filePtr.findPath(req.Path, false)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}

	file, ok := info.FileInfos[req.Filename]
	if !ok {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}

	absPath := file.AbsPath
	if !config.SaveFileMultiple {
		// 虚拟保存，修正到真实文件路径
		md5File_, ok := filePtr.MD5Files[file.FileMD5]
		if !ok {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, "Bad request")
			return
		}
		absPath = md5File_.File
	}

	//打开文件
	f, err := os.Open(absPath)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
		return
	}
	//结束后关闭文件
	defer f.Close()

	//设置响应的header头
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=\""+req.Filename+"\"")
	//将文件写至responseBody
	_, err = io.Copy(w, f)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Bad request")
	}
}
