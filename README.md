# fileCloud

仿百度网盘，实现个人文件存储云盘。


## 页面展示

![文件上传](res/file_upload.jpg)

## 分片上传、断点续传

文件过大时，上传文件需要很长时间，且中途退出将导致文件重传。

分片上传: 上传文件时，在本地将文件按照 2M 的大小将文件进行分片。在服务器端将文件组合。

断点续传: 如果文件没有上传完，关闭客户端。再一次上传文件时，对比服务器已经上传的分片，只需要上传没有的分片。

## 秒传

每一个文件都有对应的md5码。当检测上传文件时，如果本地已经存在相同md5的文件，直接将文件拷贝到对应目录。

## 启动

`go get github.com/yddeng/filecloud`

1. config.toml

```
WebAddr  = "127.0.0.1:9987"             # web 地址
WebIndex = "./app"                      # web html目录
FilePath = "cloud"                      # 文件存放目录
```

2. `go run server/main/filecloud.go config.toml` 

3. 浏览器访问 webAddr。