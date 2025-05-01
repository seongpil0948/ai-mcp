#!/bin/bash

# 스크립트 경로 설정
PROJECT_ROOT="/Users/2309-n0015/Code/Project/Util/theshop-ai"

# 필요한 Go 도구 설치
install_go_tools() {
  echo "===== Go 도구 설치 ====="
  go install github.com/google/wire/cmd/wire@latest
  go install golang.org/x/tools/cmd/goimports@latest
  go install golang.org/x/lint/golint@latest
}

# 프로젝트 빌드 설정 초기화
init_build_settings() {
  echo "===== 빌드 설정 초기화 ====="
  
  # 빌드 디렉토리 생성
  mkdir -p "$PROJECT_ROOT/bin"
  
  # 환경 변수 설정 파일 생성
  cat > "$PROJECT_ROOT/.env" << EOF
# theshop-ai 개발 환경 설정
THESHOP_SERVER_PORT=8080
THESHOP_SERVER_TIMEOUT=30s

# Jira 설정
THESHOP_JIRA_URL=https://your-jira-instance.atlassian.net
THESHOP_JIRA_USERNAME=your-email@example.com
THESHOP_JIRA_API_TOKEN=your-api-token
THESHOP_JIRA_PROJECT_KEY=PROJ

# GitLab 설정
THESHOP_GITLAB_URL=https://gitlab.com
THESHOP_GITLAB_TOKEN=your-gitlab-token
THESHOP_GITLAB_ACCESS_TOKEN=your-gitlab-access-token

# LLM 설정
THESHOP_LLM_DEFAULT_MODEL=gpt-3.5-turbo
THESHOP_LLM_DEFAULT_TEMPERATURE=0.7
THESHOP_LLM_PROVIDER=openai
THESHOP_LLM_OPENAI_API_KEY=your-openai-api-key
THESHOP_LLM_CLAUDE_API_KEY=your-claude-api-key
THESHOP_LLM_GEMINI_API_KEY=your-gemini-api-key
EOF

  echo ".env 파일이 생성되었습니다. 실제 API 키로 업데이트해주세요."
}

# 설정 파일 초기화
init_config_files() {
  echo "===== 설정 파일 초기화 ====="
  
  # configs 디렉토리 확인 및 생성
  mkdir -p "$PROJECT_ROOT/configs"
  
  # 기본 config.yaml 파일 생성
  cat > "$PROJECT_ROOT/configs/config.yaml" << EOF
# theshop-ai 서버 설정
server:
  port: 8080
  timeout: 30s

# Jira 설정
jira:
  url: "https://your-jira-instance.atlassian.net"
  username: "your-email@example.com"
  api_token: "your-api-token"
  project_key: "PROJ"

# GitLab 설정
gitlab:
  url: "https://gitlab.com"
  token: "your-gitlab-token"
  access_token: "your-gitlab-access-token"

# Confluence 설정
confluence:
  url: "https://your-confluence-instance.atlassian.net"
  username: "your-email@example.com"
  api_token: "your-api-token"

# AWS 설정
aws:
  region: "us-east-1"
  access_key_id: "your-access-key-id"
  secret_access_key: "your-secret-access-key"

# Notion 설정
notion:
  token: "your-notion-token"
  api_key: "your-notion-api-key"

# LLM 설정
llm:
  default_model: "gpt-3.5-turbo"
  default_temperature: 0.7
  provider: "openai"
  
  openai:
    api_key: "your-openai-api-key"
    org_id: "your-openai-org-id"
    default_model: "gpt-3.5-turbo"
  
  claude:
    api_key: "your-claude-api-key"
    default_model: "claude-3-haiku-20240307"
  
  gemini:
    api_key: "your-gemini-api-key"
    default_model: "gemini-pro"

# MCP 서버 설정
mcp_servers:
  openai:
    type: "http"
    url: "http://localhost:8081"
  
  claude:
    type: "http"
    url: "http://localhost:8082"
EOF

  echo "config.yaml 파일이 생성되었습니다. 실제 값으로 업데이트해주세요."
}

# Git hooks 설정
setup_git_hooks() {
  echo "===== Git Hooks 설정 ====="
  
  # pre-commit 훅 디렉토리 확인 및 생성
  mkdir -p "$PROJECT_ROOT/.git/hooks"
  
  # pre-commit 훅 생성
  cat > "$PROJECT_ROOT/.git/hooks/pre-commit" << 'EOF'
#!/bin/bash

# Go 포맷 검사
echo "Go 코드 포맷 검사..."
GOFMT_FILES=$(gofmt -l .)
if [[ -n "$GOFMT_FILES" ]]; then
    echo "다음 파일들의 형식을 수정해주세요:"
    echo "$GOFMT_FILES"
    exit 1
fi

# Go 린트 검사
echo "Go 코드 린트 검사..."
LINT_RESULT=$(golint ./... | grep -v "should have comment")
if [[ -n "$LINT_RESULT" ]]; then
    echo "다음 린트 이슈들을 수정해주세요:"
    echo "$LINT_RESULT"
    exit 1
fi

# Go vet 검사
echo "Go 코드 정적 분석 검사..."
go vet ./...
if [[ $? -ne 0 ]]; then
    echo "go vet 이슈들을 수정해주세요."
    exit 1
fi

echo "모든 검사를 통과했습니다!"
exit 0
EOF

  # pre-commit 훅에 실행 권한 부여
  chmod +x "$PROJECT_ROOT/.git/hooks/pre-commit"
  
  echo "Git hooks가 설정되었습니다."
}

# VS Code 설정
setup_vscode() {
  echo "===== VS Code 설정 ====="
  
  # .vscode 디렉토리 확인 및 생성
  mkdir -p "$PROJECT_ROOT/.vscode"
  
  # settings.json 생성
  cat > "$PROJECT_ROOT/.vscode/settings.json" << EOF
{
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golint",
    "go.buildOnSave": "workspace",
    "go.testOnSave": false,
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    },
    "files.exclude": {
        "**/.git": true,
        "**/.DS_Store": true,
        "**/bin": true,
        "**/*.sum": true
    },
    "[go]": {
        "editor.defaultFormatter": "golang.go",
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": true
        }
    }
}
EOF

  # launch.json 생성
  cat > "$PROJECT_ROOT/.vscode/launch.json" << EOF
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch API Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/api/main.go",
            "env": {},
            "args": []
        },
        {
            "name": "Test Current File",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${fileDirname}",
            "env": {},
            "args": []
        }
    ]
}
EOF

  echo "VS Code 설정이 완료되었습니다."
}

# 메인 함수
main() {
  echo "===== theshop-ai 개발 환경 설정 스크립트 ====="
  
  # Go 도구 설치
  install_go_tools
  
  # 빌드 설정 초기화
  init_build_settings
  
  # 설정 파일 초기화
  init_config_files
  
  # Git hooks 설정
  setup_git_hooks
  
  # VS Code 설정
  setup_vscode
  
  echo "===== 개발 환경 설정 완료 ====="
  echo "이제 다음 명령을 실행하여 프로젝트를 빌드하세요:"
  echo "cd $PROJECT_ROOT && make all"
}

# 스크립트 실행
main
