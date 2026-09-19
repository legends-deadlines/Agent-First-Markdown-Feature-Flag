package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/legends-deadlines/mdflag/internal/flag"
)

// verifyTool определяет инструмент проверки целостности флагов
func (s *Server) verifyTool() mcp.Tool {
	return mcp.NewTool("mdflag_verify",
		mcp.WithDescription(
			"Verify integrity of all feature flags. Checks that no flag files have been "+
				"tampered with by comparing stored hashes against computed hashes. "+
				"Use this after pulling changes from remote to ensure flags are intact.",
		),
	)
}

// handleVerify обрабатывает вызов инструмента проверки целостности
func (s *Server) handleVerify(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	results, err := flag.ValidateDir(s.flagsDir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to validate flags: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No flags found in " + s.flagsDir), nil
	}

	validCount := 0
	invalidCount := 0
	var sb strings.Builder

	for path, result := range results {
		if result.IsValid() {
			validCount++
			sb.WriteString(fmt.Sprintf("[OK]   %s\n", path))
		} else {
			invalidCount++
			sb.WriteString(fmt.Sprintf("[FAIL] %s\n", path))
			for _, errMsg := range result.Errors {
				sb.WriteString(fmt.Sprintf("       - %s\n", errMsg))
			}
		}
	}

	sb.WriteString(fmt.Sprintf("\nResults: %d valid, %d invalid, %d total\n", validCount, invalidCount, validCount+invalidCount))

	if invalidCount > 0 {
		sb.WriteString("\nWARNING: Integrity violations detected. Flag files may have been modified outside of MDFLAG tools.\n")
		sb.WriteString("Do not trust the affected flags until the issue is resolved.\n")
		return mcp.NewToolResultText(sb.String()), nil
	}

	sb.WriteString("\nAll flags are intact. No integrity violations detected.\n")
	return mcp.NewToolResultText(sb.String()), nil
}
