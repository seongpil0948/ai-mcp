package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/theshop/ai/modules/config"
	"github.com/theshop/ai/pkg/mcpclient"
	"google.golang.org/protobuf/types/known/structpb" // Import structpb
)

// GenerateOptions 텍스트 생성 옵션
type GenerateOptions struct {
	Model       string
	Temperature float64
	MaxTokens   int
	System      string
}

// Service LLM 서비스 인터페이스
type Service interface {
	GenerateText(ctx context.Context, prompt string, options GenerateOptions) (string, error)
	AnalyzeData(ctx context.Context, data []byte, query string) (string, error)
}

// MCPLLMService MCP를 통한 LLM 서비스 구현
type MCPLLMService struct {
	mcpClient    mcpclient.MCPClient
	defaultModel string
	provider     string
}

// NewMCPLLMService 새 MCP 기반 LLM 서비스 생성
func NewMCPLLMService(mcpClient mcpclient.MCPClient, config *config.LLMConfig) Service {
	return &MCPLLMService{
		mcpClient:    mcpClient,
		defaultModel: config.DefaultModel,
		provider:     config.Provider,
	}
}

// GenerateText LLM을 사용하여 텍스트 생성
func (s *MCPLLMService) GenerateText(ctx context.Context, prompt string, options GenerateOptions) (string, error) {
	model := options.Model
	if model == "" {
		model = s.defaultModel
	}

	temperature := options.Temperature
	if temperature <= 0 {
		temperature = 0.7
	}

	maxTokens := options.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 1000
	}

	systemPrompt := options.System
	if systemPrompt == "" {
		systemPrompt = "You are a helpful AI assistant."
	}

	// 프로바이더에 따른 도구 선택
	var toolName string
	switch s.provider {
	case "openai":
		toolName = "generate_text"
	case "claude":
		toolName = "anthropic_completion"
	case "gemini":
		toolName = "gemini_generate"
	default:
		toolName = "generate_text" // 기본값
	}

	// MCP 도구 호출 인자 구성
	args := map[string]interface{}{
		"model":       model,
		"prompt":      prompt,
		"temperature": temperature,
		"max_tokens":  maxTokens,
		"system":      systemPrompt,
	}

	// MCP 도구 호출
	result, err := s.mcpClient.CallTool(ctx, toolName, args)
	if err != nil {
		return "", fmt.Errorf("LLM 호출 오류: %w", err)
	}

	// 응답 처리
	if result.Error != "" {
		return "", fmt.Errorf("LLM 오류: %s", result.Error)
	}

	// 콘텐츠 추출 방식은 MCP 도구 반환 형식에 따라 다를 수 있음
	var text string

	// 텍스트 콘텐츠인 경우
	if len(result.Content) > 0 {
		text = result.Content[0].Text
	}

	// JSON 데이터가 있는 경우
	if text == "" && result.Data != nil {
		if content, ok := result.Data["text"].(string); ok {
			text = content
		} else if content, ok := result.Data["content"].(string); ok {
			text = content
		} else if content, ok := result.Data["completion"].(string); ok {
			text = content
		} else {
			// 데이터 자체를 JSON 문자열로 변환
			jsonBytes, _ := json.Marshal(result.Data)
			text = string(jsonBytes)
		}
	}

	return text, nil
}

// AnalyzeData 데이터 분석
func (s *MCPLLMService) AnalyzeData(ctx context.Context, data []byte, query string) (string, error) {
	// 데이터 유형에 따라 적절한 도구 호출
	// 여기서는 간단히 텍스트 기반 분석 가정

	// 텍스트로 변환된 데이터에 대한 쿼리 구성
	dataStr := string(data)
	if len(dataStr) > 5000 {
		// 너무 큰 데이터는 축약
		dataStr = dataStr[:5000] + "...(truncated)"
	}

	prompt := fmt.Sprintf("Analyze the following data and answer this question: %s\n\nData:\n%s", query, dataStr)

	return s.GenerateText(ctx, prompt, GenerateOptions{
		System: "You are a data analysis expert. Answer questions concisely based on the provided data.",
	})
}

