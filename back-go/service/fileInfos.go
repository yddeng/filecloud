package service

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	filePtr          *fileInfos
	saveFileMultiple = true
)

type fileInfos struct {
	FileInfo *fileInfo           `json:"fileInfo"`
	MD5Files map[string]*md5File `json:"md5Files"`
}

type md5File struct {
	File string   `json:"file"` // 原始文件
	Size int64    `json:"size"`
	MD5  string   `json:"md5"`
	Ptr  []string `json:"ptr"` // 文件引用
}

type fileInfo struct {
	Path       string               `json:"path"`           // 相对路径
	Name       string               `json:"name,omitempty"` // 名字
	AbsPath    string               `json:"absPath"`        // 绝对路径
	IsDir      bool                 `json:"isDir,omitempty"`
	ModeTime   string               `json:"modeTime"`
	FileSize   int64                `json:"fileSize"`  // 文件有值。有值表示存在文件，无值说明正在上传
	FileMD5    string               `json:"fileMd5"`   // 文件有值。有值表示存在文件，无值说明正在上传
	FileInfos  map[string]*fileInfo `json:"fileInfos"` // 文件夹有值
	FileUpload *upload              `json:"_"`         // 文件上传零时数据
}

type upload struct {
	Md5        string           `json:"md5"`        // 文件上传时的md5值
	Size       int              `json:"size"`       // 文件总大小
	SliceSize  int              `json:"sliceSize"`  // 上传的分片大小
	Total      int              `json:"total"`      // 文件上传时分片总数
	ExistSlice map[string]int64 `json:"existSlice"` // 文件上传时，已经上传的分片
	Token      string           `json:"token"`      // 上传时需要验证token
}

