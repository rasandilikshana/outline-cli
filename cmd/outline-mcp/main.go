// outline-mcp is a standalone Model Context Protocol server for Outline.
//
// It exposes the same tool surface as `outline mcp` but ships as its own
// binary with a smaller footprint (no cobra, no bubbletea/lipgloss TUI).
// Intended for users who only want MCP and not the full outline-cli.
//
// Credentials are loaded from the same sources as outline-cli:
//   1. OUTLINE_URL / OUTLINE_API_KEY environment variables
//   2. ./.outline/config.yaml  (project-scoped)
//   3. ./outline.yaml          (project-scoped alternative)
//   4. ~/.outline-cli/config.yaml (home)
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"outline-cli/internal/api"
	"outline-cli/internal/config"
	outlinemcp "outline-cli/internal/mcp"
)

// Version is injected at build time via -ldflags.
var Version = "0.3.0"

const usage = `outline-mcp — Model Context Protocol server for Outline

Usage:
  outline-mcp            Run the MCP server over stdio
  outline-mcp --version  Print version and exit
  outline-mcp --help     Print this help and exit

Credentials (any one of these):
  OUTLINE_URL / OUTLINE_API_KEY   env vars (preferred for MCP clients)
  ./.outline/config.yaml          project-scoped config
  ./outline.yaml                  project-root alternative
  ~/.outline-cli/config.yaml      personal default

Register with Claude Code:
  claude mcp add outline --scope project outline-mcp \
    --env OUTLINE_URL=https://outline.example.com \
    --env OUTLINE_API_KEY=ol_api_your_key
`

func main() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--version", "-v":
			fmt.Printf("outline-mcp v%s\n", Version)
			return
		case "--help", "-h":
			fmt.Print(usage)
			return
		default:
			fmt.Fprintf(os.Stderr, "unknown argument: %s\n\n%s", arg, usage)
			os.Exit(2)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %s\n", err)
		os.Exit(2)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(2)
	}

	client := api.NewClient(cfg.URL, cfg.APIKey)
	client.HTTPClient.Timeout = 30 * time.Second

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := outlinemcp.Run(ctx, client, Version); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %s\n", err)
		os.Exit(1)
	}
}