func (s *Service) callTool(ctx context.Context, toolCall *mcp.ToolCall) (*mcp.CallToolResult, error) {
	// ... existing code ...

	inputStruct, err := structpb.NewStruct(inputMap) // Convert map to structpb.Struct
	if err != nil {
		log.Printf("Error converting input map to struct: %v", err)
        errorStruct, _ := structpb.NewStruct(map[string]interface{}{"message": fmt.Sprintf("Internal error: failed to prepare tool input: %v", err)})
		return &mcp.CallToolResult{
            ToolName: toolCall.GetName(),
            Result: &mcp.CallToolResult_Error{Error: errorStruct},
        }, nil
	}


	req := &mcp.CallToolRequest{
		ToolName: toolCall.GetName(),
		Input:    inputStruct, // Pass structpb.Struct
	}

	log.Printf("Calling tool '%s' via MCP client '%s'", req.GetToolName(), clientName)
	result, err := client.CallTool(ctx, req)
	if err != nil {
		log.Printf("Error calling tool '%s' via MCP client '%s': %v", req.GetToolName(), clientName, err)
        // This 'err' is likely a transport error, return it directly
		return nil, fmt.Errorf("transport error calling tool %s: %w", req.GetToolName(), err)
	}

    // Check if the result itself contains an error from the tool execution
    if toolErr := result.GetError(); toolErr != nil {
         log.Printf("Tool '%s' returned error: %v", req.GetToolName(), toolErr.AsMap())
         // Return the result containing the error, but nil for the Go error
         return result, nil
    }


	log.Printf("Tool '%s' called successfully via MCP client '%s'", req.GetToolName(), clientName)
	return result, nil
}


func (s *Service) processToolCalls(ctx context.Context, toolCalls []*mcp.ToolCall) ([]*mcp.Content, error) {
	// ... existing code ...
		if result == nil {
			// Handle case where callTool returned a transport error (err != nil)
			log.Printf("Skipping result processing for tool %s due to transport error: %v", toolCall.GetName(), callErr)
			// Optionally create an error content block
			errorMap := map[string]interface{}{"message": fmt.Sprintf("Failed to call tool %s: %v", toolCall.GetName(), callErr)}
			errorStruct, _ := structpb.NewStruct(errorMap)
			toolResults = append(toolResults, mcp.NewToolResultContent(toolCall.GetId(), toolCall.GetName(), errorStruct, true)) // Indicate error
			continue
		}

        // Check if the result contains a tool execution error
        if toolErrStruct := result.GetError(); toolErrStruct != nil {
            log.Printf("Tool %s execution failed: %v", toolCall.GetName(), toolErrStruct.AsMap())
			toolResults = append(toolResults, mcp.NewToolResultContent(toolCall.GetId(), toolCall.GetName(), toolErrStruct, true)) // Indicate error
        } else if resultContent := result.GetContent(); resultContent != nil {
		    // Successful tool call, add result content
		    log.Printf("Tool %s execution successful", toolCall.GetName())
			toolResults = append(toolResults, mcp.NewToolResultContent(toolCall.GetId(), toolCall.GetName(), resultContent, false)) // Indicate success
        } else {
            // Handle unexpected case where result has neither error nor content
            log.Printf("Warning: Tool %s returned result with no error and no content", toolCall.GetName())
            errorMap := map[string]interface{}{"message": fmt.Sprintf("Tool %s returned an empty result", toolCall.GetName())}
            errorStruct, _ := structpb.NewStruct(errorMap)
            toolResults = append(toolResults, mcp.NewToolResultContent(toolCall.GetId(), toolCall.GetName(), errorStruct, true)) // Indicate error
        }
	}
	// ... existing code ...
}

// ... rest of the file ...
