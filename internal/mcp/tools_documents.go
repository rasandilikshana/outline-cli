package mcp

import (
	"context"

	"outline-cli/internal/api"
	"outline-cli/internal/models"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchIn struct {
	Query        string `json:"query" jsonschema:"search query text"`
	CollectionID string `json:"collection_id,omitempty" jsonschema:"optional collection ID to scope the search"`
	Limit        int    `json:"limit,omitempty" jsonschema:"max results (1-25, default 10)"`
	TitlesOnly   bool   `json:"titles_only,omitempty" jsonschema:"if true, only search document titles (faster)"`
}

type docHit struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Collection string  `json:"collection_id,omitempty"`
	URL        string  `json:"url,omitempty"`
	Snippet    string  `json:"snippet,omitempty"`
	Score      float64 `json:"score,omitempty"`
	Updated    string  `json:"updated,omitempty"`
}

type getDocIn struct {
	ID string `json:"id" jsonschema:"document ID or slug"`
}

type upsertDocIn struct {
	Title            string `json:"title" jsonschema:"document title (used to match existing doc in the collection)"`
	Text             string `json:"text" jsonschema:"markdown content"`
	CollectionID     string `json:"collection_id" jsonschema:"target collection ID"`
	ParentDocumentID string `json:"parent_document_id,omitempty" jsonschema:"optional parent document ID for nesting"`
	Publish          bool   `json:"publish,omitempty" jsonschema:"publish immediately (default true)"`
}

type archiveDocIn struct {
	ID string `json:"id" jsonschema:"document ID to archive"`
}

type listDocsIn struct {
	CollectionID string `json:"collection_id,omitempty" jsonschema:"filter by collection ID"`
	Limit        int    `json:"limit,omitempty" jsonschema:"max results (1-100, default 25)"`
}

func registerDocumentTools(s *sdk.Server, c *api.Client) {
	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_search",
		Description: "Search Outline documents by full-text query. Returns compact hits with id, title, snippet, score. Prefer titles_only=true for faster lookups when you only need to locate a document by name.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in searchIn) (*sdk.CallToolResult, any, error) {
		limit := in.Limit
		if limit <= 0 || limit > 25 {
			limit = 10
		}
		params := models.SearchParams{
			Query:            in.Query,
			CollectionID:     in.CollectionID,
			PaginationParams: models.PaginationParams{Limit: limit},
		}
		if in.TitlesOnly {
			docs, _, err := c.Search.Titles(ctx, params)
			if err != nil {
				return errResult(err)
			}
			hits := make([]docHit, 0, len(docs))
			for _, d := range docs {
				hits = append(hits, docHit{
					ID:         d.ID,
					Title:      d.Title,
					Collection: d.CollectionID,
					Updated:    d.UpdatedAt.Format("2006-01-02"),
				})
			}
			return textResult(map[string]any{"hits": hits})
		}
		results, _, err := c.Search.Documents(ctx, params)
		if err != nil {
			return errResult(err)
		}
		hits := make([]docHit, 0, len(results))
		for _, r := range results {
			hits = append(hits, docHit{
				ID:         r.Document.ID,
				Title:      r.Document.Title,
				Collection: r.Document.CollectionID,
				Snippet:    truncate(r.Context, 240),
				Score:      r.Ranking,
				Updated:    r.Document.UpdatedAt.Format("2006-01-02"),
			})
		}
		return textResult(map[string]any{"hits": hits})
	})

	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_get_document",
		Description: "Fetch full document content (markdown text, title, metadata) by ID. Use outline_search first to locate the ID.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in getDocIn) (*sdk.CallToolResult, any, error) {
		doc, err := c.Documents.Info(ctx, in.ID)
		if err != nil {
			return errResult(err)
		}
		out := map[string]any{
			"id":            doc.ID,
			"title":         doc.Title,
			"text":          doc.Text,
			"collection_id": doc.CollectionID,
			"revision":      doc.Revision,
			"updated_at":    doc.UpdatedAt.Format("2006-01-02 15:04"),
			"published_at":  nil,
		}
		if doc.PublishedAt != nil {
			out["published_at"] = doc.PublishedAt.Format("2006-01-02 15:04")
		}
		return textResult(out)
	})

	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_list_documents",
		Description: "List documents, optionally filtered by collection. Returns id+title+updated only — use outline_get_document to read content.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in listDocsIn) (*sdk.CallToolResult, any, error) {
		limit := in.Limit
		if limit <= 0 || limit > 100 {
			limit = 25
		}
		docs, _, err := c.Documents.List(ctx, models.DocumentListParams{
			CollectionID:     in.CollectionID,
			PaginationParams: models.PaginationParams{Limit: limit},
		})
		if err != nil {
			return errResult(err)
		}
		hits := make([]docHit, 0, len(docs))
		for _, d := range docs {
			hits = append(hits, docHit{
				ID:         d.ID,
				Title:      d.Title,
				Collection: d.CollectionID,
				Updated:    d.UpdatedAt.Format("2006-01-02"),
			})
		}
		return textResult(map[string]any{"documents": hits})
	})

	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_upsert_document",
		Description: "Create or update a document in a collection. Matches existing document by exact title within the collection; creates a new one if no match. Returns the resulting document ID and action taken.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in upsertDocIn) (*sdk.CallToolResult, any, error) {
		results, _, _ := c.Documents.Search(ctx, models.SearchParams{
			Query:            in.Title,
			CollectionID:     in.CollectionID,
			PaginationParams: models.PaginationParams{Limit: 10},
		})
		for _, r := range results {
			if r.Document.Title == in.Title && r.Document.CollectionID == in.CollectionID {
				doc, err := c.Documents.Update(ctx, models.DocumentUpdateParams{
					ID:    r.Document.ID,
					Title: in.Title,
					Text:  in.Text,
				})
				if err != nil {
					return errResult(err)
				}
				return textResult(map[string]any{
					"action":   "updated",
					"id":       doc.ID,
					"title":    doc.Title,
					"revision": doc.Revision,
				})
			}
		}
		publish := true
		if !in.Publish {
			publish = in.Publish
		}
		doc, err := c.Documents.Create(ctx, models.DocumentCreateParams{
			Title:            in.Title,
			Text:             in.Text,
			CollectionID:     in.CollectionID,
			ParentDocumentID: in.ParentDocumentID,
			Publish:          publish,
		})
		if err != nil {
			return errResult(err)
		}
		return textResult(map[string]any{
			"action":   "created",
			"id":       doc.ID,
			"title":    doc.Title,
			"revision": doc.Revision,
		})
	})

	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_archive_document",
		Description: "Archive (soft-delete) a document by ID. Use outline_restore_document to undo.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in archiveDocIn) (*sdk.CallToolResult, any, error) {
		doc, err := c.Documents.Archive(ctx, in.ID)
		if err != nil {
			return errResult(err)
		}
		return textResult(map[string]any{"id": doc.ID, "title": doc.Title, "status": "archived"})
	})
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
