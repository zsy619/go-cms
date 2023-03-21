#!/bin/bash

# 减小 golang 编译出程序的体积
# http://blog.fatedier.com/2017/02/04/reduce-golang-program-size/
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o cms main.go
# GOOS=linux GOARCH=amd64 go build -o cms main.go
 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o cms main.go

# nohup后台执行程序
# nohup /home/api/cms  >> /home/api/output.log 2>&1 &
# nohup ./cms  >> output.log 2>&1 &``