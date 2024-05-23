#!/bin/sh

# 减小 golang 编译出程序的体积
# http://blog.fatedier.com/2017/02/04/reduce-golang-program-size/

rm -r "releases"
mkdir "releases"

# 【darwin/amd64】
echo "start build darwin/amd64 ..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./releases/cms main.go

# 【linux/amd64】
echo "start build linux/amd64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./releases/cms-linux-amd64 main.go

# 【windows/amd64】
echo "start build windows/amd64 ..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./releases/cms-windows-amd64.exe main.go

echo "Congratulations,all build success!!!"
