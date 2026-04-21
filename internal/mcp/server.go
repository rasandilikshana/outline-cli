package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"outline-cli/internal/api"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const ServerName = "outline-cli"

func Run(ctx context.Context, client *api.Client, version string) error {
	server := sdk.NewServer(&sdk.Implementation{
		Name:    ServerName,
		Version: version,
	}, nil)

	registerDocumentTools(server, client)
	registerCollectionTools(server, client)
	registerCommentTools(server, client)
	registerRevisionTools(server, client)
	registerSyncTools(server, client)

	return server.Run(ctx, &sdk.StdioTransport{})
}

func textResult(v any) (*sdk.CallToolResult, any, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal result: %w", err)
	}
	return &sdk.CallToolResult{
		Content: []sdk.Content{&sdk.TextContent{Text: string(data)}},
	}, nil, nil
}

func errResult(err error) (*sdk.CallToolResult, any, error) {
	return &sdk.CallToolResult{
		IsError: true,
		Content: []sdk.Content{&sdk.TextContent{Text: err.Error()}},
	}, nil, nil
}
