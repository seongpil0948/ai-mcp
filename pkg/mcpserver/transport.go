package mcpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
)

// HTTPTransport implements a HTTP transport for the MCP server
type HTTPTransport struct {
	router *gin.Engine
	port   string
}

// NewHTTPTransport creates a new HTTP transport for the MCP server
func NewHTTPTransport(port string) *HTTPTransport {
	router := gin.Default()

	// Add CORS middleware
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

	return &HTTPTransport{
		router: router,
		port:   port,
	}
}

// Start starts the HTTP server
func (t *HTTPTransport) Start(ctx context.Context, handler server.RequestHandler) error {
	// Set up MCP HTTP endpoint
	t.router.POST("/mcp", func(c *gin.Context) {
		// Read request body
		buf, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read request body",
			})
			return
		}

		// Process the MCP request
		response, err := handler(buf)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return the response
		c.Data(http.StatusOK, "application/json", response)
	})

	// Health check endpoint
	t.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":" + t.port,
		Handler: t.router,
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	// Start the server
	log.Printf("Starting HTTP server on port %s", t.port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// StdioTransport implements a stdio transport for the MCP server
type StdioTransport struct{}

// NewStdioTransport creates a new stdio transport for the MCP server
func NewStdioTransport() *StdioTransport {
	return &StdioTransport{}
}

// Start starts the stdio transport
func (t *StdioTransport) Start(ctx context.Context, handler server.RequestHandler) error {
	// Set up signal handling
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Create stdio server
	stdioServer := server.NewStdioServer(handler)

	// Run the server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- stdioServer.Start()
	}()

	// Wait for termination signal or error
	select {
	case <-ctx.Done():
		log.Println("Context canceled, stopping stdio server")
		return ctx.Err()
	case <-signalCh:
		log.Println("Received termination signal, stopping stdio server")
		return nil
	case err := <-errCh:
		return err
	}
}
