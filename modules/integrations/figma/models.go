package figma

// APIError represents an error returned by the Figma API
type APIError struct {
	Status int `json:"status"`
	Error  struct {
		Message string `json:"message"`
	} `json:"error"`
}

// File represents a Figma file
type File struct {
	Name          string               `json:"name"`
	LastModified  string               `json:"lastModified"`
	ThumbnailURL  string               `json:"thumbnailUrl"`
	Version       string               `json:"version"`
	Document      Document             `json:"document"`
	Components    map[string]Component `json:"components"`
	Styles        map[string]Style     `json:"styles"`
	SchemaVersion int                  `json:"schemaVersion"`
	MainFileKey   string               `json:"mainFileKey,omitempty"`
	Branches      []Branch             `json:"branches,omitempty"`
}

// Document represents a Figma document
type Document struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Children []Node `json:"children"`
}

// Node represents a Figma node (element)
type Node struct {
	ID                  string          `json:"id"`
	Name                string          `json:"name"`
	Type                string          `json:"type"`
	Visible             bool            `json:"visible"`
	Children            []Node          `json:"children,omitempty"`
	Characters          string          `json:"characters,omitempty"`
	Style               *Style          `json:"style,omitempty"`
	Fills               []Fill          `json:"fills,omitempty"`
	Strokes             []Fill          `json:"strokes,omitempty"`
	ExportSettings      []ExportSetting `json:"exportSettings,omitempty"`
	AbsoluteBoundingBox *Rectangle      `json:"absoluteBoundingBox,omitempty"`
}

// Rectangle represents a rectangle
type Rectangle struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Style represents the style of a node
type Style struct {
	ID             string      `json:"id,omitempty"`
	Key            string      `json:"key,omitempty"`
	Name           string      `json:"name,omitempty"`
	Description    string      `json:"description,omitempty"`
	FontFamily     string      `json:"fontFamily,omitempty"`
	FontWeight     int         `json:"fontWeight,omitempty"`
	FontSize       float64     `json:"fontSize,omitempty"`
	TextAlign      string      `json:"textAlign,omitempty"`
	TextDecoration string      `json:"textDecoration,omitempty"`
	LineHeight     interface{} `json:"lineHeight,omitempty"`
	LetterSpacing  interface{} `json:"letterSpacing,omitempty"`
	Fills          []Fill      `json:"fills,omitempty"`
	StrokeWeight   float64     `json:"strokeWeight,omitempty"`
	Strokes        []Fill      `json:"strokes,omitempty"`
}

// Fill represents a fill (color, gradient, etc.)
type Fill struct {
	Type            string           `json:"type"`
	Visible         bool             `json:"visible"`
	Opacity         float64          `json:"opacity,omitempty"`
	Color           *Color           `json:"color,omitempty"`
	GradientHandles *GradientHandles `json:"gradientHandles,omitempty"`
	GradientStops   []GradientStop   `json:"gradientStops,omitempty"`
}

// Color represents an RGBA color
type Color struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
	A float64 `json:"a"`
}

// GradientHandles represents handles for a gradient
type GradientHandles struct {
	Start *Vector `json:"start"`
	End   *Vector `json:"end"`
}

// Vector represents a 2D vector
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// GradientStop represents a stop in a gradient
type GradientStop struct {
	Position float64 `json:"position"`
	Color    *Color  `json:"color"`
}

// Component represents a Figma component
type Component struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Branch represents a branch in a Figma file
type Branch struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnailUrl"`
	LastModified string `json:"lastModified"`
	LinkAccess   string `json:"linkAccess"`
}

// ExportSetting represents export settings for a node
type ExportSetting struct {
	Suffix     string      `json:"suffix"`
	Format     string      `json:"format"`
	Constraint *Constraint `json:"constraint"`
}

// Constraint represents size constraints for exports
type Constraint struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

// FileNodes represents a response containing specific nodes from a file
type FileNodes struct {
	Name         string               `json:"name"`
	LastModified string               `json:"lastModified"`
	ThumbnailURL string               `json:"thumbnailUrl"`
	Version      string               `json:"version"`
	Err          string               `json:"err,omitempty"`
	Nodes        map[string]*FileNode `json:"nodes"`
}

// FileNode represents a node in a FileNodes response
type FileNode struct {
	Document      *Node                `json:"document"`
	Components    map[string]Component `json:"components"`
	SchemaVersion int                  `json:"schemaVersion"`
	Styles        map[string]Style     `json:"styles"`
}

// GetImageOptions represents options for the GetImage request
type GetImageOptions struct {
	IDs    []string
	Scale  float64
	Format string // "jpg", "png", "svg", or "pdf"
}

// ImageResponse represents a response from the GetImage request
type ImageResponse struct {
	Err    string            `json:"err,omitempty"`
	Images map[string]string `json:"images"`
}

// TeamProjectsResponse represents a response from the GetTeamProjects request
type TeamProjectsResponse struct {
	Projects []Project `json:"projects"`
}

// Project represents a Figma project
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProjectFilesResponse represents a response from the GetProjectFiles request
type ProjectFilesResponse struct {
	Files []ProjectFile `json:"files"`
}

// ProjectFile represents a file in a project
type ProjectFile struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	LastModified string `json:"last_modified"`
}

// CommentPosition represents the position of a comment in a Figma file
type CommentPosition struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	NodeID string  `json:"node_id"`
}

// Comment represents a comment in a Figma file
type Comment struct {
	ID         string           `json:"id"`
	Message    string           `json:"message"`
	ClientMeta *CommentPosition `json:"client_meta,omitempty"`
	CreatedAt  string           `json:"created_at"`
	OrderID    int              `json:"order_id"`
	User       UserInfo         `json:"user"`
}

// UserInfo represents information about a Figma user
type UserInfo struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	ImgURL string `json:"img_url"`
}

// Constants for node types
const (
	NodeTypeDocument  = "DOCUMENT"
	NodeTypeCanvas    = "CANVAS"
	NodeTypeFrame     = "FRAME"
	NodeTypeGroup     = "GROUP"
	NodeTypeVector    = "VECTOR"
	NodeTypeText      = "TEXT"
	NodeTypeRectangle = "RECTANGLE"
	NodeTypeEllipse   = "ELLIPSE"
	NodeTypeLine      = "LINE"
	NodeTypeComponent = "COMPONENT"
	NodeTypeInstance  = "INSTANCE"
)
