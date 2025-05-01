// Path: modules/integrations/gitlab/service.go
package gitlab

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/theshop/ai/pkg/mcpserver"
	"github.com/xanzy/go-gitlab"
)

// Service represents a service for interacting with GitLab
type Service struct {
	client Client
}

// NewService creates a new GitLab service
func NewService(client Client) *Service {
	return &Service{
		client: client,
	}
}

// RegisterTools registers GitLab tools with the MCP server
func (s *Service) RegisterTools(mcpServer *mcpserver.MCPServer) error {
	// Register GitLab list projects tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "gitlab_list_projects",
			Description: "List GitLab projects",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"search": map[string]interface{}{
						"type":        "string",
						"description": "Search term to filter projects",
					},
					"page": map[string]interface{}{
						"type":        "integer",
						"description": "Page number (1-based)",
					},
					"per_page": map[string]interface{}{
						"type":        "integer",
						"description": "Number of items per page (max 100)",
					},
				},
			},
		},
		s.handleGitLabListProjects,
	); err != nil {
		return err
	}

	// Register GitLab get project tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "gitlab_get_project",
			Description: "Get a GitLab project by ID or path",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID or path (e.g., 'group/project')",
					},
				},
				"required": []string{"project_id"},
			},
		},
		s.handleGitLabGetProject,
	); err != nil {
		return err
	}

	// Register GitLab get file tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "gitlab_get_file",
			Description: "Get a file from a GitLab repository",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID or path",
					},
					"file_path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the file in the repository",
					},
					"ref": map[string]interface{}{
						"type":        "string",
						"description": "Branch, tag, or commit (defaults to main/master)",
					},
				},
				"required": []string{"project_id", "file_path"},
			},
		},
		s.handleGitLabGetFile,
	); err != nil {
		return err
	}

	// Register GitLab list merge requests tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "gitlab_list_merge_requests",
			Description: "List merge requests in a GitLab project",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID or path",
					},
					"state": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"opened", "closed", "locked", "merged", "all"},
						"description": "State of merge requests to list",
					},
					"scope": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"created_by_me", "assigned_to_me", "all"},
						"description": "Scope of merge requests to list",
					},
				},
				"required": []string{"project_id"},
			},
		},
		s.handleGitLabListMergeRequests,
	); err != nil {
		return err
	}

	// Register GitLab create merge request tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "gitlab_create_merge_request",
			Description: "Create a new merge request in a GitLab project",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID or path",
					},
					"source_branch": map[string]interface{}{
						"type":        "string",
						"description": "Source branch name",
					},
					"target_branch": map[string]interface{}{
						"type":        "string",
						"description": "Target branch name",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Title of the merge request",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Description of the merge request",
					},
					"remove_source_branch": map[string]interface{}{
						"type":        "boolean",
						"description": "Remove source branch after merge",
					},
				},
				"required": []string{"project_id", "source_branch", "target_branch", "title"},
			},
		},
		s.handleGitLabCreateMergeRequest,
	); err != nil {
		return err
	}

	// Register GitLab create file tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "gitlab_create_file",
			Description: "Create a new file in a GitLab repository",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID or path",
					},
					"file_path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the file in the repository",
					},
					"branch": map[string]interface{}{
						"type":        "string",
						"description": "Branch name",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Content of the file",
					},
					"commit_message": map[string]interface{}{
						"type":        "string",
						"description": "Commit message",
					},
				},
				"required": []string{"project_id", "file_path", "branch", "content", "commit_message"},
			},
		},
		s.handleGitLabCreateFile,
	); err != nil {
		return err
	}

	// Register GitLab analyze code tool
	if err := mcpServer.RegisterTool(
		mcp.Tool{
			Name:        "gitlab_analyze_code",
			Description: "Analyze code from a GitLab repository to provide context-based recommendations",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID or path",
					},
					"file_path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the file to analyze",
					},
					"context": map[string]interface{}{
						"type":        "string",
						"description": "Context for the analysis (e.g., 'variable_names', 'query_structure', etc.)",
					},
					"additional_files": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "Additional files to include in the analysis",
					},
				},
				"required": []string{"project_id", "file_path", "context"},
			},
		},
		s.handleGitLabAnalyzeCode,
	); err != nil {
		return err
	}

	return nil
}

