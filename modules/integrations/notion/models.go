package notion

import (
	"fmt"
	"time"
)

// APIError represents an error returned by the Notion API
type APIError struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Page represents a Notion page
type Page struct {
	Object         string                 `json:"object"`
	ID             string                 `json:"id"`
	CreatedTime    time.Time              `json:"created_time"`
	LastEditedTime time.Time              `json:"last_edited_time"`
	Parent         Parent                 `json:"parent"`
	Archived       bool                   `json:"archived"`
	Properties     map[string]interface{} `json:"properties"`
	URL            string                 `json:"url"`
}

// Database represents a Notion database
type Database struct {
	Object         string                 `json:"object"`
	ID             string                 `json:"id"`
	CreatedTime    time.Time              `json:"created_time"`
	LastEditedTime time.Time              `json:"last_edited_time"`
	Title          []RichText             `json:"title"`
	Properties     map[string]interface{} `json:"properties"`
	URL            string                 `json:"url"`
}

// Block represents a Notion block
type Block struct {
	Object         string                 `json:"object"`
	ID             string                 `json:"id"`
	Type           string                 `json:"type"`
	CreatedTime    time.Time              `json:"created_time"`
	LastEditedTime time.Time              `json:"last_edited_time"`
	HasChildren    bool                   `json:"has_children"`
	Archived       bool                   `json:"archived"`
	Content        map[string]interface{} `json:"content,omitempty"`
	// Each block type has specific content in this map based on Type
}

// RichText represents a rich text object in Notion
type RichText struct {
	Type        string                 `json:"type"`
	Text        TextContent            `json:"text,omitempty"`
	Annotations map[string]interface{} `json:"annotations,omitempty"`
	PlainText   string                 `json:"plain_text"`
	Href        string                 `json:"href,omitempty"`
}

// TextContent represents the content of a text object
type TextContent struct {
	Content string `json:"content"`
	Link    *Link  `json:"link,omitempty"`
}

// Link represents a link in Notion
type Link struct {
	URL string `json:"url"`
}

// Parent represents a parent object reference
type Parent struct {
	Type       string `json:"type"` // "database_id", "page_id", "workspace"
	DatabaseID string `json:"database_id,omitempty"`
	PageID     string `json:"page_id,omitempty"`
}

// DatabaseQuery represents a database query
type DatabaseQuery struct {
	Filter      interface{}   `json:"filter,omitempty"`
	Sorts       []interface{} `json:"sorts,omitempty"`
	StartCursor string        `json:"start_cursor,omitempty"`
	PageSize    int           `json:"page_size,omitempty"`
}

// DatabaseQueryResult represents the result of a database query
type DatabaseQueryResult struct {
	Object     string `json:"object"`
	Results    []Page `json:"results"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

// CreatePageParams represents the parameters for creating a page
type CreatePageParams struct {
	Parent     Parent                 `json:"parent"`
	Properties map[string]interface{} `json:"properties"`
	Children   []interface{}          `json:"children,omitempty"`
}

// UpdatePageParams represents the parameters for updating a page
type UpdatePageParams struct {
	Properties map[string]interface{} `json:"properties"`
	Archived   *bool                  `json:"archived,omitempty"`
}

// SearchOptions represents the options for a search
type SearchOptions struct {
	Query       string      `json:"query,omitempty"`
	Sort        interface{} `json:"sort,omitempty"`
	Filter      interface{} `json:"filter,omitempty"`
	StartCursor string      `json:"start_cursor,omitempty"`
	PageSize    int         `json:"page_size,omitempty"`
}

// SearchResults represents the results of a search
type SearchResults struct {
	Object     string        `json:"object"`
	Results    []interface{} `json:"results"`
	NextCursor string        `json:"next_cursor"`
	HasMore    bool          `json:"has_more"`
}

// PaginationOptions represents pagination options
type PaginationOptions struct {
	StartCursor string `json:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
}

// BlockChildrenResults represents the results of getting block children
type BlockChildrenResults struct {
	Object     string  `json:"object"`
	Results    []Block `json:"results"`
	NextCursor string  `json:"next_cursor"`
	HasMore    bool    `json:"has_more"`
}

// CreateBlockOptions represents options for creating a block
type CreateBlockOptions struct {
	Children []Block `json:"children"`
}

// BlockTypes constants for block types
const (
	BlockTypeParagraph        = "paragraph"
	BlockTypeHeading1         = "heading_1"
	BlockTypeHeading2         = "heading_2"
	BlockTypeHeading3         = "heading_3"
	BlockTypeBulletedListItem = "bulleted_list_item"
	BlockTypeNumberedListItem = "numbered_list_item"
	BlockTypeToggle           = "toggle"
	BlockTypeCode             = "code"
	BlockTypeImage            = "image"
	BlockTypeVideo            = "video"
	BlockTypeFile             = "file"
	BlockTypeDivider          = "divider"
	BlockTypeQuote            = "quote"
	BlockTypeToDo             = "to_do"
	BlockTypeBookmark         = "bookmark"
	BlockTypeCallout          = "callout"
	BlockTypeTable            = "table"
	BlockTypeTableRow         = "table_row"
)

// Helper functions to create common block types
func NewParagraphBlock(text string) Block {
	return Block{
		Object:      "block",
		Type:        BlockTypeParagraph,
		HasChildren: false,
		Content: map[string]interface{}{
			BlockTypeParagraph: map[string]interface{}{
				"rich_text": []RichText{
					{
						Type: "text",
						Text: TextContent{
							Content: text,
						},
						PlainText: text,
					},
				},
			},
		},
	}
}

func NewHeadingBlock(text string, level int) Block {
	headingType := fmt.Sprintf("heading_%d", level)
	return Block{
		Object:      "block",
		Type:        headingType,
		HasChildren: false,
		Content: map[string]interface{}{
			headingType: map[string]interface{}{
				"rich_text": []RichText{
					{
						Type: "text",
						Text: TextContent{
							Content: text,
						},
						PlainText: text,
					},
				},
			},
		},
	}
}
