#!/bin/bash
#编译全平台
mkdir -p output
GOOS=windows go build -ldflags "-s -w" -o output/windows_ut.exe
GOOS=linux go build -ldflags "-s -w" -o output/linux_ut
GOOS=darwin go build -ldflags "-s -w" -o output/darwin_ut