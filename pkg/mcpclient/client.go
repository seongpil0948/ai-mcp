package mcpclient

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp" // Use specific package path
)

// MCPClient defines the interface for interacting with MCP servers.
type MCPClient interface {
	GetTools(ctx context.Context) ([]*mcp.Tool, error)
	CallTool(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error)
	// Add other methods like StreamChat if needed
}

// clientWrapper wraps the mcp-go client implementation.
type clientWrapper struct {
	client mcpgo.Client // Use the Client interface from the base package
}

// New creates a new MCPClient based on the target string.
func New(target string) (MCPClient, error) {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		// Assuming NewHTTPClient is now part of the base package
		c, err := mcpgo.NewHTTPClient(target)
		if err != nil {
			return nil, fmt.Errorf("failed to create HTTP MCP client: %w", err)
		}
		return &clientWrapper{client: c}, nil
	}
	if target == "stdio" {
		// Assuming NewStdioClient is now part of the base package
		c, err := mcpgo.NewStdioClient()
		if err != nil {
			return nil, fmt.Errorf("failed to create Stdio MCP client: %w", err)
		}
		return &clientWrapper{client: c}, nil
	}
	return nil, errors.New("unsupported mcp client target")
}

func (w *clientWrapper) GetTools(ctx context.Context) ([]*mcp.Tool, error) {
	return w.client.GetTools(ctx)
}

func (w *clientWrapper) CallTool(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return w.client.CallTool(ctx, req)
}

// Implement other MCPClient interface methods by calling w.client methods