// handleGitLabListProjects handles the gitlab_list_projects tool
func (s *Service) handleGitLabListProjects(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	var search string
	if searchArg, ok := args["search"].(string); ok {
		search = searchArg
	}

	page := 1
	if pageArg, ok := args["page"].(float64); ok {
		page = int(pageArg)
	}

	perPage := 20
	if perPageArg, ok := args["per_page"].(float64); ok {
		perPage = int(perPageArg)
	}

	projects, err := s.client.ListProjects(ctx, search, page, perPage)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to list GitLab projects: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(projects)
}

// handleGitLabGetProject handles the gitlab_get_project tool
func (s *Service) handleGitLabGetProject(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	projectID, ok := args["project_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Project ID is required"), nil
	}

	project, err := s.client.GetProject(ctx, projectID)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to get GitLab project: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(project)
}

// handleGitLabGetFile handles the gitlab_get_file tool
func (s *Service) handleGitLabGetFile(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	projectID, ok := args["project_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Project ID is required"), nil
	}

	filePath, ok := args["file_path"].(string)
	if !ok {
		return mcpserver.CreateToolError("File path is required"), nil
	}

	var ref string
	if refArg, ok := args["ref"].(string); ok {
		ref = refArg
	}

	content, err := s.client.GetFileContent(ctx, projectID, filePath, ref)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to get file content: %v", err)), nil
	}

	// Detect file type for better formatting
	fileExt := getFileExtension(filePath)
	fileType := determineFileType(fileExt)

	result := map[string]interface{}{
		"project_id": projectID,
		"file_path":  filePath,
		"ref":        ref,
		"content":    content,
		"file_type":  fileType,
	}

	return mcpserver.CreateToolResultJSON(result)
}

// handleGitLabListMergeRequests handles the gitlab_list_merge_requests tool
func (s *Service) handleGitLabListMergeRequests(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	projectID, ok := args["project_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Project ID is required"), nil
	}

	state := "opened"
	if stateArg, ok := args["state"].(string); ok {
		state = stateArg
	}

	// In a real implementation, this would use options from the arguments
	// to filter the merge requests
	listOptions := &gitlab.ListProjectMergeRequestsOptions{
		State: &state,
	}

	mergeRequests, err := s.client.ListMergeRequests(ctx, projectID, listOptions)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to list merge requests: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(mergeRequests)
}

// handleGitLabCreateMergeRequest handles the gitlab_create_merge_request tool
func (s *Service) handleGitLabCreateMergeRequest(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	projectID, ok := args["project_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Project ID is required"), nil
	}

	sourceBranch, ok := args["source_branch"].(string)
	if !ok {
		return mcpserver.CreateToolError("Source branch is required"), nil
	}

	targetBranch, ok := args["target_branch"].(string)
	if !ok {
		return mcpserver.CreateToolError("Target branch is required"), nil
	}

	title, ok := args["title"].(string)
	if !ok {
		return mcpserver.CreateToolError("Title is required"), nil
	}

	description := ""
	if descArg, ok := args["description"].(string); ok {
		description = descArg
	}

	removeSourceBranch := false
	if removeArg, ok := args["remove_source_branch"].(bool); ok {
		removeSourceBranch = removeArg
	}

	opts := MergeRequestOptions{
		SourceBranch:       sourceBranch,
		TargetBranch:       targetBranch,
		Title:              title,
		Description:        description,
		RemoveSourceBranch: removeSourceBranch,
	}

	mr, err := s.client.CreateMergeRequest(ctx, projectID, opts)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to create merge request: %v", err)), nil
	}

	return mcpserver.CreateToolResultJSON(mr)
}

// handleGitLabCreateFile handles the gitlab_create_file tool
func (s *Service) handleGitLabCreateFile(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	projectID, ok := args["project_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Project ID is required"), nil
	}

	filePath, ok := args["file_path"].(string)
	if !ok {
		return mcpserver.CreateToolError("File path is required"), nil
	}

	branch, ok := args["branch"].(string)
	if !ok {
		return mcpserver.CreateToolError("Branch is required"), nil
	}

	content, ok := args["content"].(string)
	if !ok {
		return mcpserver.CreateToolError("Content is required"), nil
	}

	commitMessage, ok := args["commit_message"].(string)
	if !ok {
		return mcpserver.CreateToolError("Commit message is required"), nil
	}

	// Create the file in GitLab
	err := s.client.CreateFile(ctx, projectID, filePath, content, commitMessage, branch)
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to create file: %v", err)), nil
	}

	result := map[string]interface{}{
		"status":         "success",
		"project_id":     projectID,
		"file_path":      filePath,
		"branch":         branch,
		"commit_message": commitMessage,
	}

	return mcpserver.CreateToolResultJSON(result)
}

