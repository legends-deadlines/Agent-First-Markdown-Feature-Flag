package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/legends-deadlines/mdflag/internal/flag"
)

// createTool определяет инструмент создания флага
func (s *Server) createTool() mcp.Tool {
	return mcp.NewTool("mdflag_create",
		mcp.WithDescription(
			"Create a new feature flag for experimental code changes. "+
				"The flag is stored as a Markdown file in the .mdflag/ directory. "+
				"IMPORTANT: After creating a flag, wrap your new code with a check "+
				"using the runtime library. The flag starts disabled by default (percentage: 0). "+
				"A human reviewer will gradually enable it.",
		),
		mcp.WithString("name",
			mcp.Description("Unique name for the flag. Use kebab-case (e.g., 'new-checkout-flow'). Must not match an existing flag name."),
			mcp.Required(),
		),
		mcp.WithString("hypothesis",
			mcp.Description("What you expect this change to achieve. Example: 'Increase checkout conversion by 5%'. This helps reviewers evaluate the experiment."),
			mcp.Required(),
		),
		mcp.WithNumber("percentage",
			mcp.Description("Initial enable percentage (0-100). Default: 0 (disabled). Recommended: start at 0 and let a human increase it gradually."),
		),
		mcp.WithString("targeting",
			mcp.Description("Targeting method: 'user_id' (same user always gets same result), 'session_id' (per session), or 'random'. Default: 'user_id'."),
		),
		mcp.WithString("author",
			mcp.Description("Your identifier as the flag creator. Example: 'claude-code-agent' or 'cursor-agent'."),
		),
		mcp.WithString("description",
			mcp.Description("Human-readable description of what code changes this flag controls. Include affected files and the reason for the change."),
		),
		mcp.WithString("metrics",
			mcp.Description("Comma-separated list of metrics to track for this experiment. Example: 'conversion_rate,checkout_time_seconds,error_rate'."),
		),
		mcp.WithString("expires",
			mcp.Description("Expiration date in RFC3339 format. After this date the flag is automatically disabled. Example: '2026-12-17T00:00:00Z'. Recommended: set to 2-4 weeks from now."),
		),
	)
}

// handleCreate обрабатывает вызов инструмента создания флага
func (s *Server) handleCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments

	// Извлекаем обязательные параметры
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return newToolResultError("parameter 'name' is required and must be a non-empty string"), nil
	}

	hypothesis, ok := args["hypothesis"].(string)
	if !ok || hypothesis == "" {
		return newToolResultError("parameter 'hypothesis' is required and must be a non-empty string"), nil
	}

	// Извлекаем опциональные параметры со значениями по умолчанию
	percentage := 0
	if p, ok := args["percentage"].(float64); ok {
		percentage = int(p)
	}

	targeting := "user_id"
	if t, ok := args["targeting"].(string); ok && t != "" {
		targeting = t
	}

	author := ""
	if a, ok := args["author"].(string); ok {
		author = a
	}

	description := ""
	if d, ok := args["description"].(string); ok {
		description = d
	}

	metricsStr := ""
	if m, ok := args["metrics"].(string); ok {
		metricsStr = m
	}

	expiresStr := ""
	if e, ok := args["expires"].(string); ok {
		expiresStr = e
	}

	// Валидация процента
	if percentage < 0 || percentage > 100 {
		return newToolResultError(fmt.Sprintf("percentage must be 0-100, got %d", percentage)), nil
	}

	// Валидация таргетинга
	if !flag.IsValidTargeting(targeting) {
		return newToolResultError(fmt.Sprintf("invalid targeting %q: must be 'user_id', 'session_id', or 'random'", targeting)), nil
	}

	// Парсим метрики
	var metricList []string
	if metricsStr != "" {
		for _, m := range strings.Split(metricsStr, ",") {
			trimmed := strings.TrimSpace(m)
			if trimmed != "" {
				metricList = append(metricList, trimmed)
			}
		}
	}

	// Парсим срок действия
	var expiresTime time.Time
	if expiresStr != "" {
		var err error
		expiresTime, err = time.Parse(time.RFC3339, expiresStr)
		if err != nil {
			return newToolResultError(fmt.Sprintf("invalid expires format %q: expected RFC3339 (e.g., 2026-12-17T00:00:00Z)", expiresStr)), nil
		}
	}

	// Формируем мета-данные флага
	meta := flag.FlagMeta{
		Name:       name,
		Percentage: percentage,
		Targeting:  targeting,
		Status:     flag.StatusActive,
		Created:    time.Now(),
		Author:     author,
		Expires:    expiresTime,
		Hypothesis: hypothesis,
		Metrics:    metricList,
	}

	// Формируем Markdown-тело
	var bodyText string
	if description != "" {
		bodyText = "## Description\n" + description + "\n"
	}

	// Создаём флаг через хранилище
	f, err := s.store.Create(meta, bodyText)
	if err != nil {
		return newToolResultError(fmt.Sprintf("failed to create flag: %v", err)), nil
	}

	// Формируем сообщение об успехе
	result := fmt.Sprintf(
		"Flag '%s' created successfully.\n\n"+
			"File: %s/%s.md\n"+
			"Percentage: %d%%\n"+
			"Targeting: %s\n"+
			"Hypothesis: %s\n\n"+
			"NEXT STEPS:\n"+
			"1. Wrap your new code with a runtime check:\n"+
			"   if client.Enabled(\"%s\", entityID) { /* new code */ }\n"+
			"2. Commit both the flag file and your code changes.\n"+
			"3. A human reviewer will gradually enable the flag using 'mdflag rollout'.\n\n"+
			"IMPORTANT: Do NOT modify the flag file directly. Integrity hashes protect it.",
		f.Meta.Name, s.flagsDir, f.Meta.Name,
		f.Meta.Percentage, f.Meta.Targeting, f.Meta.Hypothesis,
		f.Meta.Name,
	)

	return mcp.NewToolResultText(result), nil
}
