package service

import (
	"io"
	"net/http"
	"os"
	"path"
)

type fileListData struct {
	Total int     `json:"total"`
	Items []*item `json:"items"`
}

type item struct {
	Filename string `json:"filename"`
	IsDir    bool   `json:"isDir"`
	Size     string `json:"size"`
	Date     string `json:"date"`
}

type fileHandler struct {
}

func (*fileHandler) list(wait *WaitConn, req struct {
	Path string `json:"path"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	info, err := findDir(req.Path, false)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	items := make([]*item, 0, len(info.FileInfos))
	for _, info := range info.FileInfos {
		// 正在上传中的文件不同步
		if info.IsDir || info.FileUpload == nil {
			_item := &item{
				Filename: info.Name,
				IsDir:    info.IsDir,
				Date:     info.ModeTime,
			}
			if info.IsDir {
				_item.Size = "-"
			} else {
				_item.Size = ConvertBytesString(uint64(info.FileSize))
			}

			items = append(items, _item)
		}
	}
	wait.SetResult("", &fileListData{
		Total: len(items),
		Items: items,
	})
}

func (*fileHandler) delete(wait *WaitConn, req struct {
	Path     string   `json:"path"`
	Filename []string `json:"filename"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	if req.Path == "" || len(req.Filename) == 0 {
		wait.SetResult("请求参数错误!", nil)
		return
	}

	info, err := findDir(req.Path, false)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	for _, filename := range req.Filename {
		if err = remove(info, filename); err != nil {
			wait.SetResult(err.Error(), nil)
			return
		}
	}
}

func (*fileHandler) rename(wait *WaitConn, req struct {
	Path    string `json:"path"`
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	if req.Path == "" || req.OldName == "" || req.NewName == "" {
		wait.SetResult("请求参数错误!", nil)
		return
	}

	parent, err := findDir(req.Path, false)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	info, ok := parent.FileInfos[req.OldName]
	if !ok {
		wait.SetResult("文件不存在", nil)
		return
	}

	info.Name = req.NewName
	info.Path = path.Join(parent.Path, parent.Name)
	info.AbsPath = path.Join(parent.AbsPath, info.Name)

	delete(parent.FileInfos, req.OldName)
	parent.FileInfos[req.NewName] = info

}

func (*fileHandler) copy(wait *WaitConn, req struct {
	SrcPath  string `json:"srcPath"`
	DestPath string `json:"destPath"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	if req.SrcPath == "" || req.DestPath == "" {
		wait.SetResult("请求参数错误!", nil)
		return
	}

}

func (*fileHandler) move(wait *WaitConn, req struct {
	SrcPath  string `json:"srcPath"`
	DestPath string `json:"destPath"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	if req.SrcPath == "" || req.DestPath == "" {
		wait.SetResult("请求参数错误!", nil)
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

	info, err := findDir(req.Path, false)
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
