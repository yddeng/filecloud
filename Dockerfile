FROM golang:1.15.1-alpine3.12
LABEL MAINTAINER="a nice guy"
WORKDIR $GOPATH/src/github.com/yddeng/filecloud
ADD . $GOPATH/src/github.com/yddeng/filecloud
RUN go env -w GO111MODULE=auto && go env -w GOPROXY=https://goproxy.io && go build  -o filecloud server/main/filecloud.go
ENTRYPOINT ["./filecloud", "config.toml"] 