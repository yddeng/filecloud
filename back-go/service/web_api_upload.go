package service

import (
	"path"
)

type uploadHandler struct {
}

func (*uploadHandler) check(wait *WaitConn, req struct {
	Path     string            `json:"path"`
	Filename string            `json:"filename"`
	MD5      string            `json:"md5"`
	Size     int               `json:"size"`
	SliceMd5 map[string]string `json:"sliceMd5"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)
	defer func() { wait.Done() }()

	if req.Path == "" || req.Filename == "" || req.MD5 == "" || req.Size == 0 || len(req.SliceMd5) == 0 {
		wait.SetResult("请求参数错误", nil)
		return
	}

	info, err := filePtr.findPath(req.Path, true)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	up := &upload{
		Size:     int64(req.Size),
		MD5:      req.MD5,
		SliceCnt: len(req.SliceMd5),
		UpSlice:  map[string]string{},
	}

	resp := struct {
		Need       bool              `json:"need"` // 需要上传,不需要上传
		ExistSlice map[string]string `json:"existSlice"`
	}{}

	absPath := path.Join(info.AbsPath, req.Filename)
	file, ok := info.FileInfos[req.Filename]
	if !ok {
		newInfo := &fileInfo{
			Path:    path.Join(info.Path, info.Name),
			Name:    req.Filename,
			AbsPath: absPath,
			FileMD5: req.MD5,
		}

		files, md5Ok := filePtr.MD5Files[req.MD5]
		if md5Ok {
			// 已存在md5文件
			if config.SaveFileMultiple {
				// 真实保存,拷贝文件
				if _, err := CopyFile(files.Ptr[0], absPath); err != nil {
					logger.Error(err)
					wait.SetResult(err.Error(), nil)
					return
				}
			}
			newInfo.FileOk = true
			newInfo.FileSize = files.Size
			newInfo.FileDate = nowFormat()
			info.FileInfos[newInfo.Name] = newInfo
			filePtr.addMD5File(newInfo.FileMD5, newInfo)

			resp.Need = false
			wait.SetResult("", resp)

		} else {
			// 不存在md5文件，新建
			newInfo.Upload = up
			info.FileInfos[newInfo.Name] = newInfo

			resp.Need = true
			wait.SetResult("", resp)
		}
	} else {
		if file.IsDir {
			wait.SetResult("已存在同名文件夹", nil)
			return
		}
		if file.FileMD5 != req.MD5 {
			// 原文件已经改变，需要上传
			if file.Upload == nil {
				file.Upload = up

				resp.Need = true
				wait.SetResult("", resp)
			} else {
				if file.Upload.MD5 == req.MD5 {
					// 新文件已经上传了一部分
					//file.mergeUpload()

					for i, sliceMd5 := range req.SliceMd5 {
						if existSliceMd5, ok := file.Upload.UpSlice[i]; ok && existSliceMd5 != sliceMd5 {
							// 不一致
							delete(file.Upload.UpSlice, i)
						}
					}

					resp.Need = true
					resp.ExistSlice = file.Upload.UpSlice
					wait.SetResult("", resp)
				} else {
					// 新文件没有上传完，但又上传不同md5文件
					file.clearUpload()

					resp.Need = true
					wait.SetResult("", resp)
				}
			}
		} else {
			if file.Upload == nil {
				// 已经上传完成,不需要上传
				resp.Need = false
				wait.SetResult("", resp)
			} else {
				if file.Upload.MD5 == req.MD5 {
					// 新文件已经上传了一部分
					for i, sliceMd5 := range req.SliceMd5 {
						if existSliceMd5, ok := file.Upload.UpSlice[i]; ok && existSliceMd5 != sliceMd5 {
							// 不一致
							delete(file.Upload.UpSlice, i)
						}
					}
					resp.Need = true
					resp.ExistSlice = file.Upload.UpSlice
					wait.SetResult("", resp)
				} else {
					// 新文件没有上传完，但又上传不同md5文件
					file.clearUpload()
					file.Upload = up

					resp.Need = true
					wait.SetResult("", resp)
				}
			}
		}
	}

}

func (*uploadHandler) upload(wait *WaitConn) {

	defer func() { wait.Done() }()

	ctx := wait.Context()
	filePath := ctx.PostForm("path")
	filename := ctx.PostForm("filename")
	md5 := ctx.PostForm("md5")
	current := ctx.PostForm("current")
	sliceMd5 := ctx.PostForm("sliceMd5")

	logger.Infof("%s %s %s %s %s %s", wait.GetRoute(), filePath, filename, md5, current, sliceMd5)

	info, err := filePtr.findPath(filePath, false)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

	file, ok := info.FileInfos[filename]
	if !ok || file.Upload == nil || file.Upload.MD5 != md5 {
		wait.SetResult("上传流程错误，check！", nil)
		return
	}

	_, ok = file.Upload.UpSlice[current]
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

	file.Upload.UpSlice[current] = sliceMd5
	file.mergeUpload()

}
