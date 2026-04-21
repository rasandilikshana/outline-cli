package mcp

import (
	"context"

	"outline-cli/internal/api"
	"outline-cli/internal/sync"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type pullIn struct {
	Collection string `json:"collection" jsonschema:"collection name to pull"`
	LocalPath  string `json:"local_path" jsonschema:"absolute local directory path to write markdown files into"`
	DryRun     bool   `json:"dry_run,omitempty" jsonschema:"if true, preview without writing any files"`
}

type pushIn struct {
	Collection string `json:"collection" jsonschema:"target collection name (created if missing)"`
	FolderPath string `json:"folder_path" jsonschema:"absolute path of the local folder to upload"`
	DryRun     bool   `json:"dry_run,omitempty" jsonschema:"if true, preview without uploading"`
}

func registerSyncTools(s *sdk.Server, c *api.Client) {
	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_pull_collection",
		Description: "Download an Outline collection to a local folder as markdown files, preserving the document hierarchy. Use dry_run=true to preview.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in pullIn) (*sdk.CallToolResult, any, error) {
		result, err := sync.Pull(ctx, c, in.Collection, in.LocalPath, sync.PullOptions{DryRun: in.DryRun})
		if err != nil {
			return errResult(err)
		}
		errs := make([]string, 0, len(result.Errors))
		for _, e := range result.Errors {
			errs = append(errs, e.Error())
		}
		return textResult(map[string]any{
			"downloaded": result.Downloaded,
			"errors":     errs,
			"dry_run":    in.DryRun,
		})
	})

	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_push_folder",
		Description: "Upload a local folder of markdown files to an Outline collection, preserving directory hierarchy as nested documents. Existing documents with matching titles are updated. Use dry_run=true to preview.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in pushIn) (*sdk.CallToolResult, any, error) {
		result, err := sync.Push(ctx, c, in.Collection, in.FolderPath, sync.PushOptions{DryRun: in.DryRun})
		if err != nil {
			return errResult(err)
		}
		errs := make([]string, 0, len(result.Errors))
		for _, e := range result.Errors {
			errs = append(errs, e.Error())
		}
		return textResult(map[string]any{
			"created": result.Created,
			"updated": result.Updated,
			"skipped": result.Skipped,
			"errors":  errs,
			"dry_run": in.DryRun,
		})
	})
}
