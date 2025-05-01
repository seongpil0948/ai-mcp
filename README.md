# theshop-ai MCP 서버

사내 데이터 연동 AI 어시스턴트를 위한 MCP(Model Context Protocol) 서버입니다. 이 서버는 GitLab, Notion, Figma 등 사내에서 사용 중인 주요 도구들과 연동하여 컨텍스트 기반의 정보 제공 및 자동 생성 기능을 수행합니다.

## 주요 기능

- **GitLab 연동**: 코드 컨텍스트 기반 추천, 변수명, 쿼리 구조 등 제안
- **Notion 연동**: Notion 페이지 조회 및 분석
- **Figma 연동**: Notion 콘텐츠 기반 Figma 디자인 자동 생성
- **LLM 통합**: OpenAI, Claude, Gemini 등 다양한 LLM 서비스 지원
- **MCP 프로토콜**: 표준 MCP 프로토콜을 통한 LLM 클라이언트 연동

## 시스템 요구사항

- Go 1.21 이상
- 외부 서비스 API 액세스 토큰 (GitLab, Notion, Figma, OpenAI 등)
- 운영체제: Linux, macOS, Windows

## 설치 및 설정

### 1. 저장소 클론

```bash
git clone https://github.com/seongpil0948/theshop-ai.git
cd theshop-ai
```

### 2. 초기화 스크립트 실행

```bash
chmod +x scripts/init.sh
./scripts/init.sh
```

이 스크립트는 필요한 디렉토리를 생성하고, 의존성을 설치하며, 기본 설정 파일과 환경 변수 파일을 생성합니다.

### 3. 설정 파일 수정

`configs/config.yaml` 파일과 `.env` 파일을 편집하여 서비스 API 키 및 서버 설정을 구성합니다.

```bash
# .env 파일 예시
THESHOP_GITLAB_TOKEN=your_gitlab_token
THESHOP_NOTION_TOKEN=your_notion_token
THESHOP_FIGMA_ACCESS_TOKEN=your_figma_access_token
THESHOP_LLM_OPENAI_API_KEY=your_openai_api_key
```

### 4. 서버 빌드

```bash
go build -o bin/mcpserver ./cmd/mcpserver
```

## 실행 방법

### 명령줄에서 직접 실행

```bash
# HTTP 모드로 실행 (기본 포트: 8080)
./bin/mcpserver --mode http

# 특정 포트로 실행
./bin/mcpserver --mode http --port 9000

# 디버그 모드 활성화
./bin/mcpserver --mode http --debug

# 다른 설정 파일 사용
./bin/mcpserver --config ./configs/custom-config.yaml
```

### 실행 스크립트 사용

```bash
# 기본 실행 (HTTP 모드, 8080 포트)
./scripts/run.sh

# 다른 모드나 포트로 실행
./scripts/run.sh --mode stdio
./scripts/run.sh --port 9000

# 디버그 모드 활성화
./scripts/run.sh --debug

# 다른 설정 파일 사용
./scripts/run.sh --config ./configs/custom-config.yaml
```

### Docker 실행

```bash
# Docker 이미지 빌드
docker build -t theshop-ai-mcp-server .

# Docker 컨테이너 실행
docker run -p 8080:8080 \
  -e THESHOP_GITLAB_TOKEN=your_gitlab_token \
  -e THESHOP_NOTION_TOKEN=your_notion_token \
  -e THESHOP_FIGMA_ACCESS_TOKEN=your_figma_access_token \
  -e THESHOP_LLM_OPENAI_API_KEY=your_openai_api_key \
  theshop-ai-mcp-server
```

### Docker Compose 실행

```bash
# 환경 변수 설정 (.env 파일)
cp .env.example .env
# .env 파일 편집 후

# Docker Compose로 실행
docker-compose up -d
```

## 사용 방법

### 1. HTTP MCP 서버 엔드포인트

서버가 HTTP 모드로 실행 중인 경우, 다음 엔드포인트를 통해 접근할 수 있습니다:

- MCP 엔드포인트: `http://localhost:8080/mcp`
- 상태 확인: `http://localhost:8080/health`

