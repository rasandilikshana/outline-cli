package mcp

import (
	"context"

	"outline-cli/internal/api"
	"outline-cli/internal/models"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listCollectionsIn struct {
	Limit int `json:"limit,omitempty" jsonschema:"max results (1-100, default 50)"`
}

type collectionSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Permission  string `json:"permission,omitempty"`
	Updated     string `json:"updated,omitempty"`
}

type getTreeIn struct {
	CollectionID string `json:"collection_id,omitempty" jsonschema:"collection ID (takes precedence over name)"`
	Name         string `json:"name,omitempty" jsonschema:"collection name lookup when ID is unknown"`
}

type treeNode struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Children []treeNode `json:"children,omitempty"`
}

func registerCollectionTools(s *sdk.Server, c *api.Client) {
	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_list_collections",
		Description: "List workspace collections. Returns id, name, description, permission. Use outline_get_collection_tree to explore the document hierarchy inside a collection.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in listCollectionsIn) (*sdk.CallToolResult, any, error) {
		limit := in.Limit
		if limit <= 0 || limit > 100 {
			limit = 50
		}
		colls, _, err := c.Collections.List(ctx, models.CollectionListParams{
			PaginationParams: models.PaginationParams{Limit: limit},
		})
		if err != nil {
			return errResult(err)
		}
		out := make([]collectionSummary, 0, len(colls))
		for _, cc := range colls {
			out = append(out, collectionSummary{
				ID:          cc.ID,
				Name:        cc.Name,
				Description: truncate(cc.Description, 160),
				Permission:  cc.Permission,
				Updated:     cc.UpdatedAt.Format("2006-01-02"),
			})
		}
		return textResult(map[string]any{"collections": out})
	})

	sdk.AddTool(s, &sdk.Tool{
		Name:        "outline_get_collection_tree",
		Description: "Return the nested document tree (id + title + children) for a collection. Accepts either collection_id or name.",
	}, func(ctx context.Context, req *sdk.CallToolRequest, in getTreeIn) (*sdk.CallToolResult, any, error) {
		collID := in.CollectionID
		if collID == "" {
			if in.Name == "" {
				return errResult(errMissing("collection_id or name"))
			}
			coll, err := c.Collections.FindByName(ctx, in.Name)
			if err != nil {
				return errResult(err)
			}
			if coll == nil {
				return errResult(errNotFound("collection", in.Name))
			}
			collID = coll.ID
		}
		nodes, err := c.Collections.Documents(ctx, collID)
		if err != nil {
			return errResult(err)
		}
		return textResult(map[string]any{
			"collection_id": collID,
			"tree":          convertNodes(nodes),
		})
	})
}

func convertNodes(nodes []models.NavigationNode) []treeNode {
	out := make([]treeNode, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, treeNode{
			ID:       n.ID,
			Title:    n.Title,
			Children: convertNodes(n.Children),
		})
	}
	return out
}
