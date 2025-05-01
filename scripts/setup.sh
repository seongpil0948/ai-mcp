#!/bin/bash

echo "===== 모듈 설정 초기화 ====="
go install github.com/google/wire/cmd/wire@latest
go get -u github.com/mark3labs/mcp-go@latest
