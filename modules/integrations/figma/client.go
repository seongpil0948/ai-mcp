package figma

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/theshop/ai/modules/config"
)

const (
	figmaAPIBaseURL = "https://api.figma.com/v1"
)

// Client is a client for the Figma API
type Client interface {
	GetFile(ctx context.Context, fileKey string) (*File, error)
	GetFileNodes(ctx context.Context, fileKey string, nodeIDs []string) (*FileNodes, error)
	GetImage(ctx context.Context, fileKey string, opts *GetImageOptions) (*ImageResponse, error)
	GetTeamProjects(ctx context.Context, teamID string) (*TeamProjectsResponse, error)
	GetProjectFiles(ctx context.Context, projectID string) (*ProjectFilesResponse, error)
	CreateComment(ctx context.Context, fileKey string, message string, position *CommentPosition) (*Comment, error)
	CreateFile(ctx context.Context, name string, teamID string) (*File, error)
}

// FigmaClient implements the Client interface
type FigmaClient struct {
	httpClient  *http.Client
	accessToken string
	baseURL     string
	rateLimiter *time.Ticker // Simple rate limiter
}

// NewClient creates a new Figma API client
func NewClient(cfg *config.FigmaConfig) (Client, error) {
	if cfg.AccessToken == "" {
		return nil, fmt.Errorf("Figma access token is required")
	}

	return &FigmaClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		accessToken: cfg.AccessToken,
		baseURL:     figmaAPIBaseURL,
		rateLimiter: time.NewTicker(500 * time.Millisecond), // 2 requests per second
	}, nil
}

// doRequest makes a request to the Figma API
func (c *FigmaClient) doRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
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
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Figma-Client", "theshop-ai-assistant")

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
		return nil, fmt.Errorf("API error: %s (status: %d)", apiErr.Error.Message, resp.StatusCode)
	}

	return respBody, nil
}

// GetFile gets a Figma file
func (c *FigmaClient) GetFile(ctx context.Context, fileKey string) (*File, error) {
	respBody, err := c.doRequest(ctx, http.MethodGet, "/files/"+fileKey, nil)
	if err != nil {
		return nil, err
	}

	var file File
	if err := json.Unmarshal(respBody, &file); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file response: %w", err)
	}

	return &file, nil
}

// GetFileNodes gets specific nodes from a Figma file
func (c *FigmaClient) GetFileNodes(ctx context.Context, fileKey string, nodeIDs []string) (*FileNodes, error) {
	// Create comma-separated list of node IDs
	nodeIDsStr := ""
	for i, id := range nodeIDs {
		if i > 0 {
			nodeIDsStr += ","
		}
		nodeIDsStr += id
	}

	// Build query parameters
	params := url.Values{}
	params.Add("ids", nodeIDsStr)

	endpoint := fmt.Sprintf("/files/%s/nodes?%s", fileKey, params.Encode())
	respBody, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var fileNodes FileNodes
	if err := json.Unmarshal(respBody, &fileNodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file nodes response: %w", err)
	}

	return &fileNodes, nil
}

// GetImage gets images from a Figma file
func (c *FigmaClient) GetImage(ctx context.Context, fileKey string, opts *GetImageOptions) (*ImageResponse, error) {
	// Build query parameters
	params := url.Values{}
	if opts != nil {
		if len(opts.IDs) > 0 {
			idsStr := ""
			for i, id := range opts.IDs {
				if i > 0 {
					idsStr += ","
				}
				idsStr += id
			}
			params.Add("ids", idsStr)
		}
		if opts.Scale != 0 {
			params.Add("scale", fmt.Sprintf("%f", opts.Scale))
		}
		if opts.Format != "" {
			params.Add("format", opts.Format)
		}
	}

	endpoint := fmt.Sprintf("/images/%s?%s", fileKey, params.Encode())
	respBody, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var imageResponse ImageResponse
	if err := json.Unmarshal(respBody, &imageResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal image response: %w", err)
	}

	return &imageResponse, nil
}

// GetTeamProjects gets projects for a team
func (c *FigmaClient) GetTeamProjects(ctx context.Context, teamID string) (*TeamProjectsResponse, error) {
	endpoint := fmt.Sprintf("/teams/%s/projects", teamID)
	respBody, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response TeamProjectsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team projects response: %w", err)
	}

	return &response, nil
}

// GetProjectFiles gets files for a project
func (c *FigmaClient) GetProjectFiles(ctx context.Context, projectID string) (*ProjectFilesResponse, error) {
	endpoint := fmt.Sprintf("/projects/%s/files", projectID)
	respBody, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response ProjectFilesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project files response: %w", err)
	}

	return &response, nil
}

// CreateComment creates a comment on a file
func (c *FigmaClient) CreateComment(ctx context.Context, fileKey string, message string, position *CommentPosition) (*Comment, error) {
	body := map[string]interface{}{
		"message": message,
	}
	if position != nil {
		body["client_meta"] = map[string]interface{}{
			"x":       position.X,
			"y":       position.Y,
			"node_id": position.NodeID,
		}
	}

	endpoint := fmt.Sprintf("/files/%s/comments", fileKey)
	respBody, err := c.doRequest(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}

	var comment Comment
	if err := json.Unmarshal(respBody, &comment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal comment response: %w", err)
	}

	return &comment, nil
}

// CreateFile creates a new Figma file
func (c *FigmaClient) CreateFile(ctx context.Context, name string, teamID string) (*File, error) {
	body := map[string]interface{}{
		"name": name,
	}

	endpoint := fmt.Sprintf("/teams/%s/files", teamID)
	respBody, err := c.doRequest(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}

	var file File
	if err := json.Unmarshal(respBody, &file); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file response: %w", err)
	}

	return &file, nil
}
