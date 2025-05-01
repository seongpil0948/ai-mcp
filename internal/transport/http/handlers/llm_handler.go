package handlers

import (
	"io"
	"net/http"

	"github.com/theshop/ai/modules/integrations/llm"

	"github.com/gin-gonic/gin"
)

// LLMHandler LLM API 핸들러
type LLMHandler struct {
	llmService llm.Service
}

// NewLLMHandler 새 LLM 핸들러 생성
func NewLLMHandler(llmService llm.Service) *LLMHandler {
	return &LLMHandler{
		llmService: llmService,
	}
}

// GenerateTextRequest 텍스트 생성 요청
type GenerateTextRequest struct {
	Prompt      string  `json:"prompt" binding:"required"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	System      string  `json:"system"`
}

// GenerateTextResponse 텍스트 생성 응답
type GenerateTextResponse struct {
	Text string `json:"text"`
}

// GenerateText 텍스트 생성 핸들러
func (h *LLMHandler) GenerateText(c *gin.Context) {
	var req GenerateTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "유효하지 않은 요청 데이터: " + err.Error(),
		})
		return
	}

	options := llm.GenerateOptions{
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		System:      req.System,
	}

	generatedText, err := h.llmService.GenerateText(c, req.Prompt, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "텍스트 생성 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GenerateTextResponse{
		Text: generatedText,
	})
}

// AnalyzeDataRequest 데이터 분석 요청
type AnalyzeDataRequest struct {
	Query  string `json:"query" binding:"required"`
	Model  string `json:"model"`
	System string `json:"system"`
}

// AnalyzeDataResponse 데이터 분석 응답
type AnalyzeDataResponse struct {
	Result string `json:"result"`
}

// AnalyzeData 데이터 분석 핸들러
func (h *LLMHandler) AnalyzeData(c *gin.Context) {
	// 파일 처리
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "파일을 찾을 수 없음: " + err.Error(),
		})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "파일 읽기 실패: " + err.Error(),
		})
		return
	}

	// 쿼리 필드 처리
	query := c.PostForm("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "쿼리가 필요합니다",
		})
		return
	}

	// 데이터 분석 실행
	result, err := h.llmService.AnalyzeData(c, fileBytes, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "데이터 분석 실패: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AnalyzeDataResponse{
		Result: result,
	})
}
