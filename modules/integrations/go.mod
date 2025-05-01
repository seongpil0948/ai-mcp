module github.com/theshop/ai/modules/integrations

go 1.24.2

replace (
	github.com/theshop/ai => ../..
	github.com/theshop/ai/internal => ../../internal
	github.com/theshop/ai/modules/config => ../config
	github.com/theshop/ai/modules/integrations => .
)
