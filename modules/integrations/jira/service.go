// modules/integrations/jira/service.go
package jira

import (
	"context"

	"github.com/theshop/ai/internal/domain"
)

// JiraService provides higher-level operations on top of the Jira client
type JiraService struct {
	client Client
}

// NewJiraService creates a new Jira service
func NewJiraService(client Client) *JiraService {
	return &JiraService{
		client: client,
	}
}

// GetIssue retrieves an issue from Jira
func (s *JiraService) GetIssue(ctx context.Context, issueKey string) (*domain.JiraIssue, error) {
	return s.client.GetIssue(ctx, issueKey)
}

// CreateIssueWithService creates a Jira issue using the service
// This is a different method from the client's CreateIssue
func (s *JiraService) CreateIssueWithService(ctx context.Context, projectKey, summary, description string) (*domain.JiraIssue, error) {
	// ADF description creation logic moved to client.CreateIssue
	// Remove unused variable:
	// adfDesc := map[string]interface{}{ ... }

	// You can add additional service-level logic here if needed
	// e.g., validation, logging, etc.

	return s.client.CreateIssue(ctx, projectKey, summary, description)
}
