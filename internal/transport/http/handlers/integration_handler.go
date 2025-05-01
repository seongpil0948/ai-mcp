package handlers

import (
	"net/http"

	"github.com/theshop/ai/internal/app"
	"github.com/theshop/ai/modules/integrations/gitlab"
	"github.com/theshop/ai/modules/integrations/jira"

	"github.com/gin-gonic/gin"
)

// IntegrationHandler 통합 API 핸들러
type IntegrationHandler struct {
	workspaceMgr   app.WorkspaceManager
	jiraClient     jira.Client
	gitlabClient   gitlab.Client
	integrationSvc *app.IntegrationService
}

// NewIntegrationHandler 새 통합 핸들러 생성
func NewIntegrationHandler(
	workspaceMgr app.WorkspaceManager,
	jiraClient jira.Client,
	gitlabClient gitlab.Client,
	integrationSvc *app.IntegrationService,
) *IntegrationHandler {
	return &IntegrationHandler{
		workspaceMgr:   workspaceMgr,
		jiraClient:     jiraClient,
		gitlabClient:   gitlabClient,
		integrationSvc: integrationSvc,
	}
}

// CreateJiraIssueRequest Jira 이슈 생성 요청
type CreateJiraIssueRequest struct {
	ProjectKey  string `json:"project_key"`
	Summary     string `json:"summary" binding:"required"`
	Description string `json:"description"`
}

// CreateJiraIssue Jira 이슈 생성 핸들러
func (h *IntegrationHandler) CreateJiraIssue(c *gin.Context) {
	workspaceID := c.Param("id")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "워크스페이스 ID가 필요합니다",
		})
		return
	}

	workspace, err := h.workspaceMgr.GetWorkspace(c, workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "워크스페이스를 찾을 수 없음: " + err.Error(),
		})
		return
	}

	var req CreateJiraIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "유효하지 않은 요청 데이터: " + err.Error(),
		})
		return
	}

	projectKey := req.ProjectKey
	if projectKey == "" {
		projectKey = workspace.Settings.JiraProjectKey
	}

	issue, err := h.jiraClient.CreateIssue(c, projectKey, req.Summary, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "이슈 생성 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, issue)
}

// SearchJiraIssuesRequest Jira 이슈 검색 요청
type SearchJiraIssuesRequest struct {
	JQL        string `json:"jql" binding:"required"`
	MaxResults int    `json:"max_results"`
}

// SearchJiraIssues Jira 이슈 검색 핸들러
func (h *IntegrationHandler) SearchJiraIssues(c *gin.Context) {
	workspaceID := c.Param("id")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "워크스페이스 ID가 필요합니다",
		})
		return
	}

	_, err := h.workspaceMgr.GetWorkspace(c, workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "워크스페이스를 찾을 수 없음: " + err.Error(),
		})
		return
	}

	var req SearchJiraIssuesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "유효하지 않은 요청 데이터: " + err.Error(),
		})
		return
	}

	maxResults := req.MaxResults
	if maxResults <= 0 {
		maxResults = 10
	}

	issues, err := h.jiraClient.SearchIssues(c, req.JQL, maxResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "이슈 검색 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, issues)
}

// ListGitLabProjectsRequest GitLab 프로젝트 목록 요청
type ListGitLabProjectsRequest struct {
	Search  string `json:"search"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

// ListGitLabProjects GitLab 프로젝트 목록 핸들러
func (h *IntegrationHandler) ListGitLabProjects(c *gin.Context) {
	workspaceID := c.Param("id")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "워크스페이스 ID가 필요합니다",
		})
		return
	}

	_, err := h.workspaceMgr.GetWorkspace(c, workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "워크스페이스를 찾을 수 없음: " + err.Error(),
		})
		return
	}

	var req ListGitLabProjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// JSON 바인딩 실패 시 기본값 사용
		req = ListGitLabProjectsRequest{
			Page:    1,
			PerPage: 20,
		}
	}

	projects, err := h.gitlabClient.ListProjects(c, req.Search, req.Page, req.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "프로젝트 목록 조회 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// CreateMRFromIssueRequest 이슈에서 MR 생성 요청
type CreateMRFromIssueRequest struct {
	ProjectID          interface{} `json:"project_id" binding:"required"`
	SourceBranch       string      `json:"source_branch" binding:"required"`
	TargetBranch       string      `json:"target_branch" binding:"required"`
	IssueKey           string      `json:"issue_key" binding:"required"`
	Description        string      `json:"description"`
	RemoveSourceBranch bool        `json:"remove_source_branch"`
}

// CreateMRFromIssue 이슈에서 MR 생성 핸들러
func (h *IntegrationHandler) CreateMRFromIssue(c *gin.Context) {
	workspaceID := c.Param("id")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "워크스페이스 ID가 필요합니다",
		})
		return
	}

	_, err := h.workspaceMgr.GetWorkspace(c, workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "워크스페이스를 찾을 수 없음: " + err.Error(),
		})
		return
	}

	var req CreateMRFromIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "유효하지 않은 요청 데이터: " + err.Error(),
		})
		return
	}

	// 이슈 정보 조회
	issue, err := h.jiraClient.GetIssue(c, req.IssueKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "이슈 조회 실패: " + err.Error(),
		})
		return
	}

	// MR 설명 생성
	description := req.Description
	if description == "" {
		description = issue.Description
	}
	description = "Resolves " + req.IssueKey + "\n\n" + description

	// MR 생성
	mergeRequestOpts := gitlab.MergeRequestOptions{
		SourceBranch:       req.SourceBranch,
		TargetBranch:       req.TargetBranch,
		Title:              issue.Key + ": " + issue.Summary,
		Description:        description,
		RemoveSourceBranch: req.RemoveSourceBranch,
	}

	mr, err := h.gitlabClient.CreateMergeRequest(c, req.ProjectID, mergeRequestOpts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "MR 생성 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, mr)
}