### 2. 지원되는 MCP 도구

#### GitLab 도구

- `gitlab_list_projects`: GitLab 프로젝트 목록 조회
- `gitlab_get_project`: 프로젝트 정보 조회
- `gitlab_get_file`: 파일 내용 조회
- `gitlab_list_merge_requests`: 머지 리퀘스트 목록 조회
- `gitlab_create_merge_request`: 머지 리퀘스트 생성
- `gitlab_create_file`: 파일 생성
- `gitlab_analyze_code`: 코드 분석 및 추천

#### Notion 도구

- `notion_search`: Notion 검색
- `notion_get_page`: 페이지 정보 조회
- `notion_get_database`: 데이터베이스 정보 조회
- `notion_query_database`: 데이터베이스 쿼리
- `notion_create_page`: 페이지 생성
- `notion_get_block_children`: 블록 자식 요소 조회
- `notion_append_block_children`: 블록에 자식 요소 추가

#### Figma 도구

- `figma_get_file`: Figma 파일 정보 조회
- `figma_get_nodes`: 노드 정보 조회
- `figma_create_file`: 파일 생성
- `figma_get_team_projects`: 팀 프로젝트 목록 조회
- `figma_get_project_files`: 프로젝트 파일 목록 조회
- `figma_create_from_notion`: Notion 페이지에서 Figma 파일 생성

#### LLM 도구

- `llm_generate_text`: 텍스트 생성
- `llm_analyze_code`: 코드 분석
- `llm_summarize`: 내용 요약
- `llm_suggest_variable_names`: 변수명 추천
- `llm_generate_sql`: SQL 쿼리 생성

## 통합 클라이언트

이 서버는 다음과 같은 MCP 클라이언트와 통합할 수 있습니다:

- Claude Client
- OpenAI 호환 클라이언트
- 기타 MCP 프로토콜 지원 클라이언트

## 프로젝트 구조

```
theshop-ai/
├── cmd/                   # 애플리케이션 진입점
│   └── mcpserver/         # MCP 서버 메인 코드
├── configs/               # 설정 파일
├── internal/              # 내부 패키지
│   ├── app/               # 애플리케이션 로직
│   ├── domain/            # 도메인 모델
│   └── transport/         # 통신 계층
├── modules/               # 기능 모듈
│   ├── config/            # 설정 관리
│   └── integrations/      # 외부 서비스 통합
│       ├── converter/     # 데이터 변환
│       ├── figma/         # Figma API 통합
│       ├── gitlab/        # GitLab API 통합
│       ├── llm/           # LLM 서비스 통합
│       └── notion/        # Notion API 통합
├── pkg/                   # 공유 패키지
│   └── mcpserver/         # MCP 서버 핵심 로직
├── scripts/               # 유틸리티 스크립트
├── .env                   # 환경 변수 파일
├── docker-compose.yml     # Docker Compose 설정
└── Dockerfile             # Docker 이미지 빌드 설정
```

## 개발

### 1. 코드 작성 규칙

- Go 표준 코딩 스타일 및 관례 준수
- 모든 API 엔드포인트에 대한 적절한 로깅 및 오류 처리 포함
- 새로운 기능은 해당 모듈에 적절히 분류하여 추가

### 2. 테스트

```bash
# 전체 테스트 실행
go test ./...

# 특정 패키지 테스트 실행
go test ./pkg/mcpserver
```

### 3. 기여

1. 변경 사항에 대한 이슈 생성
2. 브랜치 생성 (`feature/your-feature` 또는 `fix/your-fix`)
3. 변경 사항 커밋
4. 테스트 실행 및 코드 품질 확인
5. Pull Request 제출

## 문제 해결

- **서버 시작 실패**: 환경 변수 및 설정 파일 확인
- **API 연결 오류**: 각 서비스의 API 키 및 액세스 토큰 확인
- **로그 확인**: `logs/server.log` 파일에서 자세한 로그 확인

## 라이센스

이 프로젝트는 내부용으로 개발되었으며, 모든 권리는 TheShop에 있습니다.