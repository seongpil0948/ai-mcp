package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig는 애플리케이션 설정을 로드하는 함수입니다.
// 설정 파일(config.yaml), 환경 변수, 기본값을 순차적으로 참조합니다.
func LoadConfig() (*Config, error) {
	// Viper 인스턴스 생성
	v := viper.New()

	// 설정 파일 이름, 경로, 포맷 지정
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// 환경 변수 자동 로드 설정
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 기본값 설정
	setDefaultConfig(v)

	// 설정 파일 읽기 (파일이 없어도 에러를 반환하지 않음)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("설정 파일 읽기 오류: %w", err)
		}
		fmt.Println("설정 파일을 찾을 수 없습니다. 기본값과 환경 변수를 사용합니다.")
	}

	// 설정 구조체로 언마샬링
	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("설정 언마샬링 오류: %w", err)
	}

	// 환경 변수에서 민감 정보 로드 (환경 변수가 우선)
	loadSensitiveInfoFromEnv(config)

	return config, nil
}

// setDefaultConfig는 기본 설정값을 지정합니다.
func setDefaultConfig(v *viper.Viper) {
	// 서버 기본 설정
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.timeout", "30s")
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.idle_timeout", "60s")

	// AWS 기본 설정
	v.SetDefault("aws.region", "ap-northeast-2")

	// S3 기본 설정
	v.SetDefault("s3.prefix", "github.com/theshop/ai/results")

	// 파일 시스템 기본 설정
	v.SetDefault("filesystem.max_upload_size", 10*1024*1024) // 10MB

	// LLM 기본 설정
	v.SetDefault("llm.provider", "openai")
	v.SetDefault("llm.openai.model", "gpt-3.5-turbo")
	v.SetDefault("llm.claude.model", "claude-3-haiku-20240307")
	v.SetDefault("llm.gemini.model", "gemini-pro")

	// 로그 레벨 기본 설정
	v.SetDefault("log_level", "info")
}

// loadSensitiveInfoFromEnv는 환경 변수에서 민감 정보를 로드합니다.
func loadSensitiveInfoFromEnv(config *Config) {
	// Jira 관련 민감 정보
	if envVal := os.Getenv("JIRA_USERNAME"); envVal != "" {
		config.Jira.Username = envVal
	}
	if envVal := os.Getenv("JIRA_API_TOKEN"); envVal != "" {
		config.Jira.APIToken = envVal
	}

	// Confluence 관련 민감 정보
	if envVal := os.Getenv("CONFLUENCE_USERNAME"); envVal != "" {
		config.Confluence.Username = envVal
	}
	if envVal := os.Getenv("CONFLUENCE_API_TOKEN"); envVal != "" {
		config.Confluence.APIToken = envVal
	}

	// AWS 관련 민감 정보
	if envVal := os.Getenv("AWS_ACCESS_KEY_ID"); envVal != "" {
		config.AWS.AccessKeyID = envVal
	}
	if envVal := os.Getenv("AWS_SECRET_ACCESS_KEY"); envVal != "" {
		config.AWS.SecretAccessKey = envVal
	}

	// GitLab 관련 민감 정보
	if envVal := os.Getenv("GITLAB_ACCESS_TOKEN"); envVal != "" {
		config.GitLab.AccessToken = envVal
	}

	// Notion 관련 민감 정보
	if envVal := os.Getenv("NOTION_API_KEY"); envVal != "" {
		config.Notion.APIKey = envVal
	}

	// LLM 관련 민감 정보
	if envVal := os.Getenv("OPENAI_API_KEY"); envVal != "" {
		config.LLM.OpenAI.APIKey = envVal
	}
	if envVal := os.Getenv("CLAUDE_API_KEY"); envVal != "" {
		config.LLM.Claude.APIKey = envVal
	}
	if envVal := os.Getenv("GEMINI_API_KEY"); envVal != "" {
		config.LLM.Gemini.APIKey = envVal
	}
}
