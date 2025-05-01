//go:build wireinject
// +build wireinject

package core

import (
	"context"
	"fmt"

	"github.com/theshop/ai/internal/app"
	"github.com/theshop/ai/modules/config"
	"github.com/theshop/ai/modules/integrations/gitlab"
	"github.com/theshop/ai/modules/integrations/jira"
	"github.com/theshop/ai/modules/integrations/llm"
	"github.com/theshop/ai/pkg/mcpclient"

	"github.com/google/wire"
)

// 애플리케이션 구조체
type Application struct {
	Config             *config.Config
	JiraClient         jira.Client
	GitLabClient       gitlab.Client
	LLMService         llm.Service
	WorkspaceManager   app.WorkspaceManager
	IntegrationService *app.IntegrationService
	MCPClients         map[string]mcpclient.MCPClient
}

// ProvideWorkspaceManager 워크스페이스 관리자 생성
func ProvideWorkspaceManager() app.WorkspaceManager {
	return app.NewInMemoryWorkspaceManager()
}

// ProvideMCPClients MCP 클라이언트 맵 생성
func ProvideMCPClients(cfg *config.Config) (map[string]mcpclient.MCPClient, error) {
	clients := make(map[string]mcpclient.MCPClient)

	// 설정에서 MCP 서버 정보를 기반으로 클라이언트 생성
	for name, serverCfg := range cfg.MCPServers {
		clientCfg := mcpclient.MCPClientConfig{
			Type:    serverCfg.Type,
			URL:     serverCfg.URL,
			Command: serverCfg.Command,
			Args:    serverCfg.Args,
			Env:     serverCfg.Env,
		}

		client, err := mcpclient.NewMCPClient(name, clientCfg)
		if err != nil {
			return nil, err
		}

		// 클라이언트 초기화
		if err := client.Initialize(context.Background()); err != nil {
			return nil, err
		}

		clients[name] = client
	}

	return clients, nil
}

// 의존성 주입 세트
var ApplicationSet = wire.NewSet(
	config.NewConfig,
	jira.NewClient,
	ProvideMCPClients,
	wire.Bind(new(llm.Service), new(*llm.MCPLLMService)),

	// GitLab 클라이언트 제공
	wire.Bind(new(gitlab.Client), new(*gitlab.GitLabClient)),
	ProvideGitLabClient,

	// LLM 서비스 제공
	ProvideLLMService,

	// 워크스페이스 관리자 제공
	ProvideWorkspaceManager,

	// 통합 서비스 제공
	app.NewIntegrationService,

	// 애플리케이션 제공
	ProvideApplication,
)

// ProvideGitLabClient 깃랩 클라이언트 제공 함수
func ProvideGitLabClient(cfg *config.Config) (*gitlab.GitLabClient, error) {
	client, err := gitlab.NewClient(cfg.GitLab)
	return client.(*gitlab.GitLabClient), err
}

// ProvideLLMService LLM 서비스 제공 함수
func ProvideLLMService(mcpClients map[string]mcpclient.MCPClient, cfg *config.Config) (*llm.MCPLLMService, error) {
	// 설정된 LLM 제공자에 해당하는 MCP 클라이언트 선택
	var mcpClient mcpclient.MCPClient
	var exists bool

	mcpClient, exists = mcpClients[cfg.LLM.Provider]
	if !exists {
		// 기본 제공자가 없으면 첫 번째 MCP 클라이언트 사용
		for _, client := range mcpClients {
			mcpClient = client
			break
		}
	}

	if mcpClient == nil {
		return nil, fmt.Errorf("LLM 서비스용 MCP 클라이언트를 찾을 수 없음")
	}

	service := llm.NewMCPLLMService(mcpClient, &cfg.LLM)
	return service.(*llm.MCPLLMService), nil
}

// ProvideApplication 애플리케이션 제공 함수
func ProvideApplication(
	cfg *config.Config,
	jiraClient jira.Client,
	gitlabClient gitlab.Client,
	llmService llm.Service,
	workspaceMgr app.WorkspaceManager,
	integrationSvc *app.IntegrationService,
	mcpClients map[string]mcpclient.MCPClient,
) *Application {
	return &Application{
		Config:             cfg,
		JiraClient:         jiraClient,
		GitLabClient:       gitlabClient,
		LLMService:         llmService,
		WorkspaceManager:   workspaceMgr,
		IntegrationService: integrationSvc,
		MCPClients:         mcpClients,
	}
}

// InitializeApp 애플리케이션 초기화 함수 (Wire가 구현)
func InitializeApp() (*Application, error) {
	wire.Build(ApplicationSet)
	return nil, nil
}
