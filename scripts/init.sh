#!/bin/bash
# Path: scripts/init.sh

echo "===== theshop-ai-mcp-server 초기화 ====="

# 환경 설정
echo "환경 설정 중..."
mkdir -p bin
mkdir -p logs

# Go 모듈 초기화
echo "Go 모듈 초기화 중..."
go mod tidy
go work sync

# 의존성 설치
echo "의존성 설치 중..."
go get -u github.com/mark3labs/mcp-go
go get -u github.com/xanzy/go-gitlab
go get -u github.com/gin-gonic/gin
go get -u github.com/spf13/viper

# 환경 변수 파일 생성
if [ ! -f .env ]; then
    echo "환경 변수 파일 생성 중..."
    cat > .env << EOF
# theshop-ai-mcp-server 환경 변수

# 서버 설정
THESHOP_SERVER_PORT=8080
THESHOP_SERVER_TIMEOUT=30s

# MCP 서버 설정
MCP_SERVER_ENABLED=true
MCP_SERVER_MODE=http
MCP_SERVER_PORT=8080

# GitLab 설정
THESHOP_GITLAB_URL=https://gitlab.com
THESHOP_GITLAB_TOKEN=your_gitlab_token

# Notion 설정
THESHOP_NOTION_TOKEN=your_notion_token
THESHOP_NOTION_API_KEY=your_notion_api_key

# Figma 설정
THESHOP_FIGMA_ACCESS_TOKEN=your_figma_access_token
THESHOP_FIGMA_TEAM_ID=your_figma_team_id

# LLM 설정
THESHOP_LLM_DEFAULT_MODEL=gpt-3.5-turbo
THESHOP_LLM_DEFAULT_TEMPERATURE=0.7
THESHOP_LLM_PROVIDER=openai
THESHOP_LLM_OPENAI_API_KEY=your_openai_api_key
THESHOP_LLM_CLAUDE_API_KEY=your_claude_api_key
THESHOP_LLM_GEMINI_API_KEY=your_gemini_api_key
EOF
    echo ".env 파일이 생성되었습니다. 실제 API 키로 수정해주세요."
fi

# 설정 파일 확인
if [ ! -f configs/config.yaml ]; then
    echo "configs/config.yaml 파일이 존재하지 않습니다. 기본 파일을 복사합니다."
    mkdir -p configs
    cp configs/config.yaml.example configs/config.yaml
    echo "configs/config.yaml 파일이 생성되었습니다. 필요한 설정을 수정해주세요."
fi

# 빌드
echo "빌드 중..."
go build -o bin/mcpserver ./cmd/mcpserver

# 스크립트에 실행 권한 부여
chmod +x scripts/*.sh

echo "===== 초기화 완료 ====="
echo "서버를 실행하려면: ./scripts/run.sh"