package service

import (
	"path"
	"strings"
	"time"
)

type fileShare struct {
	Key      string   `json:"key"` // 路由
	Path     string   `json:"path"`
	Filename []string `json:"filename"` // 分享的文件、文件夹。
	Token    string   `json:"token"`    // 访问密码
	Deadline int64    `json:"deadline"` // 到期时间戳 单位秒
}

// 服务器重启后，清空
var fileShared = map[string]*fileShare{}

type shareHandler struct {
}

func (*shareHandler) getSharedRoute(shared *fileShare) string {
	return "http://" + config.WebAddr + "/shared/s/" + shared.Key
}

func (this *shareHandler) create(wait *WaitConn, req struct {
	Path     string   `json:"path"`
	Filename []string `json:"filename"`
	Deadline int      `json:"deadline"` // 分享时间，单位天。0 永久
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)

	defer func() { wait.Done() }()

	if req.Path == "" || len(req.Filename) == 0 {
		wait.SetResult("请求参数错误!", nil)
		return
	}

	dirInfo, err := filePtr.FileInfo.findDir(req.Path, false)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	for _, filename := range req.Filename {
		_, ok := dirInfo.FileInfos[filename]
		if !ok {
			wait.SetResult("文件(夹)不存在", nil)
			return
		}
	}

	key := GenToken(12)
	for {
		if _, ok := fileShared[key]; !ok {
			break
		}
	}

	shared := &fileShare{
		Key:      key,
		Path:     req.Path,
		Filename: req.Filename,
		Token:    GenToken(4),
		Deadline: 0,
	}
	if req.Deadline > 0 {
		shared.Deadline = time.Now().Add(time.Hour * time.Duration(24*req.Deadline)).Unix()
	}

	fileShared[key] = shared

	wait.SetResult("", struct {
		Route    string `json:"route"`
		Token    string `json:"token"`
		Deadline int64  `json:"deadline"`
	}{
		Route:    this.getSharedRoute(shared),
		Token:    shared.Token,
		Deadline: shared.Deadline,
	})

	// 链接: https://pan.baidu.com/s/1ilIvYpJc2i6mkrfq3P_UNQ 提取码: 43cc

}

type sharedArg struct {
	Key   string `json:"key"`
	Token string `json:"token"`
	Path  string `json:"path"`
}

// 动态路由
func (*shareHandler) list(wait *WaitConn, req *sharedArg) {
	logger.Infof("%s %v", wait.GetRoute(), req)

	defer func() { wait.Done() }()

	shared, ok := fileShared[req.Key]
	if !ok || (shared.Deadline != 0 && time.Now().Unix() > shared.Deadline) {
		wait.SetResult("分享链接已过期", nil)
		delete(fileShared, req.Key)
		return
	}

	if shared.Token != req.Token {
		wait.SetResult("提取码错误", nil)
		return
	}

	if req.Path == shared.Path {
		// 根
		dirInfo, err := filePtr.FileInfo.findDir(req.Path, false)
		if err != nil {
			wait.SetResult(err.Error(), nil)
			return
		}

		items := make([]*item, 0, len(shared.Filename))
		for _, name := range shared.Filename {
			info, ok := dirInfo.FileInfos[name]
			if ok && (info.IsDir || info.FileSize != 0) {
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
			Total: len(items),
			Items: items,
		})

	} else {
		// 子目录
		children := false
		for _, name := range shared.Filename {
			if strings.Contains(req.Path, path.Join(shared.Path, name)) {
				children = true
				break
			}
		}

		if !children {
			wait.SetResult("路径不存在", nil)
			return
		}

		dirInfo, err := filePtr.FileInfo.findDir(req.Path, false)
		if err != nil {
			wait.SetResult(err.Error(), nil)
			return
		}

		items := make([]*item, 0, len(dirInfo.FileInfos))
		for _, info := range dirInfo.FileInfos {
			if info.IsDir || info.FileSize != 0 {
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

}
