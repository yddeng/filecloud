package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	filePtr *fileInfos
)

type fileInfos struct {
	mtx      sync.RWMutex        `json:"-"`
	FileInfo *fileInfo           `json:"file_info"`
	MD5File  map[string][]string `json:"md5_file"` // md5 -> absFilename ,存在相同的md5值直接拷贝
}

type fileInfo struct {
	Path      string               `json:"path"`           // 相对路径
	Name      string               `json:"name,omitempty"` // 名字
	AbsPath   string               `json:"abs_path"`       // 绝对路径
	IsDir     bool                 `json:"is_dir,omitempty"`
	FileOk    bool                 `json:"file_ok"`    // 当前文件是否已经写入
	FileSize  int64                `json:"file_size"`  // 文件有值
	FileMD5   string               `json:"file_md_5"`  // 文件有值
	FileDate  string               `json:"file_date"`  // 文件有值
	FileInfos map[string]*fileInfo `json:"file_infos"` // 文件夹有值
	Upload    *upload
}

type upload struct {
	Size     int64             `json:"size,omitempty"` // 文件有值
	MD5      string            `json:"md5,omitempty"`  // 文件有值
	SliceCnt int               `json:"slice_cnt"`      // 文件有值，文件上传时总文件数。 为0时，表示已传输完成。
	UpSlice  map[string]string `json:"up_slice"`       // 文件有值，文件上传时，已经上传的分片
}

func (this *fileInfo) clearUpload() {
	if this.Upload != nil {
		for part := range this.Upload.UpSlice {
			filename := makeFilePart(this.AbsPath, part)
			_ = os.RemoveAll(filename)
		}
	}
}

func (this *fileInfo) mergeUpload() {
	if this.Upload == nil {
		return
	}
	if this.Upload.SliceCnt != len(this.Upload.UpSlice) {
		return
	}
	f, err := os.Create(this.AbsPath)
	if err != nil {
		logger.Errorln(err)
		return
	}
	defer f.Close()

	for i := 1; i <= this.Upload.SliceCnt; i++ {
		partFile := makeFilePart(this.AbsPath, strconv.Itoa(i))
		pf, err := os.Open(partFile)
		if err != nil {
			logger.Errorln(err)
			return
		}
		written, err := io.Copy(f, pf)
		_ = pf.Close()
		if err != nil {
			logger.Errorln(err)
			return
		}
		logger.Infof("input %s from %s written %d ", this.AbsPath, partFile, written)
	}

	this.clearUpload()

	this.FileOk = true
	this.FileSize = this.Upload.Size
	this.FileMD5 = this.Upload.MD5
	this.FileDate = nowFormat()
	this.Upload = nil

	filePtr.addMD5(this.FileMD5, this.AbsPath)

}

func (this *fileInfos) addMD5(md5, file string) {
	files, ok := this.MD5File[md5]
	if !ok {
		files = []string{}
	}
	files = append(files, file)
	this.MD5File[md5] = files
}

func (this *fileInfos) removeMD5(md5, file string) {
	// 删除md5指向
	files, ok := this.MD5File[md5]
	if ok {
		idx := -1
		for i := 0; i < len(files); i++ {
			if files[i] == file {
				idx = i
				break
			}
		}
		if idx != -1 {
			files = append(files[:idx], files[idx+1:]...)
			if len(files) > 0 {
				this.MD5File[md5] = files
			} else {
				delete(this.MD5File, md5)
			}
		}
	}
}

func (this *fileInfos) remove(parent *fileInfo, name string) error {
	info, ok := parent.FileInfos[name]
	if !ok {
		return fmt.Errorf("%s 文件不存在", name)
	}

	// 遍历文件
	if err := walk(info, func(file *fileInfo) error {
		// 删除md5指向
		this.removeMD5(file.FileMD5, file.AbsPath)
		// 清理上传的分片
		file.clearUpload()

		return nil
	}); err != nil {
		return err
	}
	// 删除文件、文件夹
	if err := os.RemoveAll(info.AbsPath); err != nil {
		logger.Errorln(err)
	}
	// 删除info
	delete(parent.FileInfos, info.Name)

	return nil
}

func (this *fileInfos) findPath(filePath string, mkdir bool) (*fileInfo, error) {
	paths := splitPath(filePath)

	info := filePtr.FileInfo
	for i := 1; i < len(paths); i++ {
		dname := paths[i]
		cInfo, ok := info.FileInfos[dname]
		if ok {
			if !cInfo.IsDir {
				return nil, fmt.Errorf("已存在同名文件！")
			}
		} else {
			cInfo = &fileInfo{
				Path:      path.Join(info.Path, info.Name),
				Name:      dname,
				AbsPath:   path.Join(info.AbsPath, dname),
				IsDir:     true,
				FileInfos: map[string]*fileInfo{},
			}
			if err := os.MkdirAll(path.Join(cInfo.Path, cInfo.Name), os.ModePerm); err != nil {
				return nil, err
			}
			info.FileInfos[cInfo.Name] = cInfo
		}
		info = cInfo
	}
	return info, nil
}

// 遍历info，调用文件
func walk(info *fileInfo, f func(file *fileInfo) error) (err error) {
	if info == nil {
		return
	}
	if !info.IsDir {
		return f(info)
	}
	for _, cInfo := range info.FileInfos {
		if cInfo.IsDir {
			err = walk(cInfo, f)
		} else {
			err = f(cInfo)
		}
		if err != nil {
			return
		}
	}
	return
}

func loadFilePath(filePath string) {
	_ = os.MkdirAll(config.FilePath, os.ModePerm)
	sdir, dname := path.Split(filePath)
	filePtr = &fileInfos{
		mtx: sync.RWMutex{},
		FileInfo: &fileInfo{
			Path:      "",
			Name:      dname,
			AbsPath:   filePath,
			IsDir:     true,
			FileInfos: map[string]*fileInfo{},
		},
		MD5File: map[string][]string{},
	}

	Must(nil, filepath.Walk(filePath, func(absPath string, f os.FileInfo, err error) error {
		if err != nil {
			logger.Errorln(err)
			return err
		}

		relativePath := strings.TrimPrefix(absPath, sdir)
		if !f.IsDir() {
			// 是文件
			md5, e := fileMD5(absPath)
			if e != nil {
				logger.Errorln(e)
				return e
			}
			dir, file := path.Split(relativePath)
			info, _ := filePtr.findPath(dir, true)
			fInfo := &fileInfo{
				Path:     path.Join(info.Path, info.Name),
				Name:     file,
				AbsPath:  path.Join(info.AbsPath, file),
				IsDir:    false,
				FileSize: f.Size(),
				FileMD5:  md5,
				FileDate: f.ModTime().Format(timeFormat),
				FileOk:   true,
			}
			info.FileInfos[file] = fInfo
			filePtr.addMD5(md5, fInfo.AbsPath)
		} else {
			_, _ = filePtr.findPath(relativePath, true)
		}

		return nil
	}))

	str, _ := json.Marshal(filePtr)
	logger.Infoln(string(str))
}
