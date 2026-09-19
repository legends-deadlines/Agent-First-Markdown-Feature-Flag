// Пакет mcp реализует MCP-сервер для MDFLAG.
// Сервер предоставляет инструменты для ИИ-агентов (Cursor, Claude Code, Windsurf).
//
// Транспорт: stdio (стандарт для локальных MCP-серверов).
// Среда агента сама поднимает процесс "mdflag serve" по конфигурации.
package mcp

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/legends-deadlines/mdflag/internal/flag"
)

func newToolResultError(msg string) *mcp.CallToolResult {
	res := mcp.NewToolResultText(msg)
	res.IsError = true
	return res
}

// Server представляет MCP-сервер для MDFLAG
type Server struct {
	store    *flag.Store
	flagsDir string
}

// NewServer создаёт новый MCP-сервер
func NewServer(flagsDir string) (*Server, error) {
	store, err := flag.NewStore(flagsDir)
	if err != nil {
		return nil, fmt.Errorf("init flag store: %w", err)
	}

	return &Server{
		store:    store,
		flagsDir: flagsDir,
	}, nil
}

// Serve запускает MCP-сервер на stdio-транспорте.
// Блокируется до завершения процесса.
func (s *Server) Serve() error {
	mcpServer := server.NewMCPServer(
		"mdflag",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// Регистрируем инструменты, доступные агенту
	mcpServer.AddTool(s.createTool(), s.handleCreate)
	mcpServer.AddTool(s.listTool(), s.handleList)
	mcpServer.AddTool(s.verifyTool(), s.handleVerify)

	// ВАЖНО: инструменты "rollout" и "delete" НЕ регистрируются.
	// Человек управляет флагами через CLI:
	//   mdflag rollout --name <flag> --percentage <0-100>
	// Это обеспечивает асимметрию прав: агент создаёт, человек управляет.

	// Запускаем stdio-сервер
	stdioServer := server.NewStdioServer(mcpServer)
	return stdioServer.Listen(context.Background(), os.Stdin, os.Stdout)
}
