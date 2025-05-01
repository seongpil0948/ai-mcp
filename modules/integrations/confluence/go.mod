module github.com/theshop/ai/modules/integrations/confluence

go 1.24.2

replace (
	github.com/theshop/ai => ../../..
	github.com/theshop/ai/modules/config => ../../config
	github.com/theshop/ai/pkg/mcpserver => ../../../pkg/mcpserver
)
