module github.com/theshop/ai/pkg/mcpserver // Ensure path is standard

go 1.24.2 // Or your specific Go version

require (
	github.com/mark3labs/mcp-go v0.25.0
	google.golang.org/protobuf v1.34.1 // Keep for structpb
	github.com/theshop/ai/modules/config v0.0.0 // Add pseudo-version
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
)

// Remove all replace directives as go.work handles them -- Re-adding necessary ones
replace github.com/theshop/ai/modules/config => ../../modules/config
