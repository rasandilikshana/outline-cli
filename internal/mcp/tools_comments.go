package mcp

import (
	"context"

	"outline-cli/internal/api"
	"outline-cli/internal/models"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listCommentsIn struct {
	DocumentID string `json:"document_id" jsonschema:"document ID to list comments for"`
	Limit      int    `json:"limit,omitempty" jsonschema:"max results (1-100, default 25)"`
}

type createCommentIn struct {
	DocumentID      string `json:"document_id" jsonschema:"document ID to comment on"`
	Text            string `json:"text" jsonschema:"comment text (plain text)"`
	ParentCommentID string `json:"parent_comment_id,omitempty" jsonschema:"optional parent comment ID for threaded replies"`
}

type commentSummary struct {
	ID        string `json:"id"`
	Author    string `json:"author,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	ParentID  string `json:"parent_comment_id,omitempty"`
}

func registerCommentTools(s *sdk.Server, c *api.Client) {
	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_list_comments",
		Description: "List comments on a document. Returns author, created_at, and thread parent. Use outline_get_document for the doc body.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in listCommentsIn) (*sdk.CallToolResult, any, error) {
		limit := in.Limit
		if limit <= 0 || limit > 100 {
			limit = 25
		}
		comments, _, err := c.Comments.List(ctx, models.CommentListParams{
			DocumentID:       in.DocumentID,
			PaginationParams: models.PaginationParams{Limit: limit},
		})
		if err != nil {
			return errResult(err)
		}
		out := make([]commentSummary, 0, len(comments))
		for _, cm := range comments {
			author := ""
			if cm.CreatedBy != nil {
				author = cm.CreatedBy.Name
			}
			parent := ""
			if cm.ParentCommentID != nil {
				parent = *cm.ParentCommentID
			}
			out = append(out, commentSummary{
				ID:        cm.ID,
				Author:    author,
				CreatedAt: cm.CreatedAt.Format("2006-01-02 15:04"),
				ParentID:  parent,
			})
		}
		return textResult(map[string]any{"comments": out})
	})

	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_create_comment",
		Description: "Add a comment to a document. Pass parent_comment_id to reply within an existing thread.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in createCommentIn) (*sdk.CallToolResult, any, error) {
		cm, err := c.Comments.Create(ctx, models.CommentCreateParams{
			DocumentID:      in.DocumentID,
			ParentCommentID: in.ParentCommentID,
			Text:            in.Text,
		})
		if err != nil {
			return errResult(err)
		}
		return textResult(map[string]any{
			"id":          cm.ID,
			"document_id": cm.DocumentID,
			"created_at":  cm.CreatedAt.Format("2006-01-02 15:04"),
		})
	})
}
