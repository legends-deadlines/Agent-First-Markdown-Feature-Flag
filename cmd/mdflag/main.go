package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/legends-deadlines/mdflag/internal/cli"
	"github.com/legends-deadlines/mdflag/internal/mcp"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var err error

	switch os.Args[1] {
	case "create":
		err = cli.Create(os.Args[2:])

	case "rollout":
		err = cli.Rollout(os.Args[2:])

	case "verify":
		err = cli.Verify(os.Args[2:])

	case "list":
		err = cli.List(os.Args[2:])

	case "serve":
		// MCP-сервер для ИИ-агентов
		err = serveMCP(os.Args[2:])

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// serveMCP запускает MCP-сервер на stdio-транспорте
func serveMCP(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	dir := fs.String("dir", ".mdflag", "Директория для хранения флагов")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	srv, err := mcp.NewServer(*dir)
	if err != nil {
		return fmt.Errorf("init mcp server: %w", err)
	}

	return srv.Serve()
}

func printUsage() {
	fmt.Println(`MDFLAG — feature flags for AI agents

Usage:
  mdflag create --name <name> --hypothesis <text> [options]   Create new flag
  mdflag rollout --name <name> --percentage <0-100>           Change flag percentage
  mdflag verify [dir]                                          Check integrity of all flags
  mdflag list [dir]                                            List all flags
  mdflag serve [--dir <path>]                                  Start MCP server (stdio transport)

Create options:
  --name          Flag name (required)
  --percentage    Enable percentage 0-100 (default: 0)
  --targeting     Targeting method: user_id, session_id, random (default: user_id)
  --author        Author name (agent or human)
  --hypothesis    Experiment hypothesis (required)
  --metrics       Comma-separated metric names
  --expires       Expiration date in RFC3339 format
  --description   Human-readable description
  --dir           Flags directory (default: .mdflag)

Serve options:
  --dir           Flags directory (default: .mdflag)

Examples:
  mdflag create --name new-checkout --percentage 5 --hypothesis "Increase conversion by 5%"
  mdflag rollout --name new-checkout --percentage 25
  mdflag verify .mdflag/
  mdflag list .mdflag/
  mdflag serve --dir .mdflag/`)
}
