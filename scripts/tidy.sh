#!/bin/bash
set -e


find . -name "go.mod" -type f -not -path "*/vendor/*" | while read -r file; do
    dir=$(dirname "$file")
    if [ "$dir" != "." ]; then
        echo "정리 중: $dir"
        (cd "$dir" && go mod tidy) || echo "경고: $dir 모듈 정리 실패"
    fi
done

echo "===== 루트 go.mod 정리 ====="
go mod tidy || echo "경고: 루트 모듈 최종 정리 실패"

echo "===== Wire 코드 생성 (core) ====="
(cd modules/core && go generate ./...) || echo "경고: go generate ./... (core) 실패"

echo "===== 완료 ====="
