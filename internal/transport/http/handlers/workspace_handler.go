package handlers

import (
	"net/http"

	"github.com/theshop/ai/internal/app"
	"github.com/theshop/ai/internal/domain"

	"github.com/gin-gonic/gin"
)

// WorkspaceHandler 워크스페이스 API 핸들러
type WorkspaceHandler struct {
	workspaceMgr app.WorkspaceManager
}

// NewWorkspaceHandler 새 워크스페이스 핸들러 생성
func NewWorkspaceHandler(workspaceMgr app.WorkspaceManager) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceMgr: workspaceMgr,
	}
}

// CreateWorkspaceRequest 워크스페이스 생성 요청
type CreateWorkspaceRequest struct {
	Name        string `json:"name" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	Description string `json:"description"`
}

// UpdateWorkspaceRequest 워크스페이스 업데이트 요청
type UpdateWorkspaceRequest struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Settings    domain.WorkspaceSettings `json:"settings"`
}

// CreateWorkspace 워크스페이스 생성 핸들러
func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	var req CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "유효하지 않은 요청 데이터: " + err.Error(),
		})
		return
	}

	workspace, err := h.workspaceMgr.CreateWorkspace(c, req.Name, req.UserID, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "워크스페이스 생성 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, workspace)
}

// GetWorkspace 워크스페이스 조회 핸들러
func (h *WorkspaceHandler) GetWorkspace(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "워크스페이스 ID가 필요합니다",
		})
		return
	}

	workspace, err := h.workspaceMgr.GetWorkspace(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "워크스페이스를 찾을 수 없음: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

// ListWorkspaces 워크스페이스 목록 조회 핸들러
func (h *WorkspaceHandler) ListWorkspaces(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "사용자 ID가 필요합니다",
		})
		return
	}

	workspaces, err := h.workspaceMgr.ListWorkspaces(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "워크스페이스 목록 조회 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}

// UpdateWorkspace 워크스페이스 업데이트 핸들러
func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "워크스페이스 ID가 필요합니다",
		})
		return
	}

	var req UpdateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "유효하지 않은 요청 데이터: " + err.Error(),
		})
		return
	}

	workspace, err := h.workspaceMgr.GetWorkspace(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "워크스페이스를 찾을 수 없음: " + err.Error(),
		})
		return
	}

	// 업데이트 가능한 필드만 수정
	if req.Name != "" {
		workspace.Name = req.Name
	}
	if req.Description != "" {
		workspace.Description = req.Description
	}

	// 설정 병합
	if req.Settings.JiraProjectKey != "" {
		workspace.Settings.JiraProjectKey = req.Settings.JiraProjectKey
	}
	if req.Settings.GitLabProject != "" {
		workspace.Settings.GitLabProject = req.Settings.GitLabProject
	}
	if req.Settings.ConfluenceSpace != "" {
		workspace.Settings.ConfluenceSpace = req.Settings.ConfluenceSpace
	}
	if req.Settings.NotionDatabase != "" {
		workspace.Settings.NotionDatabase = req.Settings.NotionDatabase
	}
	if req.Settings.DefaultLLM != "" {
		workspace.Settings.DefaultLLM = req.Settings.DefaultLLM
	}

	// MCP 서버 설정 병합
	if req.Settings.MCPServers != nil {
		if workspace.Settings.MCPServers == nil {
			workspace.Settings.MCPServers = make(map[string]string)
		}
		for k, v := range req.Settings.MCPServers {
			workspace.Settings.MCPServers[k] = v
		}
	}

	if err := h.workspaceMgr.UpdateWorkspace(c, workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "워크스페이스 업데이트 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

// DeleteWorkspace 워크스페이스 삭제 핸들러
func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "워크스페이스 ID가 필요합니다",
		})
		return
	}

	if err := h.workspaceMgr.DeleteWorkspace(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "워크스페이스 삭제 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "워크스페이스가 삭제되었습니다",
	})
}