func (this *fileInfo) makeChild(name string, isDir bool) (*fileInfo, error) {
	info := &fileInfo{
		Path:     path.Join(this.Path, this.Name),
		Name:     name,
		AbsPath:  path.Join(this.AbsPath, name),
		IsDir:    isDir,
		ModeTime: nowFormat(),
	}
	if isDir {
		info.FileInfos = map[string]*fileInfo{}
		if err := os.MkdirAll(info.AbsPath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	return info, nil
}

// 查找目录
func (this *fileInfo) findDir(filePath string, mkdir bool) (*fileInfo, error) {
	paths := strings.Split(path.Clean(filePath), "/")

	info := this
	for i := 1; i < len(paths); i++ {
		dirName := paths[i]
		cInfo, ok := info.FileInfos[dirName]
		if ok {
			if !cInfo.IsDir {
				return nil, fmt.Errorf("已存在同名文件")
			}
		} else {
			if mkdir {
				var err error
				if cInfo, err = info.makeChild(dirName, true); err != nil {
					return nil, err
				}
				info.FileInfos[cInfo.Name] = cInfo
			} else {
				return nil, fmt.Errorf("路径不存在")
			}
		}
		info = cInfo
	}
	return info, nil
}

func (this *fileInfo) clearUpload() {
	if this.FileUpload != nil {
		for part := range this.FileUpload.ExistSlice {
			filename := makeFilePart(this.AbsPath, part)
			_ = os.RemoveAll(filename)
		}
	}
}

func (this *fileInfo) mergeUpload() {
	if this.FileUpload == nil || this.FileUpload.Total != len(this.FileUpload.ExistSlice) {
		return
	}
	f, err := os.Create(this.AbsPath)
	if err != nil {
		return
	}
	defer f.Close()

	for i := 0; i < this.FileUpload.Total; i++ {
		partFile := makeFilePart(this.AbsPath, strconv.Itoa(i))
		pf, err := os.Open(partFile)
		if err != nil {
			return
		}
		written, err := io.Copy(f, pf)
		_ = pf.Close()
		if err != nil {
			return
		}
		logger.Infof("input %s from %s written %d ", this.AbsPath, partFile, written)
	}

	this.clearUpload()
	if this.FileMD5 != "" {
		// 移除原文件
		removeMD5File(this.FileMD5, this.AbsPath)
	}

	this.FileMD5 = this.FileUpload.Md5
	this.FileSize = int64(this.FileUpload.Size)
	this.ModeTime = nowFormat()
	this.FileUpload = nil

	addMD5File(this.FileMD5, this)
}

func addMD5File(md5 string, info *fileInfo) {
	files, ok := filePtr.MD5Files[md5]
	if !ok {
		files = &md5File{
			File: info.AbsPath,
			MD5:  info.FileMD5,
			Size: info.FileSize,
			Ptr:  []string{},
		}
		filePtr.MD5Files[md5] = files
	}
	files.Ptr = append(files.Ptr, info.AbsPath)
}

func removeMD5File(md5, ptr string) {
	// 删除md5指向
	files, ok := filePtr.MD5Files[md5]
	if ok {
		idx := -1
		for i := 0; i < len(files.Ptr); i++ {
			if files.Ptr[i] == ptr {
				idx = i
				break
			}
		}
		if idx != -1 {
			files.Ptr = append(files.Ptr[:idx], files.Ptr[idx+1:]...)
			if len(files.Ptr) == 0 {
				delete(filePtr.MD5Files, md5)
			}
		}
	}
}

// 文件删除，
func remove(parent *fileInfo, name string) error {
	info, ok := parent.FileInfos[name]
	if !ok {
		return fmt.Errorf("%s 文件不存在", name)
	}

	delMd5 := map[string]struct{}{} // 待删除的md5文件，源文件

	// 遍历文件
	if err := walk(info, func(file *fileInfo) error {
		if !file.IsDir && file.FileMD5 != "" {
			if !saveFileMultiple {
				if md5File_, ok := filePtr.MD5Files[file.FileMD5]; ok {
					if md5File_.File == file.AbsPath {
						// 此文件为源文件
						delMd5[file.FileMD5] = struct{}{}
					}
				}
			}
			// 删除md5指向
			removeMD5File(file.FileMD5, file.AbsPath)
			// 清理上传的分片
			file.clearUpload()
		}

		return nil
	}); err != nil {
		return err
	}

	// 删除info
	delete(parent.FileInfos, info.Name)

	if !saveFileMultiple {
		// 文件夹中包含源文件需要拷贝到他处
		for md5 := range delMd5 {
			md5File_, ok := filePtr.MD5Files[md5]
			if ok {
				// 还存在他处引用
				_ = os.Rename(md5File_.File, md5File_.Ptr[0])
				md5File_.File = md5File_.Ptr[0]
			}
		}
	}

	// 删除文件、文件夹
	if err := os.RemoveAll(info.AbsPath); err != nil {
		return err
	}

	return nil
}

// 拷贝到目标目录下
func copy2(src, destParent *fileInfo, destName string) error {
	srcPath := path.Join(src.Path, src.Name)
	return walk(src, func(file *fileInfo) error {
		var fileName string
		var dirInfo, newInfo *fileInfo
		var err error
		if file.AbsPath == src.AbsPath {
			// 自己
			fileName = destName
			dirInfo = destParent
		} else {
			// 当前分支 目录拷贝
			filePath := path.Join(file.Path, file.Name)
			revPath := strings.TrimPrefix(filePath, srcPath+"/")
			revPath = path.Dir(revPath)

			if destParent.Path == "" {
				// 根目录
				revPath = path.Join("cloud", destName, revPath)
			} else {
				revPath = path.Join(destParent.Path, destName, revPath)
			}

			fileName = file.Name

			//fmt.Println(22222, filePath, srcPath, revPath, fileName)
			if dirInfo, err = destParent.findDir(revPath, true); err != nil {
				return err
			}

		}

		if newInfo, err = dirInfo.makeChild(fileName, file.IsDir); err != nil {
			return err
		}
		//fmt.Println(srcPath, file.Path, file.Name, newInfo.Path, newInfo.Name)
		if !file.IsDir && file.FileMD5 != "" {
			if saveFileMultiple {
				// 真实保存,拷贝文件
				files, _ := filePtr.MD5Files[file.FileMD5]
				//logger.Info(files, newInfo)
				if _, err := CopyFile(files.Ptr[0], newInfo.AbsPath); err != nil {
					return err
				}
			}

			newInfo.FileSize = file.FileSize
			newInfo.FileMD5 = file.FileMD5
			addMD5File(newInfo.FileMD5, newInfo)
		}
		dirInfo.FileInfos[newInfo.Name] = newInfo

		return nil
	})
}

// 遍历info 包含自己
func walk(info *fileInfo, f func(file *fileInfo) error) (err error) {
	if info == nil {
		return
	}
	if err = f(info); err != nil {
		return
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
	filePath = path.Clean(filePath)
	_ = os.MkdirAll(filePath, os.ModePerm)
	dirPrefix, _ := path.Split(filePath)
	filePtr = &fileInfos{
		FileInfo: &fileInfo{
			Path:      "",
			Name:      "cloud",
			AbsPath:   filePath,
			IsDir:     true,
			FileInfos: map[string]*fileInfo{},
		},
		MD5Files: map[string]*md5File{},
	}

	Must(nil, filepath.Walk(filePath, func(absPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relativePath := strings.TrimPrefix(absPath, dirPrefix)
		fmt.Println(dirPrefix, relativePath, absPath, f.Name(), f.ModTime(), f.IsDir(), f.Size())
		if !f.IsDir() {
			// 是文件

			_, filename := path.Split(absPath)
			//fmt.Println(absPath, f.Name(), filename)
			if strings.Contains(filename, ".part") {
				// 是上传时的文件分片，删除
				_ = os.RemoveAll(absPath)
			} else {
				md5, e := fileMD5(absPath)
				if e != nil {
					return e
				}
				dir, file := path.Split(relativePath)
				dirInfo, _ := filePtr.FileInfo.findDir(dir, true)
				if fileInfo, err := dirInfo.makeChild(file, false); err != nil {
					return err
				} else {
					fileInfo.FileSize = f.Size()
					fileInfo.FileMD5 = md5
					fileInfo.ModeTime = f.ModTime().Format(timeFormat)
					dirInfo.FileInfos[file] = fileInfo
					addMD5File(md5, fileInfo)
				}
			}
		} else {
			dirInfo, _ := filePtr.FileInfo.findDir(relativePath, true)
			dirInfo.ModeTime = f.ModTime().Format(timeFormat)
		}

		return nil
	}))

}
