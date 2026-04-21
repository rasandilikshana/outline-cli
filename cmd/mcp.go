package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	outlinemcp "outline-cli/internal/mcp"

	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run an MCP server exposing Outline to MCP-capable clients",
	Long: `Start a Model Context Protocol server over stdio.

Configure your MCP client (e.g. Claude Code) to launch this binary with:
    command: outline
    args:    ["mcp"]

Credentials are read from the same config as the rest of the CLI:
  ~/.outline-cli/config.yaml  OR  OUTLINE_URL / OUTLINE_API_KEY env vars.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		return outlinemcp.Run(ctx, client, Version)
	},
}
