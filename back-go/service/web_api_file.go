package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strings"
)

type fileListData struct {
	DiskUsed     uint64  `json:"diskUsed"`
	DiskUsedStr  string  `json:"diskUsedStr"`
	DiskTotal    uint64  `json:"diskTotal"`
	DiskTotalStr string  `json:"diskTotalStr"`
	Total        int     `json:"total"`
	Items        []*item `json:"items"`
}

type item struct {
	Filename string `json:"filename"`
	IsDir    bool   `json:"isDir"`
	Size     string `json:"size"`
	Date     string `json:"date"`
}

type fileHandler struct {
}

func (*fileHandler) mkdir(wait *WaitConn, req struct {
	Path string `json:"path"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)

	defer func() { wait.Done() }()

	if req.Path == "" {
		wait.SetResult("创建路径错误", nil)
		return
	}

	_, err := filePtr.FileInfo.findDir(req.Path, true)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

}

func (*fileHandler) list(wait *WaitConn, req struct {
	Path string `json:"path"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	info, err := filePtr.FileInfo.findDir(req.Path, false)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	items := make([]*item, 0, len(info.FileInfos))
	for _, info := range info.FileInfos {
		if info.IsDir || info.FileMD5 != "" {
			_item := &item{
				Filename: info.Name,
				IsDir:    info.IsDir,
				Date:     info.ModeTime,
			}
			if info.IsDir {
				_item.Size = "-"
			} else {
				_item.Size = ConvertBytesString(info.FileSize)
			}

			items = append(items, _item)
		}
	}
	wait.SetResult("", &fileListData{
		DiskTotal:    fileDiskTotal,
		DiskTotalStr: ConvertBytesString(fileDiskTotal),
		DiskUsed:     filePtr.UsedDisk,
		DiskUsedStr:  ConvertBytesString(filePtr.UsedDisk),
		Total:        len(items),
		Items:        items,
	})
}

func (*fileHandler) remove(wait *WaitConn, req struct {
	Path     string   `json:"path"`
	Filename []string `json:"filename"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	if req.Path == "" || len(req.Filename) == 0 {
		wait.SetResult("请求参数错误!", nil)
		return
	}

	info, err := filePtr.FileInfo.findDir(req.Path, false)
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

	calUsedDisk()
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

	if req.OldName == req.NewName {
		return
	}

	if strings.Contains(req.NewName, "/") {
		wait.SetResult("文件名不能含有'/'", nil)
		return
	}

	dirInfo, err := filePtr.FileInfo.findDir(req.Path, false)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	srcInfo, ok := dirInfo.FileInfos[req.OldName]
	if !ok {
		wait.SetResult("文件不存在", nil)
		return
	}

	if err = copy2(srcInfo, dirInfo, req.NewName); err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	// 移除原文件
	_ = remove(dirInfo, req.OldName)

	calUsedDisk()
}

// 移动、复制 文件或文件夹
func (*fileHandler) mvcp(wait *WaitConn, req struct {
	Source []string `json:"source"`
	Target string   `json:"target"`
	Move   bool     `json:"move"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	if len(req.Source) == 0 || req.Target == "" {
		wait.SetResult("请求参数错误!", nil)
		return
	}

	tarDir, err := filePtr.FileInfo.findDir(req.Target, false)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	for _, source := range req.Source {
		srcPath, srcName := path.Split(source)
		srcDir, err := filePtr.FileInfo.findDir(srcPath, false)
		if err != nil {
			wait.SetResult(err.Error(), nil)
			return
		}

		srcInfo, ok := srcDir.FileInfos[srcName]
		if !ok {
			wait.SetResult("文件不存在", nil)
			return
		}

		// 不能移动到自身或子目录下
		if tarDir.AbsPath == srcDir.AbsPath ||
			strings.Contains(tarDir.AbsPath, srcInfo.AbsPath) {
			wait.SetResult("不能拷贝、移动文件夹到自身目录或子目录", nil)
			return
		}

		if _, ok = tarDir.FileInfos[srcName]; ok {
			wait.SetResult("目标目录下已存在同名文件", nil)
			return
		}

		if err = copy2(srcInfo, tarDir, srcInfo.Name); err != nil {
			wait.SetResult(err.Error(), nil)
			return
		}

		if req.Move {
			// 移除原文件
			_ = remove(srcDir, srcName)
		}
	}

	calUsedDisk()

}

type downloadArg struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

func download(ctx *gin.Context, req *downloadArg) {
	logger.Infof("%s %v", getCurrentRoute(ctx), req)

	info, err := filePtr.FileInfo.findDir(req.Path, false)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	file, ok := info.FileInfos[req.Filename]
	if !ok || file.FileMD5 == "" {
		// 需要是一个存在的文件类型
		ctx.Status(http.StatusBadRequest)
		return
	}

	absPath := file.AbsPath
	if !saveFileMultiple {
		// 虚拟保存，修正到真实文件路径
		md5File_, ok := filePtr.MD5Files[file.FileMD5]
		if !ok {
			ctx.Status(http.StatusBadRequest)
			return
		}
		absPath = md5File_.File
	}

	//设置响应的header头
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", req.Filename))
	ctx.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.File(absPath)
	ctx.Status(200)
}
