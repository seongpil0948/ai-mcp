package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ServerConfig 서버 관련 설정
type ServerConfig struct {
	Port    string        `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// JiraConfig Jira 연동 관련 설정
type JiraConfig struct {
	URL        string `mapstructure:"url"`
	Username   string `mapstructure:"username"`
	APIToken   string `mapstructure:"api_token"`
	ProjectKey string `mapstructure:"project_key"`
}

// GitLabConfig GitLab 연동 관련 설정
type GitLabConfig struct {
	BaseURL     string `mapstructure:"base_url"` // Added BaseURL
	Token       string `mapstructure:"token"`
	AccessToken string `mapstructure:"access_token"`
}

// ConfluenceConfig Confluence 연동 관련 설정
type ConfluenceConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	APIToken string `mapstructure:"api_token"`
}

// AWSConfig AWS 관련 설정
type AWSConfig struct {
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

// NotionConfig Notion 연동 관련 설정
type NotionConfig struct {
	Token  string `mapstructure:"token"`
	APIKey string `mapstructure:"api_key"`
}

// LLMConfig LLM 관련 설정
type LLMConfig struct {
	DefaultModel       string       `mapstructure:"default_model"`
	DefaultTemperature float64      `mapstructure:"default_temperature"`
	Provider           string       `mapstructure:"provider"`
	OpenAI             OpenAIConfig `mapstructure:"openai"`
	Claude             ClaudeConfig `mapstructure:"claude"`
	Gemini             GeminiConfig `mapstructure:"gemini"`
}

// OpenAIConfig OpenAI 설정
type OpenAIConfig struct {
	APIKey       string `mapstructure:"api_key"`
	OrgID        string `mapstructure:"org_id"`
	DefaultModel string `mapstructure:"default_model"`
}

// ClaudeConfig Claude 설정
type ClaudeConfig struct {
	APIKey       string `mapstructure:"api_key"`
	DefaultModel string `mapstructure:"default_model"`
}

// GeminiConfig Gemini 설정
type GeminiConfig struct {
	APIKey       string `mapstructure:"api_key"`
	DefaultModel string `mapstructure:"default_model"`
}

type Config struct {
	Server     ServerConfig                     `mapstructure:"server"`
	Jira       JiraConfig                       `mapstructure:"jira"`
	GitLab     GitLabConfig                     `mapstructure:"gitlab"`
	Confluence ConfluenceConfig                 `mapstructure:"confluence"`
	AWS        AWSConfig                        `mapstructure:"aws"`
	Notion     NotionConfig                     `mapstructure:"notion"`
	LLM        LLMConfig                        `mapstructure:"llm"`
	OpenAI     OpenAIConfig                     `mapstructure:"openai"`
	Claude     ClaudeConfig                     `mapstructure:"claude"`
	Gemini     GeminiConfig                     `mapstructure:"gemini"`
	Figma      FigmaConfig                      `mapstructure:"figma"`
	MCPServer  MCPServerConfig                  `mapstructure:"mcp_server"`
	MCPClients map[string]MCPServerClientConfig `mapstructure:"mcp_clients"`
}

// FigmaConfig Figma 연동 관련 설정
type FigmaConfig struct {
	AccessToken string `mapstructure:"access_token"`
	TeamID      string `mapstructure:"team_id"`
	ProjectID   string `mapstructure:"project_id"`
}

// NewConfig 설정 파일 로드 함수
func NewConfig() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.timeout", 30*time.Second)
	v.SetDefault("llm.provider", "openai")

	// Setup viper
	v.SetConfigName("config")
	v.SetConfigType("yaml") // or json, toml, etc.
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/theshop-ai/") // Example additional path
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // For environment variable mapping

	if err := v.ReadInConfig(); err != nil {
		// Ignore if config file not found, rely on defaults/env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		fmt.Println("Config file not found, using defaults and environment variables.")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Initialize MCP related fields if necessary (based on mcp_config.go logic)
	// cfg.Initialize() // Assuming an Initialize method exists

	return &cfg, nil
}
