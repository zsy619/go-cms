#!/bin/bash

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o cms.exe main.go
# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
