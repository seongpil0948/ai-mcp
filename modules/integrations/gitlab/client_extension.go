// Path: modules/integrations/gitlab/client_extension.go
package gitlab

import (
	"context"

	"github.com/theshop/ai/internal/domain"
	"github.com/xanzy/go-gitlab"
)

// Extended interfaces and types for GitLab client

// ListMergeRequestsOptions represents options for listing merge requests
type ListMergeRequestsOptions = gitlab.ListProjectMergeRequestsOptions

// MergeRequestOptions represents options for creating a merge request
type MergeRequestOptions struct {
	SourceBranch       string
	TargetBranch       string
	Title              string
	Description        string
	RemoveSourceBranch bool
	Squash             bool
}

// Client extends the existing Client interface with additional methods
type Client interface {
	// Existing methods from domain/models.go
	GetProject(ctx context.Context, projectID interface{}) (*domain.GitLabProject, error)
	ListProjects(ctx context.Context, search string, page, perPage int) ([]domain.GitLabProject, error)
	GetMergeRequest(ctx context.Context, projectID interface{}, mergeRequestIID int) (*domain.GitLabMergeRequest, error)
	CreateMergeRequest(ctx context.Context, projectID interface{}, opts MergeRequestOptions) (*domain.GitLabMergeRequest, error)
	CreateFile(ctx context.Context, projectID interface{}, filePath, content, commitMessage, branch string) error
	GetFileContent(ctx context.Context, projectID interface{}, filePath, ref string) (string, error)

	// Additional methods for extended functionality
	ListMergeRequests(ctx context.Context, projectID interface{}, options *ListMergeRequestsOptions) ([]domain.GitLabMergeRequest, error)
	GetRepository(ctx context.Context, projectID interface{}) (*gitlab.Repository, error)
	ListBranches(ctx context.Context, projectID interface{}) ([]*gitlab.Branch, error)
	CreateBranch(ctx context.Context, projectID interface{}, branch, ref string) (*gitlab.Branch, error)
	ListCommits(ctx context.Context, projectID interface{}, options *gitlab.ListCommitsOptions) ([]*gitlab.Commit, error)
	GetCommit(ctx context.Context, projectID interface{}, sha string) (*gitlab.Commit, error)
	ListTags(ctx context.Context, projectID interface{}, options *gitlab.ListTagsOptions) ([]*gitlab.Tag, error)
	CreateTag(ctx context.Context, projectID interface{}, tag, ref string, message string) (*gitlab.Tag, error)
}

// Ensure GitLabClient implements the extended Client interface
var _ Client = (*GitLabClient)(nil)

// ListMergeRequests lists merge requests for a project
func (c *GitLabClient) ListMergeRequests(ctx context.Context, projectID interface{}, options *ListMergeRequestsOptions) ([]domain.GitLabMergeRequest, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	mrs, _, err := c.client.MergeRequests.ListProjectMergeRequests(pid, options, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	result := make([]domain.GitLabMergeRequest, len(mrs))
	for i, mr := range mrs {
		result[i] = domain.GitLabMergeRequest{
			ID:          mr.ID,
			IID:         mr.IID,
			Title:       mr.Title,
			Description: mr.Description,
			WebURL:      mr.WebURL,
			State:       mr.State,
			CreatedAt:   *mr.CreatedAt,
		}
	}

	return result, nil
}

// GetRepository gets the repository for a project
func (c *GitLabClient) GetRepository(ctx context.Context, projectID interface{}) (*gitlab.Repository, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	repo, _, err := c.client.Repositories.GetRepository(pid, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// ListBranches lists branches for a project
func (c *GitLabClient) ListBranches(ctx context.Context, projectID interface{}) ([]*gitlab.Branch, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	branches, _, err := c.client.Branches.ListBranches(pid, &gitlab.ListBranchesOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return branches, nil
}

// CreateBranch creates a new branch
func (c *GitLabClient) CreateBranch(ctx context.Context, projectID interface{}, branch, ref string) (*gitlab.Branch, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	createOpts := &gitlab.CreateBranchOptions{
		Branch: gitlab.String(branch),
		Ref:    gitlab.String(ref),
	}

	newBranch, _, err := c.client.Branches.CreateBranch(pid, createOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return newBranch, nil
}

// ListCommits lists commits for a project
func (c *GitLabClient) ListCommits(ctx context.Context, projectID interface{}, options *gitlab.ListCommitsOptions) ([]*gitlab.Commit, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	commits, _, err := c.client.Commits.ListCommits(pid, options, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return commits, nil
}

// GetCommit gets a specific commit
func (c *GitLabClient) GetCommit(ctx context.Context, projectID interface{}, sha string) (*gitlab.Commit, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	commit, _, err := c.client.Commits.GetCommit(pid, sha, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return commit, nil
}

// ListTags lists tags for a project
func (c *GitLabClient) ListTags(ctx context.Context, projectID interface{}, options *gitlab.ListTagsOptions) ([]*gitlab.Tag, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	tags, _, err := c.client.Tags.ListTags(pid, options, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// CreateTag creates a new tag
func (c *GitLabClient) CreateTag(ctx context.Context, projectID interface{}, tag, ref string, message string) (*gitlab.Tag, error) {
	pid, err := parseProjectID(projectID)
	if err != nil {
		return nil, err
	}

	createOpts := &gitlab.CreateTagOptions{
		TagName: gitlab.String(tag),
		Ref:     gitlab.String(ref),
		Message: gitlab.String(message),
	}

	newTag, _, err := c.client.Tags.CreateTag(pid, createOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return newTag, nil
}
