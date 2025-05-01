package mcpclient

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCPClient MCP 서버와의 통신을 위한 인터페이스
type MCPClient interface {
	Initialize(ctx context.Context) error
	ListTools(ctx context.Context) ([]mcp.Tool, error)
	CallTool(ctx context.Context, name string, args map[string]interface{}) (*mcp.CallToolResult, error)
	ReadResource(ctx context.Context, uri string) ([]mcp.ResourceContents, error)
	Close() error
}

// MCPClientConfig MCP 클라이언트 설정
type MCPClientConfig struct {
	Type    string            // "http", "sse", "stdio"
	URL     string            // HTTP 또는 SSE 타입일 경우 사용
	Command string            // stdio 타입일 경우 사용
	Args    []string          // stdio 타입일 경우 사용
	Env     map[string]string // stdio 타입일 경우 사용
	Headers map[string]string // HTTP 또는 SSE 타입일 경우 사용
}

// MCPClientImpl MCP 클라이언트 구현체
type MCPClientImpl struct {
	client *client.Client
	name   string
}

// NewMCPClient 새 MCP 클라이언트 생성
func NewMCPClient(name string, config MCPClientConfig) (MCPClient, error) {
	var mcpClient *client.Client
	var err error

	switch strings.ToLower(config.Type) {
	case "http":
		var options []transport.ClientOption
		if len(config.Headers) > 0 {
			options = append(options, client.WithHeaders(config.Headers))
		}
		mcpClient, err = client.NewHTTPMCPClient(config.URL, options...)

	case "sse":
		var options []transport.ClientOption
		if len(config.Headers) > 0 {
			options = append(options, client.WithHeaders(config.Headers))
		}
		mcpClient, err = client.NewSSEMCPClient(config.URL, options...)

	case "stdio":
		envVars := make([]string, 0, len(config.Env))
		for k, v := range config.Env {
			envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
		}
		mcpClient, err = client.NewStdioMCPClient(config.Command, envVars, config.Args...)

	default:
		return nil, fmt.Errorf("지원하지 않는 MCP 클라이언트 타입: %s", config.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("MCP 클라이언트 생성 오류: %w", err)
	}

	return &MCPClientImpl{
		client: mcpClient,
		name:   name,
	}, nil
}

// Initialize MCP 클라이언트 초기화
func (c *MCPClientImpl) Initialize(ctx context.Context) error {
	// 클라이언트 정보 생성
	clientInfo := mcp.Implementation{
		Name:    "TheSHOP-AI",
		Version: "1.0.0",
	}

	// 초기화 요청 생성
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = clientInfo
	initRequest.Params.Capabilities = mcp.ClientCapabilities{
		Experimental: make(map[string]interface{}),
	}

	// 초기화 요청 전송
	_, err := c.client.Initialize(ctx, initRequest)
	if err != nil {
		return fmt.Errorf("MCP 클라이언트 초기화 오류: %w", err)
	}

	return nil
}

// ListTools MCP 서버가 제공하는 도구 목록 조회
func (c *MCPClientImpl) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := c.client.ListTools(ctx, toolsRequest)
	if err != nil {
		return nil, fmt.Errorf("도구 목록 조회 오류: %w", err)
	}

	return tools.Tools, nil
}

// CallTool MCP 도구 호출
func (c *MCPClientImpl) CallTool(ctx context.Context, name string, args map[string]interface{}) (*mcp.CallToolResult, error) {
	request := mcp.CallToolRequest{}
	request.Params.Name = name
	request.Params.Arguments = args

	result, err := c.client.CallTool(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("도구 호출 오류 (%s): %w", name, err)
	}

	return result, nil
}

// ReadResource MCP 리소스 읽기
func (c *MCPClientImpl) ReadResource(ctx context.Context, uri string) ([]mcp.ResourceContents, error) {
	request := mcp.ReadResourceRequest{}
	request.Params.URI = uri

	result, err := c.client.ReadResource(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("리소스 읽기 오류 (%s): %w", uri, err)
	}

	if len(result.Contents) == 0 {
		return nil, errors.New("리소스에 콘텐츠가 없습니다")
	}

	return result.Contents, nil
}

// Close MCP 클라이언트 연결 종료
func (c *MCPClientImpl) Close() error {
	return c.client.Close()
}
