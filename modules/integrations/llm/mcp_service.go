// Path: modules/integrations/llm/mcp_service.go
package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/theshop/ai/pkg/mcpserver"
)

// MCPService represents a service for exposing LLM capabilities via MCP
type MCPService struct {
	service Service
}

// NewMCPService creates a new LLM MCP service
func NewMCPService(service Service) *MCPService {
	return &MCPService{
		service: service,
	}
}

// RegisterTools registers LLM tools with the MCP server
func (s *MCPService) RegisterTools(mcpServer *mcpserver.MCPServer) error {
	// Register text generation tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "llm_generate_text",
			Description: "Generate text using an LLM",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"prompt": map[string]interface{}{
						"type":        "string",
						"description": "The text prompt to generate from",
					},
					"model": map[string]interface{}{
						"type":        "string",
						"description": "The model to use (optional, uses default if not specified)",
					},
					"temperature": map[string]interface{}{
						"type":        "number",
						"description": "Controls randomness (0.0-1.0)",
					},
					"max_tokens": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of tokens to generate",
					},
					"system": map[string]interface{}{
						"type":        "string",
						"description": "System instructions for the model",
					},
				},
				"required": []string{"prompt"},
			},
		},
		s.handleLLMGenerateText,
	); err != nil {
		return err
	}

	// Register code analysis tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "llm_analyze_code",
			Description: "Analyze code using an LLM",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code": map[string]interface{}{
						"type":        "string",
						"description": "The code to analyze",
					},
					"language": map[string]interface{}{
						"type":        "string",
						"description": "The programming language of the code",
					},
					"question": map[string]interface{}{
						"type":        "string",
						"description": "Specific question about the code",
					},
					"model": map[string]interface{}{
						"type":        "string",
						"description": "The model to use (optional)",
					},
				},
				"required": []string{"code"},
			},
		},
		s.handleLLMAnalyzeCode,
	); err != nil {
		return err
	}

	// Register content summarization tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "llm_summarize",
			Description: "Summarize content using an LLM",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"content": map[string]interface{}{
						"type":        "string",
						"description": "The content to summarize",
					},
					"length": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"short", "medium", "long"},
						"description": "The desired length of the summary",
					},
					"format": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"paragraph", "bullets", "outline"},
						"description": "The format of the summary",
					},
					"model": map[string]interface{}{
						"type":        "string",
						"description": "The model to use (optional)",
					},
				},
				"required": []string{"content"},
			},
		},
		s.handleLLMSummarize,
	); err != nil {
		return err
	}

	// Register variable name suggestion tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "llm_suggest_variable_names",
			Description: "Suggest variable names based on code context",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code_context": map[string]interface{}{
						"type":        "string",
						"description": "The code context",
					},
					"variable_type": map[string]interface{}{
						"type":        "string",
						"description": "The variable type",
					},
					"variable_purpose": map[string]interface{}{
						"type":        "string",
						"description": "Description of what the variable is used for",
					},
					"language": map[string]interface{}{
						"type":        "string",
						"description": "The programming language",
					},
					"naming_convention": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"camelCase", "snake_case", "PascalCase", "kebab-case"},
						"description": "The naming convention to use",
					},
				},
				"required": []string{"variable_purpose"},
			},
		},
		s.handleLLMSuggestVariableNames,
	); err != nil {
		return err
	}

	// Register SQL query suggestion tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "llm_generate_sql",
			Description: "Generate SQL queries based on a description",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Description of the query to generate",
					},
					"dialect": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"MySQL", "PostgreSQL", "SQLite", "MS SQL", "Oracle"},
						"description": "SQL dialect to use",
					},
					"schema": map[string]interface{}{
						"type":        "string",
						"description": "Database schema information (table definitions)",
					},
				},
				"required": []string{"description"},
			},
		},
		s.handleLLMGenerateSQL,
	); err != nil {
		return err
	}

	return nil
}

// handleLLMGenerateText handles the llm_generate_text tool
func (s *MCPService) handleLLMGenerateText(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	prompt, ok := args["prompt"].(string)
	if !ok {
		return mcpserver.CreateToolError("Prompt is required"), nil
	}

	options := GenerateOptions{}

	if modelArg, ok := args["model"].(string); ok {
		options.Model = modelArg
	}

	if tempArg, ok := args["temperature"].(float64); ok {
		options.Temperature = tempArg
	}

	if maxTokensArg, ok := args["max_tokens"].(float64); ok {
		options.MaxTokens = int(maxTokensArg)
	}

	if systemArg, ok := args["system"].(string); ok {
		options.System = systemArg
	}

	generatedText, err := s.service.GenerateText(ctx, prompt, options)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Text generation failed: %v", err)), nil
	}

	return mcpserver.CreateToolResult(generatedText, map[string]interface{}{
		"model": options.Model,
	}), nil
}

