#!/bin/bash

# 스크립트 경로 설정
PROJECT_ROOT="/Users/2309-n0015/Code/Project/Util/theshop-ai"

# 함수: 모듈 문제 디버깅
debug_module_issues() {
  local module_path=$1
  echo "======= $module_path 모듈 디버깅 ======="
  cd "$PROJECT_ROOT/$module_path"
  
  # go.mod 파일 검사
  echo "go.mod 파일 내용:"
  cat go.mod
  echo ""
  
  # 모듈 확인
  echo "go mod verify 결과:"
  go mod verify
  echo ""
  
  # 의존성 그래프 출력
  echo "의존성 그래프:"
  go mod graph
  echo ""
}

# 함수: MCPClient 패키지 수정
fix_mcpclient_package() {
  echo "===== MCPClient 패키지 수정 ====="
  cd "$PROJECT_ROOT/pkg/mcpclient"
  
  # mark3labs/mcp-go 의존성 추가
  echo "mark3labs/mcp-go 의존성 추가..."
  go get -u github.com/mark3labs/mcp-go@latest
  
  # go.mod 정리
  go mod tidy
}

# 함수: internal 모듈의 app 패키지 수정
fix_app_package() {
  echo "===== Internal/App 패키지 수정 ====="
  cd "$PROJECT_ROOT/internal/app"
  
  # service.go 파일 분석 및 수정
  echo "service.go 파일 분석..."
  if grep -q "github.com/theshop/ai/modules/integrations" service.go; then
    echo "service.go 파일에서 통합 모듈 import 확인됨. 수정 필요..."
    # 여기에 필요한 수정 코드 추가
  fi
  
  # workspace.go 파일 분석 및 수정
  echo "workspace.go 파일 분석..."
  if grep -q "github.com/theshop/ai/modules/integrations" workspace.go; then
    echo "workspace.go 파일에서 통합 모듈 import 확인됨. 수정 필요..."
    # 여기에 필요한 수정 코드 추가
  fi
}

# 함수: 모든 replace 지시문 검사 및 수정
fix_all_replace_directives() {
  echo "===== 모든 Replace 지시문 검사 및 수정 ====="
  
  # 모든 go.mod 파일 찾기
  find "$PROJECT_ROOT" -name "go.mod" | while read -r mod_file; do
    echo "검사 중: $mod_file"
    
    # 현재 디렉토리 계산
    module_dir=$(dirname "$mod_file")
    rel_path=$(realpath --relative-to="$module_dir" "$PROJECT_ROOT")
    
    # 루트 모듈에 대한 replace 추가
    if ! grep -q "replace github.com/theshop/ai =>" "$mod_file"; then
      echo "루트 모듈 replace 추가: $mod_file"
      echo "replace github.com/theshop/ai => $rel_path" >> "$mod_file"
    fi
    
    # 필요한 모든 하위 모듈에 대한 replace 추가
    # ... (필요에 따라 더 많은 로직 추가)
  done
}

# 메인 함수
main() {
  echo "===== theshop-ai 프로젝트 import 및 replace 문제 수정 스크립트 ====="
  
  # 모든 replace 지시문 수정
  fix_all_replace_directives
  
  # MCPClient 패키지 수정
  fix_mcpclient_package
  
  # app 패키지 문제 수정
  fix_app_package
  
  # 문제 있는 모듈 디버깅
  debug_module_issues "internal"
  debug_module_issues "modules/core"
  debug_module_issues "modules/integrations/jira"
  debug_module_issues "modules/integrations/gitlab"
  debug_module_issues "modules/integrations/llm"
  
  echo "===== 스크립트 완료 ====="
  echo "이제 루트 디렉토리에서 다음 명령을 실행하세요:"
  echo "make clean tidy sync build"
}

# 스크립트 실행
main
