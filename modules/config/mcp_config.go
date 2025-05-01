package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MCPServerConfig represents configuration for the MCP server
type MCPServerConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Port       string `mapstructure:"port"`
	Mode       string `mapstructure:"mode"` // "http" or "stdio"
	BasePath   string `mapstructure:"base_path"`
	EnableLogs bool   `mapstructure:"enable_logs"`
}

// MCPServerClientConfig represents a configuration for an external MCP client/server
type MCPServerClientConfig struct {
	Enabled      bool              `mapstructure:"enabled"`
	Type         string            `mapstructure:"type"`    // "http", "stdio"
	URL          string            `mapstructure:"url"`     // For HTTP clients
	Command      string            `mapstructure:"command"` // For stdio clients
	Args         []string          `mapstructure:"args"`    // For stdio clients
	Env          map[string]string `mapstructure:"env"`     // Environment variables
	Headers      map[string]string `mapstructure:"headers"` // HTTP headers
	AllowedTools []string          `mapstructure:"allowed_tools"`
	BlockedTools []string          `mapstructure:"blocked_tools"`
}

// GetMCPServerConfig returns the MCP server configuration
func (c *Config) GetMCPServerConfig() *MCPServerConfig {
	// Default configuration
	config := &MCPServerConfig{
		Enabled:    getBoolEnvVar("MCP_SERVER_ENABLED", true),
		Port:       getEnvVar("MCP_SERVER_PORT", "8085"),
		Mode:       getEnvVar("MCP_SERVER_MODE", "http"),
		BasePath:   getEnvVar("MCP_SERVER_BASE_PATH", "/mcp"),
		EnableLogs: getBoolEnvVar("MCP_SERVER_ENABLE_LOGS", true),
	}

	// If available in config, override defaults
	if c.MCPServer.Port != "" {
		config.Port = c.MCPServer.Port
	}
	if c.MCPServer.Mode != "" {
		config.Mode = c.MCPServer.Mode
	}
	if c.MCPServer.BasePath != "" {
		config.BasePath = c.MCPServer.BasePath
	}

	// Honor explicit disabled setting
	if !c.MCPServer.Enabled {
		config.Enabled = false
	}

	return config
}

// GetMCPClient returns the configuration for a specific MCP client
func (c *Config) GetMCPClient(name string) (*MCPServerClientConfig, error) {
	client, exists := c.MCPClients[name]
	if !exists {
		return nil, fmt.Errorf("MCP client '%s' not found in configuration", name)
	}

	return &client, nil
}

// Helper functions to get environment variables with defaults
func getEnvVar(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getBoolEnvVar(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// Initialize adds MCP-related fields to Config structure
func init() {
	// This is just to ensure compatibility with existing code
	// The actual structure is defined in config.go
}

// ValidateMCPClientConfig validates the MCP client configuration
func ValidateMCPClientConfig(config *MCPServerClientConfig) error {
	if config == nil {
		return fmt.Errorf("MCP client config is nil")
	}

	if !config.Enabled {
		return nil // If disabled, no need to validate further
	}

	config.Type = strings.ToLower(config.Type)

	switch config.Type {
	case "http":
		if config.URL == "" {
			return fmt.Errorf("URL is required for HTTP MCP clients")
		}
	case "stdio":
		if config.Command == "" {
			return fmt.Errorf("Command is required for stdio MCP clients")
		}
	default:
		return fmt.Errorf("unsupported MCP client type: %s", config.Type)
	}

	return nil
}
