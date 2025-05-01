package domain

import "time"

// Workspace 사용자 워크스페이스 모델
type Workspace struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	UserID      string            `json:"user_id"`
	Description string            `json:"description"`
	Settings    WorkspaceSettings `json:"settings"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// WorkspaceSettings 워크스페이스 설정
type WorkspaceSettings struct {
	JiraProjectKey  string            `json:"jira_project_key"`
	GitLabProject   string            `json:"gitlab_project"`
	ConfluenceSpace string            `json:"confluence_space"`
	NotionDatabase  string            `json:"notion_database"`
	DefaultLLM      string            `json:"default_llm"`
	MCPServers      map[string]string `json:"mcp_servers"`
}

// User 사용자 모델
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// JiraIssue Jira 이슈 모델
type JiraIssue struct {
	Key         string                 `json:"key"`
	Summary     string                 `json:"summary"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Fields      map[string]interface{} `json:"fields"`
}

// GitLabProject GitLab 프로젝트 모델
type GitLabProject struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WebURL      string    `json:"web_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// GitLabMergeRequest GitLab 머지 리퀘스트 모델
type GitLabMergeRequest struct {
	ID          int       `json:"id"`
	IID         int       `json:"iid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	WebURL      string    `json:"web_url"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
}

// NotionPage Notion 페이지 모델
type NotionPage struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	URL         string                 `json:"url"`
	Properties  map[string]interface{} `json:"properties"`
	LastUpdated time.Time              `json:"last_updated"`
}

// MCPToolResult MCP 도구 호출 결과
type MCPToolResult struct {
	Content string                 `json:"content"`
	Data    map[string]interface{} `json:"data"`
	Error   string                 `json:"error,omitempty"`
}

// ProjectOverview 프로젝트 개요
type ProjectOverview struct {
	WorkspaceID   string                 `json:"workspace_id"`
	GitLabProject GitLabProject          `json:"gitlab_project"`
	JiraIssues    []JiraIssue            `json:"jira_issues"`
	RecentCommits []GitLabCommit         `json:"recent_commits,omitempty"`
	Analytics     map[string]interface{} `json:"analytics,omitempty"`
}

// GitLabCommit GitLab 커밋 정보
type GitLabCommit struct {
	ID         string    `json:"id"`
	ShortID    string    `json:"short_id"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
}
