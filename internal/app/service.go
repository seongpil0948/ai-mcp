package app

import (
	"context"
	"fmt"

	"github.com/theshop/ai/internal/domain"
	"github.com/theshop/ai/modules/integrations/gitlab"
	"github.com/theshop/ai/modules/integrations/jira"
	"github.com/theshop/ai/modules/integrations/llm"
	"github.com/theshop/ai/pkg/mcpclient"
)

// IntegrationService 통합 서비스
type IntegrationService struct {
	gitlabClient gitlab.Client // Use the defined interface type
	jiraClient   jira.Client   // Use the defined interface type
	llmService   llm.Service   // Use the defined interface type
	workspaceMgr WorkspaceManager
	mcpClients   map[string]mcpclient.MCPClient
}

// NewIntegrationService 새 통합 서비스 생성
func NewIntegrationService(
	glClient gitlab.Client,
	jClient jira.Client,
	llmSvc llm.Service,
	wsMgr WorkspaceManager,
	mcpMap map[string]mcpclient.MCPClient,
) *IntegrationService {
	// Assuming line 15 was inside this return statement
	return &IntegrationService{
		gitlabClient: glClient, // Ensure no trailing dots or missing commas
		jiraClient:   jClient,
		llmService:   llmSvc,
		workspaceMgr: wsMgr,
		mcpClients:   mcpMap, // Check line 15 was here or nearby
	}
}

// CreateJiraIssueFromMergeRequest GitLab MR에서 Jira 이슈 생성
func (s *IntegrationService) CreateJiraIssueFromMergeRequest(
	ctx context.Context,
	workspaceID string,
	gitlabProjectID interface{},
	mergeRequestIID int,
) (string, error) {
	// 워크스페이스 로드
	workspace, err := s.workspaceMgr.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return "", fmt.Errorf("워크스페이스 로드 오류: %w", err)
	}

	// GitLab MR 가져오기
	mr, err := s.gitlabClient.GetMergeRequest(ctx, gitlabProjectID, mergeRequestIID)
	if err != nil {
		return "", fmt.Errorf("MR 로드 오류: %w", err)
	}

	// LLM으로 MR 설명 요약
	summary, err := s.llmService.GenerateText(ctx,
		fmt.Sprintf("다음 GitLab MR을 기반으로 Jira 이슈 요약을 작성해주세요: %s", mr.Description),
		llm.GenerateOptions{MaxTokens: 100},
	)
	if err != nil {
		return "", fmt.Errorf("요약 생성 오류: %w", err)
	}

	// Jira 이슈 생성
	issue, err := s.jiraClient.CreateIssue(ctx,
		workspace.Settings.JiraProjectKey,
		fmt.Sprintf("MR #%d: %s", mergeRequestIID, mr.Title),
		fmt.Sprintf("GitLab MR에서 생성됨: %s\n\n%s\n\n원본 설명:\n%s",
			mr.WebURL, summary, mr.Description),
	)
	if err != nil {
		return "", fmt.Errorf("Jira 이슈 생성 오류: %w", err)
	}

	return issue.Key, nil
}

// GetMCPClient 워크스페이스에 맞는 MCP 클라이언트 반환
func (s *IntegrationService) GetMCPClient(ctx context.Context, workspaceID, mcpServerName string) (mcpclient.MCPClient, error) {
	// 전역 MCP 클라이언트 확인
	if client, exists := s.mcpClients[mcpServerName]; exists {
		return client, nil
	}

	// 워크스페이스별 MCP 서버 설정 확인
	workspace, err := s.workspaceMgr.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("워크스페이스 로드 오류: %w", err)
	}

	serverID, exists := workspace.Settings.MCPServers[mcpServerName]
	if !exists {
		return nil, fmt.Errorf("워크스페이스에 MCP 서버 설정이 없음: %s", mcpServerName)
	}

	// 워크스페이스별 MCP 서버 설정 사용
	// 실제 구현에서는 여기서 새 클라이언트를 생성하거나 캐시된 클라이언트를 반환
	if client, exists := s.mcpClients[serverID]; exists {
		return client, nil
	}

	return nil, fmt.Errorf("MCP 클라이언트를 찾을 수 없음: %s", serverID)
}

// GenerateMRDescription MR 설명 생성
func (s *IntegrationService) GenerateMRDescription(ctx context.Context, issueKey, changes string) (string, error) {
	// 이슈 조회
	issue, err := s.jiraClient.GetIssue(ctx, issueKey)
	if err != nil {
		return "", fmt.Errorf("이슈 조회 오류: %w", err)
	}

	// LLM으로 설명 생성
	prompt := fmt.Sprintf(
		"다음 Jira 이슈와 변경사항을 바탕으로 GitLab 머지 리퀘스트 설명을 작성해주세요:\n\n"+
			"이슈 키: %s\n제목: %s\n설명: %s\n\n변경사항:\n%s",
		issue.Key, issue.Summary, issue.Description, changes,
	)

	description, err := s.llmService.GenerateText(ctx, prompt, llm.GenerateOptions{
		System: "당신은 개발자로서 변경 사항에 대한 명확하고 간결한 설명을 작성하는 전문가입니다.",
	})
	if err != nil {
		return "", fmt.Errorf("설명 생성 오류: %w", err)
	}

	// 기본 템플릿 적용 (이슈 키 참조 포함)
	return fmt.Sprintf("Resolves %s\n\n%s", issue.Key, description), nil
}

// AnalyzeCodeWithLLM LLM으로 코드 분석
func (s *IntegrationService) AnalyzeCodeWithLLM(ctx context.Context, code, query string) (string, error) {
	prompt := fmt.Sprintf("다음 코드를 분석해주세요:\n\n```\n%s\n```\n\n질문: %s", code, query)

	result, err := s.llmService.GenerateText(ctx, prompt, llm.GenerateOptions{
		System: "당신은 코드 분석 전문가입니다. 코드를 면밀히 검토하고 명확하게 설명해주세요.",
	})
	if err != nil {
		return "", fmt.Errorf("코드 분석 오류: %w", err)
	}

	return result, nil
}

// GetProjectOverview 프로젝트 개요 생성
func (s *IntegrationService) GetProjectOverview(ctx context.Context, workspaceID string, gitlabProjectID interface{}) (*domain.ProjectOverview, error) {
	// 워크스페이스 로드
	workspace, err := s.workspaceMgr.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("워크스페이스 로드 오류: %w", err)
	}

	// GitLab 프로젝트 정보 가져오기
	project, err := s.gitlabClient.GetProject(ctx, gitlabProjectID)
	if err != nil {
		return nil, fmt.Errorf("프로젝트 로드 오류: %w", err)
	}

	// Jira 이슈 목록 가져오기
	jiraProjectKey := workspace.Settings.JiraProjectKey
	jql := fmt.Sprintf("project = %s ORDER BY updated DESC", jiraProjectKey)
	issues, err := s.jiraClient.SearchIssues(ctx, jql, 10)
	if err != nil {
		return nil, fmt.Errorf("이슈 검색 오류: %w", err)
	}

	// 프로젝트 개요 구성
	overview := &domain.ProjectOverview{
		WorkspaceID:   workspaceID,
		GitLabProject: *project,
		JiraIssues:    issues,
		// 기타 필요한 정보 추가
	}

	return overview, nil
}
