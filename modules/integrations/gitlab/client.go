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
	GetMergeRequest(ctx context.Context, projectID interface{}, mergeRequestID int) (*domain.GitLabMergeRequest, error)
	ListMergeRequests(ctx context.Context, projectID interface{}, opts *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	CreateMergeRequest(ctx context.Context, projectID interface{}, opts MergeRequestOptions) (*domain.GitLabMergeRequest, error)
	CreateFile(ctx context.Context, projectID interface{}, filePath, content, commitMessage, branch string) error
	GetFileContent(ctx context.Context, projectID interface{}, filePath, ref string) (string, error)
	GetRepository(ctx context.Context, projectID interface{}) (*gitlab.Project, error)
	ListBranches(ctx context.Context, projectID interface{}, opts *gitlab.ListBranchesOptions) ([]*gitlab.Branch, error)
	CreateBranch(ctx context.Context, projectID interface{}, branch, ref string) (*gitlab.Branch, error)
	ListCommits(ctx context.Context, projectID interface{}, opts *gitlab.ListCommitsOptions) ([]*gitlab.Commit, error)
	GetCommit(ctx context.Context, projectID interface{}, sha string) (*gitlab.Commit, error)
	ListTags(ctx context.Context, projectID interface{}, opts *gitlab.ListTagsOptions) ([]*gitlab.Tag, error)
	CreateTag(ctx context.Context, projectID interface{}, tag, ref string, message string) (*gitlab.Tag, error)
}

// MergeRequestOptions MR 생성 옵션
type MergeRequestOptions struct {
	SourceBranch       string
	TargetBranch       string
	Title              string
	Description        string
	AssigneeIDs        []int
	ReviewerIDs        []int
	Labels             []string
	MilestoneID        *int
	RemoveSourceBranch bool
	Squash             bool
}

// GitLabClient GitLab API 클라이언트 구현체
type GitLabClient struct {
	client *gitlab.Client
}

// NewClient 새 GitLab 클라이언트 생성
func NewClient(cfg *config.GitLabConfig) (Client, error) {
	var cli *gitlab.Client
	var err error

	if cfg.BaseURL != "" {
		cli, err = gitlab.NewClient(cfg.Token, gitlab.WithBaseURL(cfg.BaseURL))
	} else {
		cli, err = gitlab.NewClient(cfg.Token)
	}

	if err != nil {
		return nil, fmt.Errorf("GitLab 클라이언트 생성 오류: %w", err)
	}

	return &GitLabClient{
		client: cli,
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

// ListMergeRequests GitLab 머지 리퀘스트 목록 조회
func (c *GitLabClient) ListMergeRequests(ctx context.Context, projectID interface{}, opts *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	mrs, _, err := c.client.MergeRequests.ListProjectMergeRequests(pid, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("GitLab API 오류 (ListProjectMergeRequests): %w", err)
	}
	return mrs, nil
}

// CreateMergeRequest GitLab 머지 리퀘스트 생성
func (c *GitLabClient) CreateMergeRequest(ctx context.Context, projectID interface{}, opts MergeRequestOptions) (*domain.GitLabMergeRequest, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	gitlabOpts := &gitlab.CreateMergeRequestOptions{
		Title:              &opts.Title,
		Description:        &opts.Description,
		SourceBranch:       &opts.SourceBranch,
		TargetBranch:       &opts.TargetBranch,
		AssigneeID:         nil,                        // AssigneeID is deprecated, use AssigneeIDs
		AssigneeIDs:        &opts.AssigneeIDs,          // Pass pointer to slice
		ReviewerIDs:        &opts.ReviewerIDs,          // Pass pointer to slice
		Labels:             gitlab.Labels(opts.Labels), // Convert []string to gitlab.Labels
		MilestoneID:        &opts.MilestoneID,
		RemoveSourceBranch: &opts.RemoveSourceBranch,
		Squash:             &opts.Squash,
	}

	mr, _, err := c.client.MergeRequests.CreateMergeRequest(pid, gitlabOpts, gitlab.WithContext(ctx))
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

// GetRepository GitLab 프로젝트 정보 조회
func (c *GitLabClient) GetRepository(ctx context.Context, projectID interface{}) (*gitlab.Project, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}
	project, _, err := c.client.Projects.GetProject(pid, nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("프로젝트 조회 오류: %w", err)
	}
	return project, nil
}

// ListBranches GitLab 브랜치 목록 조회
func (c *GitLabClient) ListBranches(ctx context.Context, projectID interface{}, opts *gitlab.ListBranchesOptions) ([]*gitlab.Branch, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}
	branches, _, err := c.client.Branches.ListBranches(pid, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("브랜치 목록 조회 오류: %w", err)
	}
	return branches, nil
}

// CreateBranch GitLab 브랜치 생성
func (c *GitLabClient) CreateBranch(ctx context.Context, projectID interface{}, branch, ref string) (*gitlab.Branch, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}
	opts := &gitlab.CreateBranchOptions{
		Branch: gitlab.String(branch),
		Ref:    gitlab.String(ref),
	}
	newBranch, _, err := c.client.Branches.CreateBranch(pid, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("브랜치 생성 오류: %w", err)
	}
	return newBranch, nil
}

// ListCommits GitLab 커밋 목록 조회
func (c *GitLabClient) ListCommits(ctx context.Context, projectID interface{}, opts *gitlab.ListCommitsOptions) ([]*gitlab.Commit, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}
	commits, _, err := c.client.Commits.ListCommits(pid, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("커밋 목록 조회 오류: %w", err)
	}
	return commits, nil
}

// GetCommit GitLab 커밋 조회
func (c *GitLabClient) GetCommit(ctx context.Context, projectID interface{}, sha string) (*gitlab.Commit, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}
	commit, _, err := c.client.Commits.GetCommit(pid, sha, nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("커밋 조회 오류: %w", err)
	}
	return commit, nil
}

// ListTags GitLab 태그 목록 조회
func (c *GitLabClient) ListTags(ctx context.Context, projectID interface{}, opts *gitlab.ListTagsOptions) ([]*gitlab.Tag, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}
	tags, _, err := c.client.Tags.ListTags(pid, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("태그 목록 조회 오류: %w", err)
	}
	return tags, nil
}

// CreateTag GitLab 태그 생성
func (c *GitLabClient) CreateTag(ctx context.Context, projectID interface{}, tag, ref string, message string) (*gitlab.Tag, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}
	opts := &gitlab.CreateTagOptions{
		TagName: gitlab.String(tag),
		Ref:     gitlab.String(ref),
		Message: gitlab.String(message),
	}
	newTag, _, err := c.client.Tags.CreateTag(pid, opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("태그 생성 오류: %w", err)
	}
	return newTag, nil
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
