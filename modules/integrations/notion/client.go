package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/theshop/ai/modules/config"
)

const (
	notionAPIBaseURL = "https://api.notion.com/v1"
	notionVersion    = "2022-06-28" // Current Notion API version
)

// Client is a client for the Notion API
type Client interface {
	GetPage(ctx context.Context, pageID string) (*Page, error)
	GetDatabase(ctx context.Context, databaseID string) (*Database, error)
	QueryDatabase(ctx context.Context, databaseID string, query *DatabaseQuery) (*DatabaseQueryResult, error)
	CreatePage(ctx context.Context, params *CreatePageParams) (*Page, error)
	UpdatePage(ctx context.Context, pageID string, params *UpdatePageParams) (*Page, error)
	Search(ctx context.Context, query string, options *SearchOptions) (*SearchResults, error)
	GetBlock(ctx context.Context, blockID string) (*Block, error)
	GetBlockChildren(ctx context.Context, blockID string, options *PaginationOptions) (*BlockChildrenResults, error)
	AppendBlockChildren(ctx context.Context, blockID string, children []Block) (*BlockChildrenResults, error)
}

// NotionClient implements the Client interface
type NotionClient struct {
	httpClient  *http.Client
	apiKey      string
	apiVersion  string
	baseURL     string
	rateLimiter *time.Ticker // Simple rate limiter
}

// NewClient creates a new Notion API client
func NewClient(cfg *config.NotionConfig) (Client, error) {
	if cfg.APIKey == "" && cfg.Token == "" {
		return nil, fmt.Errorf("either Notion API key or token is required")
	}

	// Use API key if provided, otherwise fall back to token
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = cfg.Token
	}

	return &NotionClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:      apiKey,
		apiVersion:  notionVersion,
		baseURL:     notionAPIBaseURL,
		rateLimiter: time.NewTicker(200 * time.Millisecond), // 5 requests per second
	}, nil
}

// doRequest makes a request to the Notion API
func (c *NotionClient) doRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
	// Rate limiting - wait for the next tick
	<-c.rateLimiter.C

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Notion-Version", c.apiVersion)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("API error: %s (code: %s)", apiErr.Message, apiErr.Code)
	}

	return respBody, nil
}

// GetPage retrieves a page by ID
func (c *NotionClient) GetPage(ctx context.Context, pageID string) (*Page, error) {
	respBody, err := c.doRequest(ctx, http.MethodGet, "/pages/"+pageID, nil)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := json.Unmarshal(respBody, &page); err != nil {
		return nil, fmt.Errorf("failed to unmarshal page response: %w", err)
	}

	return &page, nil
}

// GetDatabase retrieves a database by ID
func (c *NotionClient) GetDatabase(ctx context.Context, databaseID string) (*Database, error) {
	respBody, err := c.doRequest(ctx, http.MethodGet, "/databases/"+databaseID, nil)
	if err != nil {
		return nil, err
	}

	var database Database
	if err := json.Unmarshal(respBody, &database); err != nil {
		return nil, fmt.Errorf("failed to unmarshal database response: %w", err)
	}

	return &database, nil
}

// QueryDatabase queries a database
func (c *NotionClient) QueryDatabase(ctx context.Context, databaseID string, query *DatabaseQuery) (*DatabaseQueryResult, error) {
	respBody, err := c.doRequest(ctx, http.MethodPost, "/databases/"+databaseID+"/query", query)
	if err != nil {
		return nil, err
	}

	var result DatabaseQueryResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal database query response: %w", err)
	}

	return &result, nil
}

// CreatePage creates a new page
func (c *NotionClient) CreatePage(ctx context.Context, params *CreatePageParams) (*Page, error) {
	respBody, err := c.doRequest(ctx, http.MethodPost, "/pages", params)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := json.Unmarshal(respBody, &page); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create page response: %w", err)
	}

	return &page, nil
}

// UpdatePage updates a page
func (c *NotionClient) UpdatePage(ctx context.Context, pageID string, params *UpdatePageParams) (*Page, error) {
	respBody, err := c.doRequest(ctx, http.MethodPatch, "/pages/"+pageID, params)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := json.Unmarshal(respBody, &page); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update page response: %w", err)
	}

	return &page, nil
}

// Search searches for pages and databases
func (c *NotionClient) Search(ctx context.Context, query string, options *SearchOptions) (*SearchResults, error) {
	if options == nil {
		options = &SearchOptions{}
	}
	options.Query = query

	respBody, err := c.doRequest(ctx, http.MethodPost, "/search", options)
	if err != nil {
		return nil, err
	}

	var results SearchResults
	if err := json.Unmarshal(respBody, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	return &results, nil
}

// GetBlock retrieves a block by ID
func (c *NotionClient) GetBlock(ctx context.Context, blockID string) (*Block, error) {
	respBody, err := c.doRequest(ctx, http.MethodGet, "/blocks/"+blockID, nil)
	if err != nil {
		return nil, err
	}

	var block Block
	if err := json.Unmarshal(respBody, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block response: %w", err)
	}

	return &block, nil
}

// GetBlockChildren retrieves children of a block
func (c *NotionClient) GetBlockChildren(ctx context.Context, blockID string, options *PaginationOptions) (*BlockChildrenResults, error) {
	endpoint := "/blocks/" + blockID + "/children"

	if options != nil {
		if options.StartCursor != "" {
			endpoint += "?start_cursor=" + options.StartCursor
		}
		if options.PageSize > 0 {
			if options.StartCursor != "" {
				endpoint += "&"
			} else {
				endpoint += "?"
			}
			endpoint += fmt.Sprintf("page_size=%d", options.PageSize)
		}
	}

	respBody, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var results BlockChildrenResults
	if err := json.Unmarshal(respBody, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block children response: %w", err)
	}

	return &results, nil
}

// AppendBlockChildren appends children to a block
func (c *NotionClient) AppendBlockChildren(ctx context.Context, blockID string, children []Block) (*BlockChildrenResults, error) {
	body := map[string]interface{}{
		"children": children,
	}

	respBody, err := c.doRequest(ctx, http.MethodPatch, "/blocks/"+blockID+"/children", body)
	if err != nil {
		return nil, err
	}

	var results BlockChildrenResults
	if err := json.Unmarshal(respBody, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal append block children response: %w", err)
	}

	return &results, nil
}
