package service

type uploadHandler struct {
}

func (*uploadHandler) check(wait *WaitConn, req struct {
	Path       string `json:"path"`
	Filename   string `json:"filename"`
	MD5        string `json:"md5"`
	Size       int    `json:"size"`
	SliceTotal int    `json:"sliceTotal"`
	SliceSize  int    `json:"sliceSize"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	if req.Path == "" || req.Filename == "" || req.MD5 == "" || req.Size == 0 || req.SliceTotal == 0 || req.SliceSize == 0 {
		wait.SetResult("请求参数错误", nil)
		return
	}

	if filePtr.UsedDisk+int64(req.Size) > fileDiskTotal {
		wait.SetResult("存储空间不足", nil)
		return
	}

	dirInfo, err := filePtr.FileInfo.findDir(req.Path, true)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	up := &upload{
		Md5:        req.MD5,
		Size:       req.Size,
		SliceSize:  req.SliceSize,
		Total:      req.SliceTotal,
		ExistSlice: map[string]int64{},
		Token:      nowFormat(),
	}

	resp := struct {
		Need       bool             `json:"need"` // 需要上传,不需要上传
		Token      string           `json:"token"`
		ExistSlice map[string]int64 `json:"existSlice"`
	}{}

	file, ok := dirInfo.FileInfos[req.Filename]
	if !ok {
		newInfo, err := dirInfo.makeChild(req.Filename, false)
		if err != nil {
			logger.Error(err)
			wait.SetResult(err.Error(), nil)
			return
		}

		files, md5Ok := filePtr.MD5Files[req.MD5]
		if md5Ok {
			// 已存在md5文件
			if saveFileMultiple {
				// 真实保存,拷贝文件
				if _, err := CopyFile(files.Ptr[0], newInfo.AbsPath); err != nil {
					logger.Error(err)
					wait.SetResult(err.Error(), nil)
					return
				}
			}
			newInfo.FileSize = int64(req.Size)
			newInfo.FileMD5 = req.MD5

			addMD5File(newInfo.FileMD5, newInfo)
			resp.Need = false

		} else {
			// 不存在md5文件，新建
			newInfo.FileUpload = up
			resp.Need = true
			resp.Token = up.Token
		}

		dirInfo.FileInfos[newInfo.Name] = newInfo
		wait.SetResult("", resp)
	} else {
		if file.IsDir {
			wait.SetResult("已存在同名文件夹", nil)
			return
		}

		if file.FileUpload != nil {
			// 文件正在上传

			if file.FileUpload.Md5 != req.MD5 {
				// 上传不同内容的同名文件，拒绝
				wait.SetResult("正在上传同名文件", nil)
				return
			} else {
				// 内容相同，上传了一部分
				// 存在两处同时上传，但分片大小不同，导致相互覆盖
				if file.FileUpload.SliceSize != req.SliceSize {
					// 分片大小不同，拒绝
					wait.SetResult("正在上传同名文件", nil)
					return
				}

				// 判断分片完整性
				for idx, size := range file.FileUpload.ExistSlice {
					if size != int64(req.SliceSize) {
						delete(file.FileUpload.ExistSlice, idx)
					}
				}

				resp.Need = true
				resp.Token = file.FileUpload.Token
				resp.ExistSlice = file.FileUpload.ExistSlice
				wait.SetResult("", resp)
			}

		} else {
			if file.FileMD5 == req.MD5 {
				// 已经上传完成,不需要上传
				resp.Need = false
			} else {
				file.FileUpload = up
				resp.Need = true
				resp.Token = up.Token
			}
			wait.SetResult("", resp)
		}

	}
}

func (*uploadHandler) upload(wait *WaitConn) {

	defer func() { wait.Done() }()

	ctx := wait.Context()
	filePath := ctx.PostForm("path")
	filename := ctx.PostForm("filename")
	token := ctx.PostForm("token")
	current := ctx.PostForm("current")

	logger.Infof("%s %s %s %s %s", wait.GetRoute(), filePath, filename, current, token)

	dirInfo, err := filePtr.FileInfo.findDir(filePath, false)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	file, ok := dirInfo.FileInfos[filename]
	if !ok || file.FileUpload == nil || file.FileUpload.Token != token {
		wait.SetResult("上传流程错误，check！", nil)
		return
	}

	_, ok = file.FileUpload.ExistSlice[current]
	if ok {
		// 当前分片已经上传
		return
	}

	gFile, err := ctx.FormFile("file")
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	src, err := gFile.Open()
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}
	defer src.Close()

	partFilename := makeFilePart(file.AbsPath, current)
	if _, err = WriteFile(partFilename, src); err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	file.FileUpload.ExistSlice[current] = gFile.Size
	file.mergeUpload()

}
