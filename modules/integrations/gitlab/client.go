package gitlab

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/theshop/ai/internal/domain"
	"github.com/theshop/ai/modules/config"

	gitlab "github.com/xanzy/go-gitlab"
)

// Client GitLab API와의 상호작용을 위한 인터페이스
type Client interface {
	GetProject(ctx context.Context, projectID interface{}) (*domain.GitLabProject, error)
	ListProjects(ctx context.Context, search string, page, perPage int) ([]domain.GitLabProject, error)
	GetMergeRequest(ctx context.Context, projectID interface{}, mergeRequestIID int) (*domain.GitLabMergeRequest, error)
	CreateMergeRequest(ctx context.Context, projectID interface{}, opts MergeRequestOptions) (*domain.GitLabMergeRequest, error)
	CreateFile(ctx context.Context, projectID interface{}, filePath, content, commitMessage, branch string) error
	GetFileContent(ctx context.Context, projectID interface{}, filePath, ref string) (string, error)
}

// MergeRequestOptions MR 생성 옵션
type MergeRequestOptions struct {
	SourceBranch       string
	TargetBranch       string
	Title              string
	Description        string
	RemoveSourceBranch bool
	Squash             bool
}

// GitLabClient GitLab API 클라이언트 구현체
type GitLabClient struct {
	client *gitlab.Client
}

// NewClient 새 GitLab 클라이언트 생성
func NewClient(cfg *config.GitLabConfig) (Client, error) {
	client, err := gitlab.NewClient(cfg.Token, gitlab.WithBaseURL(cfg.URL))
	if err != nil {
		return nil, fmt.Errorf("GitLab 클라이언트 생성 오류: %w", err)
	}

	return &GitLabClient{
		client: client,
	}, nil
}

// GetProject GitLab 프로젝트 조회
func (c *GitLabClient) GetProject(ctx context.Context, projectID interface{}) (*domain.GitLabProject, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	project, _, err := c.client.Projects.GetProject(pid, nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("프로젝트 조회 오류: %w", err)
	}

	return &domain.GitLabProject{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		WebURL:      project.WebURL,
		CreatedAt:   *project.CreatedAt,
	}, nil
}

// ListProjects GitLab 프로젝트 목록 조회
func (c *GitLabClient) ListProjects(ctx context.Context, search string, page, perPage int) ([]domain.GitLabProject, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 100 {
		perPage = 20
	}

	opt := &gitlab.ListProjectsOptions{
		Search: gitlab.String(search),
		ListOptions: gitlab.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}

	projects, _, err := c.client.Projects.ListProjects(opt, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("프로젝트 목록 조회 오류: %w", err)
	}

	var result []domain.GitLabProject
	for _, project := range projects {
		result = append(result, domain.GitLabProject{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			WebURL:      project.WebURL,
			CreatedAt:   *project.CreatedAt,
		})
	}

	return result, nil
}

// GetMergeRequest GitLab 머지 리퀘스트 조회
func (c *GitLabClient) GetMergeRequest(ctx context.Context, projectID interface{}, mergeRequestIID int) (*domain.GitLabMergeRequest, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	mr, _, err := c.client.MergeRequests.GetMergeRequest(pid, mergeRequestIID, nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("머지 리퀘스트 조회 오류: %w", err)
	}

	createdAt := time.Now()
	if mr.CreatedAt != nil {
		createdAt = *mr.CreatedAt
	}

	return &domain.GitLabMergeRequest{
		ID:          mr.ID,
		IID:         mr.IID,
		Title:       mr.Title,
		Description: mr.Description,
		WebURL:      mr.WebURL,
		State:       mr.State,
		CreatedAt:   createdAt,
	}, nil
}

// CreateMergeRequest GitLab 머지 리퀘스트 생성
func (c *GitLabClient) CreateMergeRequest(ctx context.Context, projectID interface{}, opts MergeRequestOptions) (*domain.GitLabMergeRequest, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	createOpts := &gitlab.CreateMergeRequestOptions{
		Title:              gitlab.String(opts.Title),
		Description:        gitlab.String(opts.Description),
		SourceBranch:       gitlab.String(opts.SourceBranch),
		TargetBranch:       gitlab.String(opts.TargetBranch),
		RemoveSourceBranch: gitlab.Bool(opts.RemoveSourceBranch),
		Squash:             gitlab.Bool(opts.Squash),
	}

	mr, _, err := c.client.MergeRequests.CreateMergeRequest(pid, createOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("머지 리퀘스트 생성 오류: %w", err)
	}

	createdAt := time.Now()
	if mr.CreatedAt != nil {
		createdAt = *mr.CreatedAt
	}

	return &domain.GitLabMergeRequest{
		ID:          mr.ID,
		IID:         mr.IID,
		Title:       mr.Title,
		Description: mr.Description,
		WebURL:      mr.WebURL,
		State:       mr.State,
		CreatedAt:   createdAt,
	}, nil
}

// CreateFile GitLab 파일 생성
func (c *GitLabClient) CreateFile(ctx context.Context, projectID interface{}, filePath, content, commitMessage, branch string) error {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return err
	}

	fileOptions := &gitlab.CreateFileOptions{
		Branch:        gitlab.String(branch),
		Content:       gitlab.String(content),
		CommitMessage: gitlab.String(commitMessage),
	}

	_, _, err = c.client.RepositoryFiles.CreateFile(pid, filePath, fileOptions, gitlab.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("파일 생성 오류: %w", err)
	}

	return nil
}

// GetFileContent GitLab 파일 내용 조회
func (c *GitLabClient) GetFileContent(ctx context.Context, projectID interface{}, filePath, ref string) (string, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return "", err
	}

	if ref == "" {
		ref = "main"
	}

	fileOptions := &gitlab.GetFileOptions{
		Ref: gitlab.String(ref),
	}

	file, _, err := c.client.RepositoryFiles.GetFile(pid, filePath, fileOptions, gitlab.WithContext(ctx))
	if err != nil {
		return "", fmt.Errorf("파일 조회 오류: %w", err)
	}

	return file.Content, nil
}

// parseProjectID ID 또는 경로를 프로젝트 ID로 변환
func parseProjectID(projectID interface{}) (string, error) {
	switch v := projectID.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		// 이미 문자열이면 그대로 사용
		return v, nil
	case *string:
		if v == nil {
			return "", fmt.Errorf("프로젝트 ID가 nil입니다")
		}
		return *v, nil
	case *int:
		if v == nil {
			return "", fmt.Errorf("프로젝트 ID가 nil입니다")
		}
		return strconv.Itoa(*v), nil
	default:
		return "", fmt.Errorf("지원하지 않는 프로젝트 ID 타입: %T", projectID)
	}
}