// handleLLMAnalyzeCode handles the llm_analyze_code tool
func (s *MCPService) handleLLMAnalyzeCode(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	code, ok := args["code"].(string)
	if !ok {
		return mcpserver.CreateToolError("Code is required"), nil
	}

	language := "unknown"
	if langArg, ok := args["language"].(string); ok {
		language = langArg
	}

	question := "Analyze this code and provide insights."
	if questionArg, ok := args["question"].(string); ok {
		question = questionArg
	}

	model := ""
	if modelArg, ok := args["model"].(string); ok {
		model = modelArg
	}

	prompt := fmt.Sprintf(`Analyze the following %s code:

```%s
%s
```

%s`, language, language, code, question)

	generatedAnalysis, err := s.service.GenerateText(ctx, prompt, GenerateOptions{
		Model: model,
		System: "You are a coding expert specialized in code analysis, readability, and best practices.",
	})
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Code analysis failed: %v", err)), nil
	}

	return mcpserver.CreateToolResult(generatedAnalysis, map[string]interface{}{
		"language": language,
	}), nil
}

// handleLLMSummarize handles the llm_summarize tool
func (s *MCPService) handleLLMSummarize(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	content, ok := args["content"].(string)
	if !ok {
		return mcpserver.CreateToolError("Content is required"), nil
	}

	length := "medium"
	if lengthArg, ok := args["length"].(string); ok {
		length = lengthArg
	}

	format := "paragraph"
	if formatArg, ok := args["format"].(string); ok {
		format = formatArg
	}

	model := ""
	if modelArg, ok := args["model"].(string); ok {
		model = modelArg
	}

	// Determine token limit based on length
	tokenLimit := 150
	switch length {
	case "short":
		tokenLimit = 100
	case "medium":
		tokenLimit = 250
	case "long":
		tokenLimit = 500
	}

	formatInstructions := ""
	switch format {
	case "paragraph":
		formatInstructions = "Summarize this in a single coherent paragraph."
	case "bullets":
		formatInstructions = "Summarize this as a bulleted list of key points."
	case "outline":
		formatInstructions = "Summarize this as a hierarchical outline with main points and sub-points."
	}

	prompt := fmt.Sprintf(`Summarize the following content:

%s

%s`, content, formatInstructions)

	summary, err := s.service.GenerateText(ctx, prompt, GenerateOptions{
		Model:     model,
		MaxTokens: tokenLimit,
		System:    "You are an expert at summarizing content clearly and concisely while retaining the most important information.",
	})
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Summarization failed: %v", err)), nil
	}

	return mcpserver.CreateToolResult(summary, map[string]interface{}{
		"format": format,
		"length": length,
	}), nil
}

// handleLLMSuggestVariableNames handles the llm_suggest_variable_names tool
func (s *MCPService) handleLLMSuggestVariableNames(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	purpose, ok := args["variable_purpose"].(string)
	if !ok {
		return mcpserver.CreateToolError("Variable purpose is required"), nil
	}

	varType := ""
	if typeArg, ok := args["variable_type"].(string); ok {
		varType = typeArg
	}

	codeContext := ""
	if contextArg, ok := args["code_context"].(string); ok {
		codeContext = contextArg
	}

	language := "Go"
	if langArg, ok := args["language"].(string); ok {
		language = langArg
	}

	convention := "camelCase"
	if convArg, ok := args["naming_convention"].(string); ok {
		convention = convArg
	}

	prompt := fmt.Sprintf(`Suggest 5 appropriate variable names for a %s variable in %s with the following purpose:
"%s"

The names should follow %s naming convention.`, varType, language, purpose, convention)

	if codeContext != "" {
		prompt += fmt.Sprintf(`

Consider this code context:
```
%s
````, codeContext)
	}

	suggestions, err := s.service.GenerateText(ctx, prompt, GenerateOptions{
		System: "You are an expert programmer who follows best practices for variable naming in various programming languages.",
	})
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Variable name suggestion failed: %v", err)), nil
	}

	// Process the response to extract just the variable names
	names := extractVariableNames(suggestions)

	result := map[string]interface{}{
		"suggestions": names,
		"naming_convention": convention,
		"language": language,
	}

	return mcpserver.CreateToolResultJSON(result)
}

