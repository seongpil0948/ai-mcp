# TheShop AI

Go 기반의 AI 통합 서비스 API 서버입니다. Jira, GitLab, Confluence 등의 서비스와 통합하여 AI 기능을 제공합니다.

## 프로젝트 구조

```
theshop-ai/
├── api/                # API 문서 및 OpenAPI 명세
├── cmd/                # 애플리케이션 진입점
│   └── api/            # API 서버 메인 코드
├── configs/            # 설정 파일
├── internal/           # 내부 패키지 (외부에서 import 불가)
│   ├── app/            # 비즈니스 로직
│   ├── domain/         # 도메인 모델
│   └── transport/      # HTTP, gRPC 등 통신 계층
├── modules/            # 재사용 가능한 모듈
│   ├── config/         # 설정 관리
│   ├── core/           # 핵심 기능
│   └── integrations/   # 외부 서비스 통합
│       ├── aws/        # AWS 통합
│       ├── confluence/ # Confluence 통합
│       ├── gitlab/     # GitLab 통합
│       ├── jira/       # Jira 통합
│       └── llm/        # LLM 서비스 통합
├── pkg/                # 외부에서 import 가능한 패키지
│   ├── mcpclient/      # MCP 클라이언트
│   └── utils/          # 유틸리티 함수
├── scripts/            # 빌드 및 배포 스크립트
└── test/               # 통합 및 E2E 테스트
```

## 기능

- 워크스페이스 관리: 프로젝트 및 사용자별 작업 환경 제공
- Jira 통합: 이슈 생성, 조회, 검색
- GitLab 통합: 프로젝트 조회, MR 생성
- LLM 통합: OpenAI, Claude, Gemini 등의 LLM 서비스 사용

## 로컬 개발 환경 설정

### 사전 요구사항

- Go 1.24 이상
- Git
- Make

### 설치 및 설정

1. 프로젝트 클론:

```bash
git clone github.com/theshop/ai
cd theshop-ai
```

2. 개발 환경 설정:

```bash
# 개발 환경 설정 스크립트 실행
chmod +x ./scripts/setup-dev-env.sh
./scripts/setup-dev-env.sh

# 설정 파일 수정
vi ./configs/config.yaml  # 필요한 API 키 등 설정
```

3. 의존성 설치 및 빌드:

```bash
# Makefile 명령어로 의존성 설치 및 빌드
make all
```

## 모듈 관리

이 프로젝트는 Go 멀티 모듈 구조를 사용합니다. 다음 명령어로 모듈을 관리할 수 있습니다:

```bash
# 모듈 의존성 정리
make tidy

# 모듈 문제 수정
make fix-modules

# go.work 동기화
make sync
```

## 문제 해결 가이드

### Import 문제 해결

```bash
# Import 문제 해결 스크립트 실행
chmod +x ./scripts/fix-imports.sh
./scripts/fix-imports.sh
```

### 모듈 초기화 문제

```bash
# 모듈 캐시 정리 및 재설치
make clean
make setup
make tidy
make sync
```

### 빌드 오류

```bash
# Wire 코드 생성 문제
make wire

# 전체 프로젝트 다시 빌드
make clean all
```

## Makefile 사용법

```bash
# 기본 빌드 (권장)
make all

# 캐시 및 임시 파일 정리
make clean

# go.mod 파일 정리
make tidy

# go.work 동기화
make sync

# 프로젝트 빌드
make build

# 테스트 실행
make test

# Wire 코드 생성
make wire

# 모든 모듈의 go.mod 파일 수정
make fix-modules

# Go 의존성 재설치
make reinstall-deps

# 버전 정보 출력
make version

# 도움말 표시
make help
```

## 환경 변수

프로젝트는 다음 환경 변수를 지원합니다:

- `THESHOP_SERVER_PORT`: 서버 포트 (기본값: 8080)
- `THESHOP_SERVER_TIMEOUT`: 서버 타임아웃 (기본값: 30s)
- `THESHOP_JIRA_URL`: Jira URL
- `THESHOP_JIRA_USERNAME`: Jira 사용자명
- `THESHOP_JIRA_API_TOKEN`: Jira API 토큰
- `THESHOP_GITLAB_URL`: GitLab URL
- `THESHOP_GITLAB_TOKEN`: GitLab 토큰
- `THESHOP_LLM_PROVIDER`: 기본 LLM 제공자 (기본값: openai)
- `THESHOP_LLM_OPENAI_API_KEY`: OpenAI API 키
- `THESHOP_LLM_CLAUDE_API_KEY`: Claude API 키
- `THESHOP_LLM_GEMINI_API_KEY`: Gemini API 키

## 라이선스

이 프로젝트는 내부용 프로젝트로 모든 권리가 TheShop에 있습니다.
