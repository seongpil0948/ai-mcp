package mcpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/protobuf/types/known/structpb"
)

// ToolHandler defines the interface required by the MCP server
// to interact with the actual tool implementations.
type ToolHandler interface {
	GetTools(ctx context.Context) ([]*mcp.Tool, error)
	CallTool(ctx context.Context, toolName string, input map[string]interface{}) (map[string]interface{}, error)
}

// Server wraps the mcp-go server implementation and acts as the RequestHandler.
type Server struct {
	handler     ToolHandler
	httpServer  *http.Server // Store the standard Go http server
	stdioServer *server.StdioServer
	// Add other server types if needed
}

// New creates a new MCP Server wrapper.
func New(handler ToolHandler) *Server {
	return &Server{
		handler: handler,
	}
}

// RegisterHTTPTransport configures the HTTP transport.
// It prepares an http.Server but doesn't start listening yet.
func (s *Server) RegisterHTTPTransport(addr string, path string) error {
	mcpHandler, err := server.NewHTTPHandler(s) // Pass the wrapper which implements server.RequestHandler
	if err != nil {
		return fmt.Errorf("failed to create mcp http handler: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle(path+"/", http.StripPrefix(path, mcpHandler)) // Handle /tools and /calls under the path

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	log.Printf("MCP HTTP transport configured for %s%s", addr, path)
	return nil
}

// RegisterStdioTransport configures the Stdio transport.
func (s *Server) RegisterStdioTransport() error {
	stdioServer, err := server.NewStdioServer(s) // Pass the wrapper
	if err != nil {
		return fmt.Errorf("failed to create stdio server: %w", err)
	}
	s.stdioServer = stdioServer
	log.Println("MCP Stdio transport registered")
	return nil
}

// Run starts the registered transports. This is a blocking call.
func (s *Server) Run() error {
	errCh := make(chan error, 2) // Channel to collect errors from goroutines

	// Start HTTP server if configured
	if s.httpServer != nil {
		go func() {
			log.Printf("Starting MCP HTTP server on %s", s.httpServer.Addr)
			if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP server error: %v", err)
				errCh <- fmt.Errorf("http server error: %w", err)
			} else {
				log.Println("HTTP server stopped.")
				errCh <- nil // Indicate clean shutdown or no error
			}
		}()
	}

	// Start Stdio server if configured
	if s.stdioServer != nil {
		go func() {
			log.Println("Starting MCP Stdio server")
			if err := s.stdioServer.Serve(); err != nil {
				log.Printf("Stdio server error: %v", err)
				errCh <- fmt.Errorf("stdio server error: %w", err)
			} else {
				log.Println("Stdio server stopped.")
				errCh <- nil // Indicate clean shutdown or no error
			}
		}()
	}

	// Wait for errors from the running servers
	var firstErr error
	serversRunning := 0
	if s.httpServer != nil {
		serversRunning++
	}
	if s.stdioServer != nil {
		serversRunning++
	}

	if serversRunning == 0 {
		return fmt.Errorf("no MCP transports registered or started")
	}

	for i := 0; i < serversRunning; i++ {
		err := <-errCh
		if err != nil && firstErr == nil {
			firstErr = err // Capture the first error encountered
		}
	}

	return firstErr // Return the first error encountered, or nil if all stopped cleanly
}

// Stop gracefully stops the server transports.
func (s *Server) Stop(ctx context.Context) error {
	log.Println("Stopping MCP server...")
	var firstErr error

	if s.httpServer != nil {
		log.Println("Shutting down HTTP server...")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	if s.stdioServer != nil {
		log.Println("Stopping Stdio server...")
		// StdioServer might need a specific Stop method or rely on closing stdin/stdout
		// Assuming it stops when Serve() returns, which happens on EOF or error.
		// If it has a Stop method:
		// if err := s.stdioServer.Stop(); err != nil {
		//     log.Printf("Stdio server stop error: %v", err)
		//     if firstErr == nil { firstErr = err }
		// }
	}
	return firstErr
}

// --- Implement server.RequestHandler interface ---

func (s *Server) GetTools(ctx context.Context) ([]*mcp.Tool, error) {
	log.Println("Handler: GetTools called")
	tools, err := s.handler.GetTools(ctx)
	if err != nil {
		log.Printf("Handler: GetTools error: %v", err)
	}
	return tools, err
}

func (s *Server) CallTool(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	toolName := req.GetToolName()
	log.Printf("Handler: CallTool called for tool: %s", toolName)

	var inputMap map[string]interface{}
	if req.GetInput() != nil {
		inputMap = req.GetInput().AsMap()
	}

	resultData, err := s.handler.CallTool(ctx, toolName, inputMap)
	if err != nil {
		log.Printf("Handler: CallTool error for %s: %v", toolName, err)
		errorStruct, _ := structpb.NewStruct(map[string]interface{}{"message": err.Error()})
		return &mcp.CallToolResult{
			ToolName: toolName,
			Result:   &mcp.CallToolResult_Error{Error: errorStruct},
		}, nil // Return nil Go error, error is in the result payload
	}

	resultStruct, err := structpb.NewStruct(resultData)
	if err != nil {
		log.Printf("Handler: Failed to convert result data to struct for %s: %v", toolName, err)
		errorStruct, _ := structpb.NewStruct(map[string]interface{}{"message": fmt.Sprintf("Internal error processing tool result: %v", err)})
		return &mcp.CallToolResult{
			ToolName: toolName,
			Result:   &mcp.CallToolResult_Error{Error: errorStruct},
		}, nil
	}

	log.Printf("Handler: CallTool success for %s", toolName)
	return &mcp.CallToolResult{
		ToolName: toolName,
		Result:   &mcp.CallToolResult_Content{Content: resultStruct},
	}, nil
}

// Ensure *Server implements server.RequestHandler
var _ server.RequestHandler = (*Server)(nil)
