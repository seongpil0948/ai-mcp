#!/bin/bash

PROJECT_ROOT="/Users/2309-n0015/Code/Project/Util/theshop-ai"
MODULES=". ./internal ./modules/config ./modules/core ./modules/integrations/aws ./modules/integrations/confluence ./modules/integrations/filesystem ././modules/integrations ./modules/integrations/jira ./modules/integrations/llm ./modules/integrations/notion ./pkg/mcpclient"

echo "===== go.mod 파일 정리 ====="
for module in $MODULES; do
  echo "정리 중: $module"
  cd "$PROJECT_ROOT/$module" && go mod tidy
done
