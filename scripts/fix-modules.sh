#!/bin/bash

PROJECT_ROOT="/Users/2309-n0015/Code/Project/Util/theshop-ai"

echo "===== 모듈 파일 수정 중... ====="

# 루트 모듈 수정
echo "루트 go.mod 파일 수정..."
sed -i '.bak' 's/^go 1.24.2$/go 1.24/' "$PROJECT_ROOT/go.mod"

# internal 모듈 수정
echo "internal 모듈 수정..."
sed -i '.bak' 's/^go 1.24.2$/go 1.24/' "$PROJECT_ROOT/internal/go.mod"
cat > "$PROJECT_ROOT/internal/go.mod" << 'EOF'
module github.com/theshop/ai/internal

go 1.24

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/google/uuid v1.6.0
	github.com/theshop/ai/modules/integrations/gitlab v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/modules/integrations/jira v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/modules/integrations/llm v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/pkg/mcpclient v0.0.0-00010101000000-000000000000
)


replace (
	github.com/theshop/ai => ../
	github.com/theshop/ai/modules/config => ../modules/config
	github.com/theshop/ai/modules/core => ../modules/core
	github.com/theshop/ai/modules/integrations/aws => ../modules/integrations/aws
	github.com/theshop/ai/modules/integrations/confluence => ../modules/integrations/confluence
	github.com/theshop/ai/modules/integrations/filesystem => ../modules/integrations/filesystem
	github.com/theshop/ai/modules/integrations/gitlab => ../modules/integrations/gitlab
	github.com/theshop/ai/modules/integrations/jira => ../modules/integrations/jira
	github.com/theshop/ai/modules/integrations/llm => ../modules/integrations/llm
	github.com/theshop/ai/modules/integrations/notion => ../modules/integrations/notion
	github.com/theshop/ai/pkg/mcpclient => ../pkg/mcpclient
)
EOF

# modules/config 모듈 수정
echo "config 모듈 수정..."
cat > "$PROJECT_ROOT/modules/config/go.mod" << 'EOF'
module github.com/theshop/ai/modules/config

go 1.24

require github.com/spf13/viper v1.18.2


replace (
	github.com/theshop/ai => ../..
	github.com/theshop/ai/internal => ../../internal
)
EOF

# modules/core 모듈 수정
echo "core 모듈 수정..."
cat > "$PROJECT_ROOT/modules/core/go.mod" << 'EOF'
module github.com/theshop/ai/modules/core

go 1.24

require (
	github.com/google/wire v0.6.0
	github.com/theshop/ai/internal v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/modules/config v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/modules/integrations/gitlab v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/modules/integrations/jira v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/modules/integrations/llm v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/pkg/mcpclient v0.0.0-00010101000000-000000000000
)


replace (
	github.com/theshop/ai => ../..
	github.com/theshop/ai/internal => ../../internal
	github.com/theshop/ai/modules/config => ../config
	github.com/theshop/ai/modules/integrations/aws => ../integrations/aws
	github.com/theshop/ai/modules/integrations/confluence => ../integrations/confluence
	github.com/theshop/ai/modules/integrations/filesystem => ../integrations/filesystem
	github.com/theshop/ai/modules/integrations/gitlab => ../integrations/gitlab
	github.com/theshop/ai/modules/integrations/jira => ../integrations/jira
	github.com/theshop/ai/modules/integrations/llm => ../integrations/llm
	github.com/theshop/ai/modules/integrations/notion => ../integrations/notion
	github.com/theshop/ai/pkg/mcpclient => ../../pkg/mcpclient
)
EOF

# modules/integrations/jira 모듈 수정
echo "jira 모듈 수정..."
cat > "$PROJECT_ROOT/modules/integrations/jira/go.mod" << 'EOF'
module github.com/theshop/ai/modules/integrations/jira

go 1.24

require (
	github.com/theshop/ai/internal v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/modules/config v0.0.0-00010101000000-000000000000
)


replace (
	github.com/theshop/ai => ../../..
	github.com/theshop/ai/internal => ../../../internal
	github.com/theshop/ai/modules/config => ../../config
)
EOF

# modules/integrations/gitlab 모듈 수정
echo "gitlab 모듈 수정..."
cat > "$PROJECT_ROOT/modules/integrations/gitlab/go.mod" << 'EOF'
module github.com/theshop/ai/modules/integrations/gitlab

go 1.24

require (
	github.com/theshop/ai/internal v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/modules/config v0.0.0-00010101000000-000000000000
	github.com/xanzy/go-gitlab v0.96.0
)


replace (
	github.com/theshop/ai => ../../..
	github.com/theshop/ai/internal => ../../../internal
	github.com/theshop/ai/modules/config => ../../config
)
EOF

# modules/integrations/llm 모듈 수정
echo "llm 모듈 수정..."
cat > "$PROJECT_ROOT/modules/integrations/llm/go.mod" << 'EOF'
module github.com/theshop/ai/modules/integrations/llm

go 1.24

require (
	github.com/theshop/ai/modules/config v0.0.0-00010101000000-000000000000
	github.com/theshop/ai/pkg/mcpclient v0.0.0-00010101000000-000000000000
)


replace (
	github.com/theshop/ai => ../../..
	github.com/theshop/ai/internal => ../../../internal
	github.com/theshop/ai/modules/config => ../../config
	github.com/theshop/ai/pkg/mcpclient => ../../../pkg/mcpclient
)
EOF

# pkg/mcpclient 모듈 수정
echo "mcpclient 모듈 수정..."
cat > "$PROJECT_ROOT/pkg/mcpclient/go.mod" << 'EOF'
module github.com/theshop/ai/pkg/mcpclient

go 1.24

require github.com/mark3labs/mcp-go v0.24.1


replace (
	github.com/theshop/ai => ../..
	github.com/theshop/ai/internal => ../../internal
)
EOF

# go.work 파일 수정
echo "go.work 파일 수정..."
cat > "$PROJECT_ROOT/go.work" << 'EOF'
go 1.24

use (
	.
	./internal
	./modules/config
	./modules/core
	./modules/integrations/aws
	./modules/integrations/confluence
	./modules/integrations/filesystem
	./modules/integrations/gitlab
	./modules/integrations/jira
	./modules/integrations/llm
	./modules/integrations/notion
	./pkg/mcpclient
)
EOF

echo "모든 모듈 파일 수정 완료!"
