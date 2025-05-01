#!/bin/bash

# 루트 디렉토리로 이동
cd /Users/2309-n0015/Code/Project/Util/theshop-ai

# Go 모듈 캐시 정리
go clean -modcache

# mcpclient 패키지에 필요한 의존성 추가
cd pkg/mcpclient
go get github.com/mark3labs/mcp-go@latest

# 루트 디렉토리로 돌아가서 go.mod 파일 정리
cd /Users/2309-n0015/Code/Project/Util/theshop-ai
go mod tidy

# wire_gen.go 파일 생성 (선택 사항)
go install github.com/google/wire/cmd/wire@latest
go run github.com/google/wire/cmd/wire ./modules/core

# 모든 모듈의 go.mod 파일 정리
cd internal
go mod tidy

cd ../modules/core
go mod tidy

cd ../config
go mod tidy

cd ../integrations/jira
go mod tidy

cd ../gitlab
go mod tidy

cd ../llm
go mod tidy

cd ../../../pkg/mcpclient
go mod tidy

# 루트 디렉토리로 돌아가기
cd /Users/2309-n0015/Code/Project/Util/theshop-ai

# 워크스페이스 동기화
go work sync

# 최종 빌드 테스트
go build -v ./cmd/api