// handleLLMGenerateSQL handles the llm_generate_sql tool
func (s *MCPService) handleLLMGenerateSQL(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	description, ok := args["description"].(string)
	if !ok {
		return mcpserver.CreateToolError("Description is required"), nil
	}

	dialect := "PostgreSQL"
	if dialectArg, ok := args["dialect"].(string); ok {
		dialect = dialectArg
	}

	schema := ""
	if schemaArg, ok := args["schema"].(string); ok {
		schema = schemaArg
	}

	prompt := fmt.Sprintf(`Generate a %s SQL query that:
%s`, dialect, description)

	if schema != "" {
		prompt += fmt.Sprintf(`

Using this database schema:
%s`, schema)
	}

	sqlQuery, err := s.service.GenerateText(ctx, prompt, GenerateOptions{
		System: "You are an expert SQL developer. Generate efficient, secure, and well-commented SQL queries.",
	})
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("SQL generation failed: %v", err)), nil
	}

	// Extract just the SQL code if it's in a code block
	sqlQuery = extractSQLCode(sqlQuery)

	return mcpserver.CreateToolResult(sqlQuery, map[string]interface{}{
		"dialect": dialect,
	}), nil
}

// GenerateCodeReviewComment generates a code review comment using LLM
func (s *MCPService) GenerateCodeReviewComment(ctx context.Context, req *mcp.GenerateCodeReviewCommentRequest) (*mcp.GenerateCodeReviewCommentResponse, error) {
	prompt := fmt.Sprintf(
		"Please review the following code diff and provide constructive feedback as a comment. Focus on potential bugs, improvements, and adherence to best practices.\n\n"+
			"File: %s\n"+
			"Diff:\n```diff\n%s\n```\n\n"+
			"Guidelines:\n"+
			"- Be specific and provide actionable suggestions.\n"+
			"- If suggesting code changes, provide clear examples.\n"+
			"- Maintain a positive and collaborative tone.\n"+
			"- Consider context if available (e.g., related issue description: %s).",
		req.GetFileName(),
		req.GetDiff(),
		req.GetContext(), // Assuming context might be relevant, adjust if needed
	)

	comment, err := s.service.GenerateText(ctx, prompt, GenerateOptions{
		Model: "code-review-model",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate code review comment: %w", err)
	}

	return &mcp.GenerateCodeReviewCommentResponse{
		Comment: comment,
	}, nil
}

// Helper functions

// extractVariableNames extracts variable names from LLM suggestions
func extractVariableNames(text string) []string {
	// A simple approach for extracting variable names
	lines := strings.Split(text, "\n")
	var names []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove leading numbers, bullets, dashes, etc.
		line = strings.TrimLeft(line, "0123456789.-*• ")
		line = strings.TrimSpace(line)

		// If there's a colon or dash with explanation, take just the first part
		if idx := strings.IndexAny(line, ":-"); idx > 0 {
			line = strings.TrimSpace(line[:idx])
		}

		// Remove any backticks that might be used for code formatting
		line = strings.Trim(line, "`")

		if line != "" {
			names = append(names, line)
		}
	}

	// Limit to 5 suggestions
	if len(names) > 5 {
		names = names[:5]
	}

	return names
}

// extractSQLCode extracts SQL code from a possibly decorated response
func extractSQLCode(text string) string {
	// Check if the response contains a code block with SQL
	codeStart := strings.Index(text, "```sql")
	if codeStart >= 0 {
		codeStart += 6 // Move past ```sql
		codeEnd := strings.Index(text[codeStart:], "```")
		if codeEnd > 0 {
			return strings.TrimSpace(text[codeStart : codeStart+codeEnd])
		}
	}

	// Check for generic code block
	codeStart = strings.Index(text, "```")
	if codeStart >= 0 {
		codeStart += 3 // Move past ```
		codeEnd := strings.Index(text[codeStart:], "```")
		if codeEnd > 0 {
			return strings.TrimSpace(text[codeStart : codeStart+codeEnd])
		}
	}

	// If no code block, return as is
	return text
}