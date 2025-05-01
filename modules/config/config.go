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
	URL         string `mapstructure:"url"`
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
	DefaultModel       string  `mapstructure:"default_model"`
	DefaultTemperature float64 `mapstructure:"default_temperature"`
	Provider           string  `mapstructure:"provider"`
	OpenAI             OpenAIConfig `mapstructure:"openai"`
	Claude             ClaudeConfig `mapstructure:"claude"`
	Gemini             GeminiConfig `mapstructure:"gemini"`
}

// OpenAIConfig OpenAI 설정
type OpenAIConfig struct {
	APIKey      string `mapstructure:"api_key"`
	OrgID       string `mapstructure:"org_id"`
	DefaultModel string `mapstructure:"default_model"`
}

// ClaudeConfig Claude 설정
type ClaudeConfig struct {
	APIKey      string `mapstructure:"api_key"`
	DefaultModel string `mapstructure:"default_model"`
}

// GeminiConfig Gemini 설정
type GeminiConfig struct {
	APIKey      string `mapstructure:"api_key"`
	DefaultModel string `mapstructure:"default_model"`
}

// MCPServerConfig MCP 서버 설정
type MCPServerConfig struct {
	URL     string            `mapstructure:"url"`
	Type    string            `mapstructure:"type"`
	Command string            `mapstructure:"command"`
	Args    []string          `mapstructure:"args"`
	Env     map[string]string `mapstructure:"env"`
}

// Config 전체 애플리케이션 설정
type Config struct {
	Server     ServerConfig               `mapstructure:"server"`
	Jira       JiraConfig                 `mapstructure:"jira"`
	GitLab     GitLabConfig               `mapstructure:"gitlab"`
	Confluence ConfluenceConfig           `mapstructure:"confluence"`
	AWS        AWSConfig                  `mapstructure:"aws"`
	Notion     NotionConfig               `mapstructure:"notion"`
	LLM        LLMConfig                  `mapstructure:"llm"`
	MCPServers map[string]MCPServerConfig `mapstructure:"mcp_servers"`
}

// NewConfig 설정 파일 로드 함수
func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// 환경 변수 매핑 설정 (THESHOP_JIRA_URL -> jira.url)
	viper.SetEnvPrefix("THESHOP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 기본값 설정
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.timeout", "30s")
	viper.SetDefault("llm.default_model", "gpt-3.5-turbo")
	viper.SetDefault("llm.default_temperature", 0.7)
	viper.SetDefault("llm.provider", "openai")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("설정 파일 읽기 오류: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("설정 언마샬링 오류: %w", err)
	}

	return &config, nil
}
