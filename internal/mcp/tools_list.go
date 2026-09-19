package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// listTool определяет инструмент просмотра списка флагов
func (s *Server) listTool() mcp.Tool {
	return mcp.NewTool("mdflag_list",
		mcp.WithDescription(
			"List all feature flags in the project. Returns name, percentage, status, "+
				"hypothesis, and metrics for each flag. Use this to understand what "+
				"experiments are currently active before making changes.",
		),
	)
}

// handleList обрабатывает вызов инструмента списка флагов
func (s *Server) handleList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	flags, err := s.store.List()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list flags: %v", err)), nil
	}

	if len(flags) == 0 {
		return mcp.NewToolResultText(
			"No flags found in " + s.flagsDir + "\n\n" +
				"To create a flag, use the mdflag_create tool.",
		), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d flag(s) in %s:\n\n", len(flags), s.flagsDir))

	for _, f := range flags {
		sb.WriteString(fmt.Sprintf("### %s\n", f.Meta.Name))
		sb.WriteString(fmt.Sprintf("- Percentage: %d%%\n", f.Meta.Percentage))
		sb.WriteString(fmt.Sprintf("- Status: %s\n", f.Meta.Status))
		sb.WriteString(fmt.Sprintf("- Targeting: %s\n", f.Meta.Targeting))
		sb.WriteString(fmt.Sprintf("- Hypothesis: %s\n", f.Meta.Hypothesis))
		if f.Meta.Author != "" {
			sb.WriteString(fmt.Sprintf("- Author: %s\n", f.Meta.Author))
		}
		if len(f.Meta.Metrics) > 0 {
			sb.WriteString(fmt.Sprintf("- Metrics: %s\n", strings.Join(f.Meta.Metrics, ", ")))
		}
		if !f.Meta.Expires.IsZero() {
			sb.WriteString(fmt.Sprintf("- Expires: %s\n", f.Meta.Expires.Format("2006-01-02")))
		}
		sb.WriteString("\n")
	}

	return mcp.NewToolResultText(sb.String()), nil
}
