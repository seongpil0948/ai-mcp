.PHONY: all setup clean tidy sync build test wire fix-modules

all: clean setup fix-modules tidy sync build

setup:
	./scripts/setup.sh

clean:
	go clean -modcache
	find . -name "*.sum" -delete
	rm -f ./cmd/api/api

tidy:
	./scripts/tidy.sh

sync:
	go work sync

build:
	go build -v ./cmd/api

test:
	go test -v ./...

wire:
	cd ./modules/core && go run github.com/google/wire/cmd/wire

fix-modules:
	./scripts/fix-modules.sh

help:
	@echo "사용 가능한 명령어:"
	@echo "  all         : 전체 프로젝트 초기화 및 빌드"
	@echo "  setup       : 모듈 설정 초기화"
	@echo "  clean       : 캐시 및 임시 파일 정리"
	@echo "  tidy        : go.mod 파일 정리"
	@echo "  sync        : go.work 동기화"
	@echo "  build       : 프로젝트 빌드"
	@echo "  test        : 테스트 실행"
	@echo "  wire        : Wire 코드 생성"
	@echo "  fix-modules : 모든 모듈의 go.mod 파일 수정"
