package notion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/theshop/ai/pkg/mcpserver"
)

// Service represents a service for interacting with Notion
type Service struct {
	client Client
}

// NewService creates a new Notion service
func NewService(client Client) *Service {
	return &Service{
		client: client,
	}
}

// RegisterTools registers Notion tools with the MCP server
func (s *Service) RegisterTools(mcpServer *mcpserver.MCPServer) error {
	// Register Notion search tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "notion_search",
			Description: "Search for pages or databases in Notion",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
					"filter_type": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"page", "database", "all"},
						"description": "Filter results by type (page, database, or all)",
					},
					"page_size": map[string]interface{}{
						"type":        "integer",
						"description": "Number of results to return (max 100)",
					},
				},
				"required": []string{"query"},
			},
		},
		s.handleNotionSearch,
	); err != nil {
		return err
	}

	// Register Notion get page tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "notion_get_page",
			Description: "Get a Notion page by ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"page_id": map[string]interface{}{
						"type":        "string",
						"description": "Page ID",
					},
				},
				"required": []string{"page_id"},
			},
		},
		s.handleNotionGetPage,
	); err != nil {
		return err
	}

	// Register Notion get database tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "notion_get_database",
			Description: "Get a Notion database by ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"database_id": map[string]interface{}{
						"type":        "string",
						"description": "Database ID",
					},
				},
				"required": []string{"database_id"},
			},
		},
		s.handleNotionGetDatabase,
	); err != nil {
		return err
	}

	// Register Notion query database tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "notion_query_database",
			Description: "Query a Notion database",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"database_id": map[string]interface{}{
						"type":        "string",
						"description": "Database ID",
					},
					"filter": map[string]interface{}{
						"type":        "object",
						"description": "Filter criteria",
					},
					"sorts": map[string]interface{}{
						"type":        "array",
						"description": "Sort criteria",
					},
					"page_size": map[string]interface{}{
						"type":        "integer",
						"description": "Number of results to return (max 100)",
					},
				},
				"required": []string{"database_id"},
			},
		},
		s.handleNotionQueryDatabase,
	); err != nil {
		return err
	}

	// Register Notion create page tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "notion_create_page",
			Description: "Create a Notion page",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"parent_type": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"database_id", "page_id"},
						"description": "Type of parent (database_id or page_id)",
					},
					"parent_id": map[string]interface{}{
						"type":        "string",
						"description": "ID of the parent database or page",
					},
					"properties": map[string]interface{}{
						"type":        "object",
						"description": "Page properties",
					},
					"content": map[string]interface{}{
						"type":        "array",
						"description": "Page content blocks",
					},
				},
				"required": []string{"parent_type", "parent_id", "properties"},
			},
		},
		s.handleNotionCreatePage,
	); err != nil {
		return err
	}

	// Register Notion get block children tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "notion_get_block_children",
			Description: "Get the children of a Notion block",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"block_id": map[string]interface{}{
						"type":        "string",
						"description": "Block ID",
					},
					"page_size": map[string]interface{}{
						"type":        "integer",
						"description": "Number of results to return (max 100)",
					},
				},
				"required": []string{"block_id"},
			},
		},
		s.handleNotionGetBlockChildren,
	); err != nil {
		return err
	}

	// Register Notion append block children tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "notion_append_block_children",
			Description: "Append children to a Notion block",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"block_id": map[string]interface{}{
						"type":        "string",
						"description": "Block ID",
					},
					"children": map[string]interface{}{
						"type":        "array",
						"description": "Block children to append",
					},
				},
				"required": []string{"block_id", "children"},
			},
		},
		s.handleNotionAppendBlockChildren,
	); err != nil {
		return err
	}

	return nil
}

// handleNotionSearch handles the notion_search tool
func (s *Service) handleNotionSearch(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	query, ok := args["query"].(string)
	if !ok {
		return mcpserver.CreateToolError("Query is required"), nil
	}

	options := &SearchOptions{
		Query: query,
	}

	if pageSizeRaw, ok := args["page_size"]; ok {
		if pageSize, ok := pageSizeRaw.(float64); ok {
			options.PageSize = int(pageSize)
		}
	}

	if filterTypeRaw, ok := args["filter_type"]; ok {
		if filterType, ok := filterTypeRaw.(string); ok && filterType != "all" {
			options.Filter = map[string]interface{}{
				"property": "object",
				"value":    filterType,
			}
		}
	}

	results, err := s.client.Search(ctx, query, options)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to search Notion: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(results)
}

