package mcp

import (
	"context"

	"outline-cli/internal/api"
	"outline-cli/internal/models"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listRevisionsIn struct {
	DocumentID string `json:"document_id" jsonschema:"document ID to list revisions for"`
	Limit      int    `json:"limit,omitempty" jsonschema:"max results (1-50, default 10)"`
}

type revisionSummary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Version   int    `json:"version"`
	Author    string `json:"author,omitempty"`
	CreatedAt string `json:"created_at"`
}

func registerRevisionTools(s *sdk.Server, c *api.Client) {
	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_list_revisions",
		Description: "List revisions for a document with version numbers and authors. Does not return revision text — use the Outline UI or a targeted API call if diffing content is needed.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in listRevisionsIn) (*sdk.CallToolResult, any, error) {
		limit := in.Limit
		if limit <= 0 || limit > 50 {
			limit = 10
		}
		revs, _, err := c.Revisions.List(ctx, models.RevisionListParams{
			DocumentID:       in.DocumentID,
			PaginationParams: models.PaginationParams{Limit: limit},
		})
		if err != nil {
			return errResult(err)
		}
		out := make([]revisionSummary, 0, len(revs))
		for _, r := range revs {
			author := ""
			if r.CreatedBy != nil {
				author = r.CreatedBy.Name
			}
			out = append(out, revisionSummary{
				ID:        r.ID,
				Title:     r.Title,
				Version:   r.Version,
				Author:    author,
				CreatedAt: r.CreatedAt.Format("2006-01-02 15:04"),
			})
		}
		return textResult(map[string]any{"revisions": out})
	})
}
