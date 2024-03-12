#!/bin/bash
version="0.0.1"
sourceVersion=$(git log --date=iso --pretty=format:"%H @%cd" -1)
buildTime="$(date '+%Y-%m-%d %H:%M:%S') by $(go version)"

cat <<EOF | gofmt >utils/version.go
package utils

var (
	Version = "$version"
	SourceVersion = "$sourceVersion"
	BuildTime = "$buildTime"
)
EOF

go fmt ./...
go build -ldflags "-w -s" -o vmmanager main.go

