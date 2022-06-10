# 构建：使用golang:1.17.7版本
FROM golang:1.17.7 AS builder
# 设置工作区
WORKDIR /go/src/github.com/yddeng/filecloud
# 把全部文件添加到/go/src/github.com/yddeng/filecloud目录
COPY . .
# 编译：把back-go/cmd/filecloud.go编译成可执行的二进制文件，filecloud
RUN go env -w GO111MODULE=off && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o filecloud back-go/cmd/filecloud.go


# step 2
# 运行：使用scratch作为基础镜像 alpine
FROM scratch
# 设置工作区
WORKDIR ./back-go/cmd
# 在builder阶段复制时区
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 在builder阶段复制可执行的go二进制文件app
COPY --from=builder /go/src/github.com/yddeng/filecloud/filecloud .
# 在builder阶段复制可执行的go二进制程序的配置文件
COPY --from=builder /go/src/github.com/yddeng/filecloud/back-go/cmd/config.toml .
# 在builder阶段复制前端项目
COPY --from=builder /go/src/github.com/yddeng/filecloud/front-vue/dist /front-vue/dist/

# 启动服务
CMD ["./filecloud"]
