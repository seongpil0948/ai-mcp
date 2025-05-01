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

	"github.com/theshop/ai/modules/config"
	"github.com/theshop/ai/modules/integrations/figma"
	"github.com/theshop/ai/modules/integrations/gitlab"
	"github.com/theshop/ai/modules/integrations/notion"
	"github.com/theshop/ai/pkg/mcpserver"
)

var (
	configPath  = flag.String("config", "configs/config.yaml", "Path to config file")
	serverMode  = flag.String("mode", "", "Server mode (http or stdio, overrides config)")
	serverPort  = flag.String("port", "", "Server port for HTTP mode (overrides config)")
	debugMode   = flag.Bool("debug", false, "Enable debug logging")
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

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Apply command line overrides
	if *serverMode != "" {
		cfg.MCPServer.Mode = *serverMode
	}
	if *serverPort != "" {
		cfg.MCPServer.Port = *serverPort
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize MCP server
	mcpSrv := mcpserver.NewMCPServer(appName, appVersion)

	// Register integrations
	if err := registerIntegrations(ctx, mcpSrv, cfg); err != nil {
		log.Fatalf("Failed to register integrations: %v", err)
	}

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down", sig)
		cancel()
	}()

	// Start the server in the appropriate mode
	switch cfg.MCPServer.Mode {
	case "http":
		transport := mcpserver.NewHTTPTransport(cfg.MCPServer.Port)
		log.Printf("Starting HTTP server on port %s", cfg.MCPServer.Port)
		if err := mcpSrv.Start(ctx, transport); err != nil && err != context.Canceled {
			log.Fatalf("HTTP server error: %v", err)
		}
	case "stdio":
		transport := mcpserver.NewStdioTransport()
		log.Printf("Starting stdio server")
		if err := mcpSrv.Start(ctx, transport); err != nil && err != context.Canceled {
			log.Fatalf("Stdio server error: %v", err)
		}
	default:
		log.Fatalf("Unknown server mode: %s", cfg.MCPServer.Mode)
	}

	log.Println("Server shutdown complete")
}

// registerIntegrations registers all integrations with the MCP server
func registerIntegrations(ctx context.Context, mcpSrv *mcpserver.MCPServer, cfg *config.Config) error {
	// Register GitLab integration
	if gitlabClient, err := gitlab.NewClient(&cfg.GitLab); err == nil {
		gitlabService := gitlab.NewService(gitlabClient)
		if err := gitlabService.RegisterTools(mcpSrv); err != nil {
			log.Printf("Warning: Failed to register GitLab tools: %v", err)
		} else {
			log.Println("Registered GitLab integration")
		}
	} else {
		log.Printf("Warning: Failed to initialize GitLab client: %v", err)
	}

	// Register Notion integration
	if notionClient, err := notion.NewClient(&cfg.Notion); err == nil {
		notionService := notion.NewService(notionClient)
		if err := notionService.RegisterTools(mcpSrv); err != nil {
			log.Printf("Warning: Failed to register Notion tools: %v", err)
		} else {
			log.Println("Registered Notion integration")
		}
	} else {
		log.Printf("Warning: Failed to initialize Notion client: %v", err)
	}

	// Register Figma integration
	if figmaClient, err := figma.NewClient(&cfg.Figma); err == nil {
		figmaService := figma.NewService(figmaClient, cfg.Figma.TeamID)
		if err := figmaService.RegisterTools(mcpSrv); err != nil {
			log.Printf("Warning: Failed to register Figma tools: %v", err)
		} else {
			log.Println("Registered Figma integration")
		}
	} else {
		log.Printf("Warning: Failed to initialize Figma client: %v", err)
	}

	return nil
}
