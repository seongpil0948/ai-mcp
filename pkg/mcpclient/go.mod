module github.com/theshop/ai/pkg/mcpclient

go 1.24.2

require github.com/mark3labs/mcp-go v0.25.0

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
)

replace (
	github.com/theshop/ai => ../..
	github.com/theshop/ai/internal => ../../internal
)