// handleGitLabAnalyzeCode handles the gitlab_analyze_code tool
func (s *Service) handleGitLabAnalyzeCode(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	projectID, ok := args["project_id"].(string)
	if !ok {
		return mcpserver.CreateToolError("Project ID is required"), nil
	}

	filePath, ok := args["file_path"].(string)
	if !ok {
		return mcpserver.CreateToolError("File path is required"), nil
	}

	context, ok := args["context"].(string)
	if !ok {
		return mcpserver.CreateToolError("Context is required"), nil
	}

	// Get the main file content
	mainContent, err := s.client.GetFileContent(ctx, projectID, filePath, "")
	if err != nil {
		return mcpserver.CreateToolError(fmt.Sprintf("Failed to get file content: %v", err)), nil
	}

	// Get additional files if specified
	additionalContents := make(map[string]string)
	if additionalFilesRaw, ok := args["additional_files"].([]interface{}); ok {
		for _, filePathRaw := range additionalFilesRaw {
			if addFilePath, ok := filePathRaw.(string); ok {
				content, err := s.client.GetFileContent(ctx, projectID, addFilePath, "")
				if err != nil {
					// Log but don't fail if additional file can't be retrieved
					fmt.Printf("Warning: Failed to get content for additional file %s: %v\n", addFilePath, err)
					continue
				}
				additionalContents[addFilePath] = content
			}
		}
	}

	// In a real implementation, this would analyze the code based on the context
	// and provide recommendations. For now, we'll return a mock result
	analysisType := ""
	switch context {
	case "variable_names":
		analysisType = "Variable naming conventions"
	case "query_structure":
		analysisType = "SQL query structure"
	case "code_style":
		analysisType = "Code style recommendations"
	default:
		analysisType = "General code analysis"
	}

	result := map[string]interface{}{
		"project_id":       projectID,
		"file_path":        filePath,
		"analysis_type":    analysisType,
		"file_type":        determineFileType(getFileExtension(filePath)),
		"additional_files": len(additionalContents),
		"recommendations":  generateMockRecommendations(context, filePath),
	}

	return mcpserver.CreateToolResultJSON(result)
}

// Helper functions

// getFileExtension extracts the file extension from a path
func getFileExtension(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
}

// determineFileType determines the file type based on extension
func determineFileType(ext string) string {
	switch strings.ToLower(ext) {
	case "go":
		return "Go"
	case "py":
		return "Python"
	case "js":
		return "JavaScript"
	case "ts":
		return "TypeScript"
	case "java":
		return "Java"
	case "c":
		return "C"
	case "cpp", "cc":
		return "C++"
	case "cs":
		return "C#"
	case "php":
		return "PHP"
	case "rb":
		return "Ruby"
	case "sql":
		return "SQL"
	case "html":
		return "HTML"
	case "css":
		return "CSS"
	case "md", "markdown":
		return "Markdown"
	case "json":
		return "JSON"
	case "yml", "yaml":
		return "YAML"
	case "xml":
		return "XML"
	default:
		return "Text"
	}
}

// generateMockRecommendations generates mock recommendations based on context
func generateMockRecommendations(contextType, filePath string) []string {
	switch contextType {
	case "variable_names":
		return []string{
			"Use camelCase for variable names",
			"Prefix interface names with 'I'",
			"Use descriptive names that reflect purpose",
		}
	case "query_structure":
		return []string{
			"Use prepared statements to prevent SQL injection",
			"Add indexes for frequently queried columns",
			"Split complex queries into simpler ones",
		}
	case "code_style":
		return []string{
			"Add comments for complex logic",
			"Keep functions small and focused",
			"Use consistent indentation",
		}
	default:
		return []string{
			"Implement error handling",
			"Add unit tests",
			"Consider breaking large files into modules",
		}
	}
}
