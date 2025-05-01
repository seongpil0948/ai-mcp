package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/theshop/ai/internal/transport/http/handlers"
	"github.com/theshop/ai/modules/core"

	"github.com/gin-gonic/gin"
)

func main() {
	// 애플리케이션 초기화
	app, err := core.InitializeApp()
	if err != nil {
		log.Fatalf("애플리케이션 초기화 오류: %v", err)
	}

	// Gin 라우터 설정
	router := gin.Default()

	// CORS 미들웨어 설정
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 워크스페이스 핸들러 생성
	workspaceHandler := handlers.NewWorkspaceHandler(app.WorkspaceManager)

	// 통합 핸들러 생성
	integrationHandler := handlers.NewIntegrationHandler(
		app.WorkspaceManager,
		app.JiraClient,
		app.GitLabClient,
		app.IntegrationService,
	)

	// LLM 핸들러 생성
	llmHandler := handlers.NewLLMHandler(app.LLMService)

	// API 버전 v1
	v1 := router.Group("/api/v1")
	{
		// 워크스페이스 관련 API
		workspaces := v1.Group("/workspaces")
		{
			workspaces.GET("", workspaceHandler.ListWorkspaces)
			workspaces.POST("", workspaceHandler.CreateWorkspace)
			workspaces.GET("/:id", workspaceHandler.GetWorkspace)
			workspaces.PUT("/:id", workspaceHandler.UpdateWorkspace)
			workspaces.DELETE("/:id", workspaceHandler.DeleteWorkspace)

			// 워크스페이스 컨텍스트 내 Jira API
			workspaces.POST("/:id/jira/issues", integrationHandler.CreateJiraIssue)
			workspaces.POST("/:id/jira/search", integrationHandler.SearchJiraIssues)

			// 워크스페이스 컨텍스트 내 GitLab API
			workspaces.POST("/:id/gitlab/projects", integrationHandler.ListGitLabProjects)
			workspaces.POST("/:id/gitlab/mr-from-issue", integrationHandler.CreateMRFromIssue)

			// 기타 통합 엔드포인트...
		}

		// LLM 관련 API
		llmApi := v1.Group("/llm")
		{
			llmApi.POST("/generate", llmHandler.GenerateText)
			llmApi.POST("/analyze", llmHandler.AnalyzeData)
		}
	}

	// 상태 확인 엔드포인트
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"version": "1.0.0",
		})
	})

	// 서버 설정
	srv := &http.Server{
		Addr:    ":" + app.Config.Server.Port,
		Handler: router,
	}

	// 비동기로 서버 시작
	go func() {
		log.Printf("서버가 http://localhost:%s에서 시작되었습니다", app.Config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("서버 시작 오류: %v", err)
		}
	}()

	// 종료 시그널 처리
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("서버를 종료합니다...")

	// 서버 종료 컨텍스트 (10초 타임아웃)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 서버 우아하게 종료
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("서버 종료 오류: %v", err)
	}

	log.Println("서버가 종료되었습니다")
}
