package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/theshop/ai/internal/domain"
	"github.com/theshop/ai/modules/config"
)

// Client Jira API와의 상호작용을 위한 인터페이스
type Client interface {
	GetIssue(ctx context.Context, issueKey string) (*domain.JiraIssue, error)
	CreateIssue(ctx context.Context, projectKey, summary, description string) (*domain.JiraIssue, error)
	UpdateIssue(ctx context.Context, issueKey string, fields map[string]interface{}) error
	SearchIssues(ctx context.Context, jql string, maxResults int) ([]domain.JiraIssue, error)
}

// JiraClient Jira API 클라이언트 구현체
type JiraClient struct {
	baseURL    string
	username   string
	apiToken   string
	httpClient *http.Client
	projectKey string
}

// NewClient 새 Jira 클라이언트 생성
func NewClient(cfg *config.JiraConfig) Client {
	return &JiraClient{
		baseURL:    cfg.URL,
		username:   cfg.Username,
		apiToken:   cfg.APIToken,
		projectKey: cfg.ProjectKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetIssue Jira 이슈 조회
func (c *JiraClient) GetIssue(ctx context.Context, issueKey string) (*domain.JiraIssue, error) {
	url := fmt.Sprintf("%s/rest/api/3/issue/%s", c.baseURL, issueKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 오류: %w", err)
	}

	req.SetBasicAuth(c.username, c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 호출 오류: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 오류 (상태 코드: %d): %s", resp.StatusCode, string(body))
	}

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("응답 파싱 오류: %w", err)
	}

	fields, ok := responseData["fields"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("fields 필드가 올바르지 않음")
	}

	summary, _ := fields["summary"].(string)

	var description string
	if descMap, ok := fields["description"].(map[string]interface{}); ok {
		// ADF 형식 처리
		description = extractTextFromADF(descMap)
	}

	status := ""
	if statusMap, ok := fields["status"].(map[string]interface{}); ok {
		if name, ok := statusMap["name"].(string); ok {
			status = name
		}
	}

	issue := &domain.JiraIssue{
		Key:         issueKey,
		Summary:     summary,
		Description: description,
		Status:      status,
		Fields:      fields,
	}

	return issue, nil
}

// CreateIssue Jira 이슈 생성
func (c *JiraClient) CreateIssue(ctx context.Context, projectKey, summary, description string) (*domain.JiraIssue, error) {
	if projectKey == "" {
		projectKey = c.projectKey
	}

	// ADF 형식의 설명 생성
	adfDesc := map[string]interface{}{
		"type":    "doc",
		"version": 1,
		"content": []map[string]interface{}{
			{
				"type": "paragraph",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": description,
					},
				},
			},
		},
	}

	// 요청 본문 구성
	requestData := map[string]interface{}{
		"fields": map[string]interface{}{
			"project": map[string]string{
				"key": projectKey,
			},
			"summary":     summary,
			"description": adfDesc,
			"issuetype": map[string]string{
				"name": "Task",
			},
		},
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("요청 마샬링 오류: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/3/issue", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("요청 생성 오류: %w", err)
	}

	req.SetBasicAuth(c.username, c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 호출 오류: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 오류 (상태 코드: %d): %s", resp.StatusCode, string(body))
	}

	var responseData struct {
		ID  string `json:"id"`
		Key string `json:"key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("응답 파싱 오류: %w", err)
	}

	// 생성된 이슈 정보 조회
	return c.GetIssue(ctx, responseData.Key)
}

// UpdateIssue Jira 이슈 업데이트
func (c *JiraClient) UpdateIssue(ctx context.Context, issueKey string, fields map[string]interface{}) error {
	requestData := map[string]interface{}{
		"fields": fields,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("요청 마샬링 오류: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/3/issue/%s", c.baseURL, issueKey)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("요청 생성 오류: %w", err)
	}

	req.SetBasicAuth(c.username, c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("API 호출 오류: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API 오류 (상태 코드: %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// SearchIssues JQL을 사용하여 이슈 검색
func (c *JiraClient) SearchIssues(ctx context.Context, jql string, maxResults int) ([]domain.JiraIssue, error) {
	if maxResults <= 0 {
		maxResults = 10
	}

	requestData := map[string]interface{}{
		"jql":        jql,
		"maxResults": maxResults,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("요청 마샬링 오류: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/3/search", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("요청 생성 오류: %w", err)
	}

	req.SetBasicAuth(c.username, c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 호출 오류: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 오류 (상태 코드: %d): %s", resp.StatusCode, string(body))
	}

	var responseData struct {
		Issues []map[string]interface{} `json:"issues"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("응답 파싱 오류: %w", err)
	}

	var issues []domain.JiraIssue
	for _, issueData := range responseData.Issues {
		key, _ := issueData["key"].(string)

		fields, ok := issueData["fields"].(map[string]interface{})
		if !ok {
			continue
		}

		summary, _ := fields["summary"].(string)

		var description string
		if descMap, ok := fields["description"].(map[string]interface{}); ok {
			description = extractTextFromADF(descMap)
		}

		status := ""
		if statusMap, ok := fields["status"].(map[string]interface{}); ok {
			if name, ok := statusMap["name"].(string); ok {
				status = name
			}
		}

		issue := domain.JiraIssue{
			Key:         key,
			Summary:     summary,
			Description: description,
			Status:      status,
			Fields:      fields,
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// extractTextFromADF ADF 형식에서 텍스트 추출 (간소화된 구현)
func extractTextFromADF(adf map[string]interface{}) string {
	var text string

	// content 배열 처리
	if content, ok := adf["content"].([]interface{}); ok {
		for _, item := range content {
			if itemMap, ok := item.(map[string]interface{}); ok {
				// 텍스트 노드 처리
				if itemType, ok := itemMap["type"].(string); ok {
					if itemType == "text" {
						if t, ok := itemMap["text"].(string); ok {
							text += t
						}
					} else if itemType == "paragraph" || itemType == "heading" {
						// 재귀적으로 내부 content 처리
						if t := extractTextFromADF(itemMap); t != "" {
							text += t + "\n"
						}
					}
				}
			}
		}
	}

	return text
}