// handleNotionGetPage handles the notion_get_page tool
func (s *Service) handleNotionGetPage(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	pageID, ok := args["page_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Page ID is required"), nil
	}

	page, err := s.client.GetPage(ctx, pageID)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to get Notion page: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(page)
}

// handleNotionGetDatabase handles the notion_get_database tool
func (s *Service) handleNotionGetDatabase(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	databaseID, ok := args["database_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Database ID is required"), nil
	}

	database, err := s.client.GetDatabase(ctx, databaseID)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to get Notion database: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(database)
}

// handleNotionQueryDatabase handles the notion_query_database tool
func (s *Service) handleNotionQueryDatabase(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	databaseID, ok := args["database_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Database ID is required"), nil
	}

	query := &DatabaseQuery{}

	if filterRaw, ok := args["filter"]; ok {
		query.Filter = filterRaw
	}

	if sortsRaw, ok := args["sorts"]; ok {
		if sorts, ok := sortsRaw.([]interface{}); ok {
			query.Sorts = sorts
		}
	}

	if pageSizeRaw, ok := args["page_size"]; ok {
		if pageSize, ok := pageSizeRaw.(float64); ok {
			query.PageSize = int(pageSize)
		}
	}

	results, err := s.client.QueryDatabase(ctx, databaseID, query)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to query Notion database: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(results)
}

// handleNotionCreatePage handles the notion_create_page tool
func (s *Service) handleNotionCreatePage(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	parentType, ok := args["parent_type"].(string)
	if !ok {
		return mcpserver.CreateToolError("Parent type is required"), nil
	}

	parentID, ok := args["parent_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Parent ID is required"), nil
	}

	propertiesRaw, ok := args["properties"].(map[string]interface{})
	if !ok {
		return mcpserver.CreateToolError("Properties are required"), nil
	}

	var children []interface{}
	if childrenRaw, ok := args["content"].([]interface{}); ok {
		children = childrenRaw
	}

	params := &CreatePageParams{
		Parent: Parent{
			Type: parentType,
		},
		Properties: propertiesRaw,
		Children:   children,
	}

	if parentType == "database_id" {
		params.Parent.DatabaseID = parentID
	} else if parentType == "page_id" {
		params.Parent.PageID = parentID
	} else {
		return mcpserver.CreateToolError("Invalid parent type, must be 'database_id' or 'page_id'"), nil
	}

	page, err := s.client.CreatePage(ctx, params)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to create Notion page: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(page)
}

// handleNotionGetBlockChildren handles the notion_get_block_children tool
func (s *Service) handleNotionGetBlockChildren(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	blockID, ok := args["block_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Block ID is required"), nil
	}

	options := &PaginationOptions{}

	if pageSizeRaw, ok := args["page_size"]; ok {
		if pageSize, ok := pageSizeRaw.(float64); ok {
			options.PageSize = int(pageSize)
		}
	}

	results, err := s.client.GetBlockChildren(ctx, blockID, options)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to get Notion block children: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(results)
}

// handleNotionAppendBlockChildren handles the notion_append_block_children tool
func (s *Service) handleNotionAppendBlockChildren(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	blockID, ok := args["block_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Block ID is required"), nil
	}

	childrenRaw, ok := args["children"].([]interface{})
	if !ok {
		return mcpserver.CreateToolError("Children are required"), nil
	}

	// Convert children to Block objects
	var children []Block
	for _, child := range childrenRaw {
		// Convert to JSON and back to handle different block formats
		childBytes, err := json.Marshal(child)
		if err != nil {
			return mcpserver.CreateToolError(fmt.Sprintf("Failed to marshal block: %v", err)), nil
		}

		var block Block
		if err := json.Unmarshal(childBytes, &block); err != nil {
			return mcpserver.CreateToolError(fmt.Sprintf("Failed to unmarshal block: %v", err)), nil
		}

		children = append(children, block)
	}

	results, err := s.client.AppendBlockChildren(ctx, blockID, children)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to append Notion block children: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(results)
}
