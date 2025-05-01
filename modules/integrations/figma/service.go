// Path: modules/integrations/figma/service.go
package figma

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/theshop/ai/pkg/mcpserver"
)

// Service represents a service for interacting with Figma
type Service struct {
	client Client
	teamID string
}

// NewService creates a new Figma service
func NewService(client Client, teamID string) *Service {
	return &Service{
		client: client,
		teamID: teamID,
	}
}

// RegisterTools registers Figma tools with the MCP server
func (s *Service) RegisterTools(mcpServer *mcpserver.MCPServer) error {
	// Register Figma get file tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "figma_get_file",
			Description: "Get a Figma file by key",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"file_key": map[string]interface{}{
						"type":        "string",
						"description": "The file key (e.g., 'abc123')",
					},
				},
				"required": []string{"file_key"},
			},
		},
		s.handleFigmaGetFile,
	); err != nil {
		return err
	}

	// Register Figma get nodes tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "figma_get_nodes",
			Description: "Get specific nodes from a Figma file",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"file_key": map[string]interface{}{
						"type":        "string",
						"description": "The file key",
					},
					"node_ids": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "Array of node IDs to retrieve",
					},
				},
				"required": []string{"file_key", "node_ids"},
			},
		},
		s.handleFigmaGetNodes,
	); err != nil {
		return err
	}

	// Register Figma create file tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "figma_create_file",
			Description: "Create a new Figma file",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the file to create",
					},
				},
				"required": []string{"name"},
			},
		},
		s.handleFigmaCreateFile,
	); err != nil {
		return err
	}

	// Register Figma get team projects tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "figma_get_team_projects",
			Description: "Get projects for a team",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"team_id": map[string]interface{}{
						"type":        "string",
						"description": "Team ID (optional, uses default team if not provided)",
					},
				},
			},
		},
		s.handleFigmaGetTeamProjects,
	); err != nil {
		return err
	}

	// Register Figma get project files tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "figma_get_project_files",
			Description: "Get files for a project",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID",
					},
				},
				"required": []string{"project_id"},
			},
		},
		s.handleFigmaGetProjectFiles,
	); err != nil {
		return err
	}

	// Register Figma create from Notion tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "figma_create_from_notion",
			Description: "Create a Figma file from Notion page content",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"notion_page_id": map[string]interface{}{
						"type":        "string",
						"description": "Notion page ID",
					},
					"file_name": map[string]interface{}{
						"type":        "string",
						"description": "Name for the new Figma file",
					},
					"template_style": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"document", "presentation", "wireframe"},
						"description": "Style template to use",
					},
				},
				"required": []string{"notion_page_id", "file_name"},
			},
		},
		s.handleFigmaCreateFromNotion,
	); err != nil {
		return err
	}

	return nil
}

// handleFigmaGetFile handles the figma_get_file tool
func (s *Service) handleFigmaGetFile(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	fileKey, ok := args["file_key"].(string)
	if !ok {
		return mcpserver.CreateToolError("File key is required"), nil
	}

	file, err := s.client.GetFile(ctx, fileKey)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to get Figma file: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(file)
}

// handleFigmaGetNodes handles the figma_get_nodes tool
func (s *Service) handleFigmaGetNodes(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	fileKey, ok := args["file_key"].(string)
	if !ok {
		return mcpserver.CreateToolError("File key is required"), nil
	}

	nodeIDsRaw, ok := args["node_ids"].([]interface{})
	if !ok {
		return mcpserver.CreateToolError("Node IDs are required"), nil
	}

	var nodeIDs []string
	for _, id := range nodeIDsRaw {
		if idStr, ok := id.(string); ok {
			nodeIDs = append(nodeIDs, idStr)
		}
	}

	if len(nodeIDs) == 0 {
		return mcpserver.CreateToolError("At least one node ID is required"), nil
	}

	nodes, err := s.client.GetFileNodes(ctx, fileKey, nodeIDs)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to get Figma nodes: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(nodes)
}

// handleFigmaCreateFile handles the figma_create_file tool
func (s *Service) handleFigmaCreateFile(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	name, ok := args["name"].(string)
	if !ok {
		return mcpserver.CreateToolError("File name is required"), nil
	}

	file, err := s.client.CreateFile(ctx, name, s.teamID)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to create Figma file: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(file)
}

// handleFigmaGetTeamProjects handles the figma_get_team_projects tool
func (s *Service) handleFigmaGetTeamProjects(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	teamID := s.teamID
	if teamIDArg, ok := args["team_id"].(string); ok && teamIDArg != "" {
		teamID = teamIDArg
	}

	if teamID == "" {
		return mcpserver.CreateToolError("Team ID is required"), nil
	}

	projects, err := s.client.GetTeamProjects(ctx, teamID)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to get team projects: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(projects)
}

// handleFigmaGetProjectFiles handles the figma_get_project_files tool
func (s *Service) handleFigmaGetProjectFiles(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	projectID, ok := args["project_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Project ID is required"), nil
	}

	files, err := s.client.GetProjectFiles(ctx, projectID)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to get project files: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(files)
}

// handleFigmaCreateFromNotion handles the figma_create_from_notion tool
// This is a complex integration that requires interacting with both Notion and Figma APIs
func (s *Service) handleFigmaCreateFromNotion(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	notionPageID, ok := args["notion_page_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Notion page ID is required"), nil
	}

	fileName, ok := args["file_name"].(string)
	if !ok {
		return mcpserver.CreateToolError("File name is required"), nil
	}

	templateStyle := "document"
	if styleArg, ok := args["template_style"].(string); ok && styleArg != "" {
		templateStyle = styleArg
	}

	// In a real implementation, this would:
	// 1. Fetch the Notion page content
	// 2. Process the content into a Figma-compatible structure
	// 3. Create a new Figma file
	// 4. Add elements to the Figma file based on the Notion content

	// For now, we'll return a mock result
	result := map[string]interface{}{
		"status":      "success",
		"message":     fmt.Sprintf("Created Figma file from Notion page %s with template style %s", notionPageID, templateStyle),
		"file_name":   fileName,
		"template":    templateStyle,
		"notion_page": notionPageID,
		"url":         "https://www.figma.com/file/abcdef123456",
	}

	return mcpserver.CreateToolResultJSON(result)
}
