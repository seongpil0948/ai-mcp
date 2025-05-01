module github.com/theshop/ai/internal

go 1.24.2

replace (
	github.com/theshop/ai => ../
	github.com/theshop/ai/internal/app => ../internal/app
	github.com/theshop/ai/internal/domain => ../internal/domain
	github.com/theshop/ai/modules/config => ../modules/config
	github.com/theshop/ai/modules/core => ../modules/core
	github.com/theshop/ai/modules/integrations => ../modules/integrations
	github.com/theshop/ai/modules/integrations/aws => ../modules/integrations/aws
	github.com/theshop/ai/modules/integrations/confluence => ../modules/integrations/confluence
	github.com/theshop/ai/modules/integrations/filesystem => ../modules/integrations/filesystem
	github.com/theshop/ai/modules/integrations/jira => ../modules/integrations/jira
	github.com/theshop/ai/modules/integrations/llm => ../modules/integrations/llm
	github.com/theshop/ai/modules/integrations/notion => ../modules/integrations/notion
	github.com/theshop/ai/pkg/mcpclient => ../pkg/mcpclient
)
