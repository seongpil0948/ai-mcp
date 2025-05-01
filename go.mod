module github.com/theshop/ai

go 1.24.2

// Prevent git repository lookups for internal modules
replace (
	github.com/theshop/ai/cmd/api => ./cmd/api
	github.com/theshop/ai/cmd/mcpserver => ./cmd/mcpserver
	github.com/theshop/ai/internal => ./internal
	github.com/theshop/ai/internal/app => ./internal/app
	github.com/theshop/ai/internal/domain => ./internal/domain
	github.com/theshop/ai/internal/transport/http/handlers => ./internal/transport/http/handlers
	github.com/theshop/ai/modules/config => ./modules/config
	github.com/theshop/ai/modules/core => ./modules/core
	github.com/theshop/ai/modules/integrations => ./modules/integrations
	github.com/theshop/ai/modules/integrations/aws => ./modules/integrations/aws
	github.com/theshop/ai/modules/integrations/confluence => ./modules/integrations/confluence
	github.com/theshop/ai/modules/integrations/converter => ./modules/integrations/converter
	github.com/theshop/ai/modules/integrations/figma => ./modules/integrations/figma
	github.com/theshop/ai/modules/integrations/gitlab => ./modules/integrations/gitlab
	github.com/theshop/ai/modules/integrations/jira => ./modules/integrations/jira
	github.com/theshop/ai/modules/integrations/llm => ./modules/integrations/llm
	github.com/theshop/ai/modules/integrations/notion => ./modules/integrations/notion
	github.com/theshop/ai/pkg/mcpclient => ./pkg/mcpclient
	github.com/theshop/ai/pkg/mcpserver => ./pkg/mcpserver
)
