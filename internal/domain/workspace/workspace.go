package workspace

import (
	"context"
	"fmt"
	"time"
)

// Workspace는 사용자의 작업 컨텍스트를 나타내는 구조체입니다.
type Workspace struct {
	ID               string                 // 워크스페이스 고유 식별자
	Name             string                 // 워크스페이스 이름
	Description      string                 // 워크스페이스 설명
	Owner            string                 // 소유자 ID
	Created          time.Time              // 생성 시간
	LastAccessed     time.Time              // 마지막 접근 시간
	JiraProjects     []string               // 연결된 Jira 프로젝트
	GitLabProjects   []string               // 연결된 GitLab 프로젝트
	ConfluenceSpaces []string               // 연결된 Confluence 스페이스
	NotionPages      []string               // 연결된 Notion 페이지
	FilePaths        []string               // 연결된 파일 경로
	MCPServers       []string               // 연결할 MCP 서버 ID 목록
	LLMProvider      string                 // 사용할 LLM 제공자 (openai, claude, gemini 등)
	Settings         map[string]interface{} // 워크스페이스 관련 추가 설정
}

// WorkspaceService는 워크스페이스 관리 서비스 인터페이스입니다.
type WorkspaceService interface {
	Create(ctx context.Context, workspace *Workspace) error
	Get(ctx context.Context, id string) (*Workspace, error)
	Update(ctx context.Context, workspace *Workspace) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, ownerID string) ([]*Workspace, error)
}

// WorkspaceRepository는 워크스페이스 저장소 인터페이스입니다.
type WorkspaceRepository interface {
	Save(ctx context.Context, workspace *Workspace) error
	FindByID(ctx context.Context, id string) (*Workspace, error)
	Update(ctx context.Context, workspace *Workspace) error
	Delete(ctx context.Context, id string) error
	FindByOwner(ctx context.Context, ownerID string) ([]*Workspace, error)
}

// WorkspaceContextKey는 컨텍스트에서 워크스페이스 정보를 가져오기 위한 키입니다.
type WorkspaceContextKey struct{}

// FromContext는 context에서 Workspace 값을 조회합니다.
func FromContext(ctx context.Context) (*Workspace, error) {
	workspace, ok := ctx.Value(WorkspaceContextKey{}).(*Workspace)
	if !ok {
		return nil, fmt.Errorf("워크스페이스 컨텍스트를 찾을 수 없습니다")
	}
	return workspace, nil
}

// WithContext는 context에 Workspace 값을 추가합니다.
func WithContext(ctx context.Context, workspace *Workspace) context.Context {
	return context.WithValue(ctx, WorkspaceContextKey{}, workspace)
}
