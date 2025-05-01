// Path: cmd/mcpserver/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/theshop/ai/modules/config"
	"github.com/theshop/ai/modules/integrations/figma"
	"github.com/theshop/ai/modules/integrations/gitlab"
	"github.com/theshop/ai/modules/integrations/notion"
	"github.com/theshop/ai/pkg/mcpserver"
)

var (
	httpAddr    = flag.String("http", "", "HTTP service address (e.g., ':8080')")
	stdio       = flag.Bool("stdio", false, "Use stdio transport")
	configPath  = flag.String("config", "config.yaml", "Path to config file")
	showVersion = flag.Bool("version", false, "Show version and exit")
)

const (
	appName    = "theshop-ai-mcp-server"
	appVersion = "0.1.0"
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	log.Printf("Starting %s v%s", appName, appVersion)

	if (*httpAddr == "" && !*stdio) || (*httpAddr != "" && *stdio) {
		fmt.Println("Error: Please specify exactly one transport: -http ADDR or -stdio")
		flag.Usage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize services/integrations based on config
	var handlers []mcpserver.ToolHandler
	if cfg.Integrations.GitLab.Token != "" {
		gitlabService, err := gitlab.NewService(&cfg.Integrations.GitLab)
		if err != nil {
			log.Printf("Failed to initialize GitLab service: %v", err)
		} else {
			handlers = append(handlers, gitlabService)
			log.Println("GitLab service initialized")
		}
	}
	if cfg.Integrations.Figma.Token != "" {
		figmaService, err := figma.NewService(&cfg.Integrations.Figma)
		if err != nil {
			log.Printf("Failed to initialize Figma service: %v", err)
		} else {
			handlers = append(handlers, figmaService)
			log.Println("Figma service initialized")
		}
	}
	if cfg.Integrations.Notion.Token != "" {
		notionService, err := notion.NewService(&cfg.Integrations.Notion)
		if err != nil {
			log.Printf("Failed to initialize Notion service: %v", err)
		} else {
			handlers = append(handlers, notionService)
			log.Println("Notion service initialized")
		}
	}

	if len(handlers) == 0 {
		log.Fatalf("No integration services initialized. Check config.")
	}

	// Combine handlers if multiple are active
	combinedHandler := mcpserver.NewCombinedToolHandler(handlers...)

	// Create MCP Server
	server := mcpserver.New(combinedHandler)

	// Register transports based on flags
	if *httpAddr != "" {
		if err := server.RegisterHTTPTransport(*httpAddr, "/mcp"); err != nil {
			log.Fatalf("Failed to register HTTP transport: %v", err)
		}
	}
	if *stdio {
		if err := server.RegisterStdioTransport(); err != nil {
			log.Fatalf("Failed to register Stdio transport: %v", err)
		}
	}

	// Start server in a goroutine
	go func() {
		if err := server.Run(); err != nil {
			log.Printf("Server run error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// Helper function to combine multiple ToolHandlers (add this to mcpserver package or here)
// Example implementation:
/*
package mcpserver

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"sync"
)

type CombinedToolHandler struct {
	handlers []ToolHandler
	toolMap  map[string]ToolHandler
	mapOnce  sync.Once
	mapErr   error
}

func NewCombinedToolHandler(handlers ...ToolHandler) *CombinedToolHandler {
	return &CombinedToolHandler{handlers: handlers}
}

func (c *CombinedToolHandler) buildMap(ctx context.Context) error {
	c.mapOnce.Do(func() {
		c.toolMap = make(map[string]ToolHandler)
		for _, h := range c.handlers {
			tools, err := h.GetTools(ctx)
			if err != nil {
				c.mapErr = fmt.Errorf("failed to get tools from a handler: %w", err)
				return
			}
			for _, t := range tools {
				if _, exists := c.toolMap[t.GetName()]; exists {
					c.mapErr = fmt.Errorf("duplicate tool name detected: %s", t.GetName())
					return
				}
				c.toolMap[t.GetName()] = h
			}
		}
	})
	return c.mapErr
}


func (c *CombinedToolHandler) GetTools(ctx context.Context) ([]*mcp.Tool, error) {
	if err := c.buildMap(ctx); err != nil {
		return nil, err
	}
	var allTools []*mcp.Tool
	// This could be optimized by storing the combined list during buildMap
	for _, h := range c.handlers {
        tools, err := h.GetTools(ctx) // Call again or use stored list
        if err != nil {
            return nil, fmt.Errorf("failed to get tools from a handler during combined get: %w", err)
        }
        allTools = append(allTools, tools...)
    }
	return allTools, nil
}

func (c *CombinedToolHandler) CallTool(ctx context.Context, toolName string, input map[string]interface{}) (map[string]interface{}, error) {
	if err := c.buildMap(ctx); err != nil {
		return nil, err
	}
	handler, ok := c.toolMap[toolName]
	if !ok {
		return nil, fmt.Errorf("tool '%s' not found", toolName)
	}
	return handler.CallTool(ctx, toolName, input)
}

var _ ToolHandler = (*CombinedToolHandler)(nil)

*/
