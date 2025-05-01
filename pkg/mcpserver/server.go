package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPServer represents the main MCP server that communicates with LLM clients
type MCPServer struct {
	name      string
	version   string
	server    *server.MCPServer
	tools     map[string]mcp.Tool
	toolLock  sync.RWMutex
	toolImpls map[string]ToolImplementation
}

// ToolImplementation is a function that implements a tool's functionality
type ToolImplementation func(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error)

// NewMCPServer creates a new MCP server instance
func NewMCPServer(name, version string) *MCPServer {
	mcpSrv := server.NewMCPServer(name, version, server.WithResourceCapabilities(true, true))

	return &MCPServer{
		name:      name,
		version:   version,
		server:    mcpSrv,
		tools:     make(map[string]mcp.Tool),
		toolImpls: make(map[string]ToolImplementation),
	}
}

// RegisterTool registers a new tool with the MCP server
func (s *MCPServer) RegisterTool(tool mcp.Tool, impl ToolImplementation) error {
	s.toolLock.Lock()
	defer s.toolLock.Unlock()

	if _, exists := s.tools[tool.Name]; exists {
		return fmt.Errorf("tool with name %s already registered", tool.Name)
	}

	s.tools[tool.Name] = tool
	s.toolImpls[tool.Name] = impl

	// Register the tool with the MCP server
	s.server.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		impl, ok := s.toolImpls[request.Params.Name]
		if !ok {
			return nil, fmt.Errorf("tool implementation for %s not found", request.Params.Name)
		}

		return impl(ctx, request.Params.Arguments)
	})

	log.Printf("Registered tool: %s", tool.Name)
	return nil
}

// Start starts the MCP server using the provided transport
func (s *MCPServer) Start(ctx context.Context, transport server.Transport) error {
	log.Printf("Starting MCP server %s v%s", s.name, s.version)

	// Start the server
	if err := s.server.Start(ctx, transport); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	return nil
}

// Stop stops the MCP server
func (s *MCPServer) Stop(ctx context.Context) error {
	return s.server.Stop(ctx)
}

// CreateToolResult creates a result object for tool calls
func CreateToolResult(content string, data map[string]interface{}) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			{
				Type: "text",
				Text: content,
			},
		},
		Data: data,
	}
}

// CreateToolResultJSON creates a result object with JSON-formatted content
func CreateToolResultJSON(obj interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result to JSON: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			{
				Type: "text",
				Text: string(jsonBytes),
			},
		},
		Data: map[string]interface{}{
			"result": obj,
		},
	}, nil
}

// CreateToolError creates an error result for tool calls
func CreateToolError(errMsg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Error: errMsg,
	}
}
