package main

import (
	"fmt"
	"os"

	"github.com/legends-deadlines/mdflag/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcmd := os.Args[1]
	args := os.Args[2:]

	switch subcmd {
	case "create":
		if err := cli.Create(args); err != nil {
			fmt.Fprintf(os.Stderr, "create error: %v\n", err)
			os.Exit(1)
		}

	case "rollout":
		if err := cli.Rollout(args); err != nil {
			fmt.Fprintf(os.Stderr, "rollout error: %v\n", err)
			os.Exit(1)
		}

	case "verify":
		if err := cli.Verify(args); err != nil {
			fmt.Fprintf(os.Stderr, "verify error: %v\n", err)
			os.Exit(1)
		}

	case "list":
		if err := cli.List(args); err != nil {
			fmt.Fprintf(os.Stderr, "list error: %v\n", err)
			os.Exit(1)
		}

	case "serve":
		// Реализация в этапе 7 (MCP)
		fmt.Println("serve: MCP server start routine (stage 7)")

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", subcmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`MDFLAG — feature flags for AI agents

Usage:
  mdflag create <args>      Create new flag
  mdflag rollout <args>     Change flag percentage
  mdflag verify [dir]       Check integrity of all flags
  mdflag list [dir]         List all flags
  mdflag serve              Start MCP server (stdio transport)`)
}
