package mcp

import (
	"context"
	"fmt"

	"github.com/legends-deadlines/mdflag/internal/cli"
)

// Server handles Model Context Protocol stdio tools for AI coding agents.
type Server struct {
	StoreDir string
}

// NewServer initializes a new MCP server.
func NewServer(storeDir string) *Server {
	if storeDir == "" {
		storeDir = ".mdflag"
	}
	return &Server{StoreDir: storeDir}
}

// HandleCreateFlag handles the mdflag_create MCP tool invocation from an agent.
func (s *Server) HandleCreateFlag(name, hypothesis, description, author string, metrics []string) error {
	opts := cli.CreateOptions{
		Name:        name,
		Hypothesis:  hypothesis,
		Description: description,
		Author:      author,
		Targeting:   "user_id",
		Metrics:     metrics,
		Dir:         s.StoreDir,
	}
	return cli.CreateFlag(opts)
}

// Start launches the stdio transport listener for MCP.
func (s *Server) Start(ctx context.Context) error {
	fmt.Println("MDFLAG MCP Server running on stdio...")
	<-ctx.Done()
	return nil
}
