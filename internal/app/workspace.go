package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/theshop/ai/internal/domain"

	"github.com/google/uuid"
)

// WorkspaceManager 워크스페이스 관리 인터페이스
type WorkspaceManager interface {
	CreateWorkspace(ctx context.Context, name, userID, description string) (*domain.Workspace, error)
	GetWorkspace(ctx context.Context, id string) (*domain.Workspace, error)
	ListWorkspaces(ctx context.Context, userID string) ([]*domain.Workspace, error)
	UpdateWorkspace(ctx context.Context, workspace *domain.Workspace) error
	DeleteWorkspace(ctx context.Context, id string) error
}

// InMemoryWorkspaceManager 메모리 기반 워크스페이스 관리자 구현
// 실제 애플리케이션에서는 데이터베이스 기반으로 교체 가능
type InMemoryWorkspaceManager struct {
	workspaces map[string]*domain.Workspace
	mutex      sync.RWMutex
}

// NewInMemoryWorkspaceManager 새 인메모리 워크스페이스 관리자 생성
func NewInMemoryWorkspaceManager() *InMemoryWorkspaceManager {
	return &InMemoryWorkspaceManager{
		workspaces: make(map[string]*domain.Workspace),
	}
}

// CreateWorkspace 새 워크스페이스 생성
func (m *InMemoryWorkspaceManager) CreateWorkspace(ctx context.Context, name, userID, description string) (*domain.Workspace, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	id := uuid.New().String()
	now := time.Now()

	workspace := &domain.Workspace{
		ID:          id,
		Name:        name,
		UserID:      userID,
		Description: description,
		Settings: domain.WorkspaceSettings{
			MCPServers: make(map[string]string),
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	m.workspaces[id] = workspace
	return workspace, nil
}

// GetWorkspace 워크스페이스 조회
func (m *InMemoryWorkspaceManager) GetWorkspace(ctx context.Context, id string) (*domain.Workspace, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	workspace, exists := m.workspaces[id]
	if !exists {
		return nil, fmt.Errorf("워크스페이스를 찾을 수 없음: %s", id)
	}

	return workspace, nil
}

// ListWorkspaces 사용자의 워크스페이스 목록 조회
func (m *InMemoryWorkspaceManager) ListWorkspaces(ctx context.Context, userID string) ([]*domain.Workspace, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var workspaces []*domain.Workspace
	for _, ws := range m.workspaces {
		if ws.UserID == userID {
			workspaces = append(workspaces, ws)
		}
	}

	return workspaces, nil
}

// UpdateWorkspace 워크스페이스 업데이트
func (m *InMemoryWorkspaceManager) UpdateWorkspace(ctx context.Context, workspace *domain.Workspace) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, exists := m.workspaces[workspace.ID]
	if !exists {
		return errors.New("워크스페이스를 찾을 수 없음")
	}

	workspace.UpdatedAt = time.Now()
	m.workspaces[workspace.ID] = workspace
	return nil
}

// DeleteWorkspace 워크스페이스 삭제
func (m *InMemoryWorkspaceManager) DeleteWorkspace(ctx context.Context, id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, exists := m.workspaces[id]
	if !exists {
		return errors.New("워크스페이스를 찾을 수 없음")
	}

	delete(m.workspaces, id)
	return nil
}
