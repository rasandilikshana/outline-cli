package sync

import (
	"context"
	"fmt"
	"outline-cli/internal/api"
	"outline-cli/internal/models"
)

type PushOptions struct {
	DryRun bool
	Delete bool
}

type PushResult struct {
	Created int
	Updated int
	Skipped int
	Deleted int
	Errors  []error
}

func Push(ctx context.Context, client *api.Client, collectionName string, folderPath string, opts PushOptions) (*PushResult, error) {
	result := &PushResult{}

	localTree, err := BuildLocalTree(folderPath)
	if err != nil {
		return nil, fmt.Errorf("scan folder: %w", err)
	}

	coll, err := client.Collections.FindByName(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("find collection: %w", err)
	}

	if coll == nil {
		if opts.DryRun {
			fmt.Printf("[dry-run] Would create collection: %s\n", collectionName)
		} else {
			created, err := client.Collections.Create(ctx, models.CollectionCreateParams{
				Name: collectionName,
			})
			if err != nil {
				return nil, fmt.Errorf("create collection: %w", err)
			}
			coll = created
			fmt.Printf("Created collection: %s (%s)\n", coll.Name, coll.ID)
		}
	} else {
		fmt.Printf("Using existing collection: %s (%s)\n", coll.Name, coll.ID)
	}

	if opts.DryRun && coll == nil {
		printDryRunTree(localTree, 0)
		return result, nil
	}

	var remoteNodes []models.NavigationNode
	if coll != nil {
		remoteNodes, err = client.Collections.Documents(ctx, coll.ID)
		if err != nil {
			return nil, fmt.Errorf("get collection documents: %w", err)
		}
	}

	if err := pushChildren(ctx, client, coll.ID, "", localTree.Children, remoteNodes, opts, result); err != nil {
		return result, err
	}

	return result, nil
}

func pushChildren(ctx context.Context, client *api.Client, collectionID, parentDocID string, localNodes []*LocalNode, remoteNodes []models.NavigationNode, opts PushOptions, result *PushResult) error {
	for _, local := range localNodes {
		remoteDoc := findRemoteByTitle(local.Title, remoteNodes)

		if local.IsDir {
			var docID string
			if remoteDoc != nil {
				docID = remoteDoc.ID
				result.Skipped++
				if opts.DryRun {
					fmt.Printf("[dry-run] Skip existing folder doc: %s\n", local.Title)
				}
			} else {
				if opts.DryRun {
					fmt.Printf("[dry-run] Would create folder doc: %s\n", local.Title)
					result.Created++
				} else {
					params := models.DocumentCreateParams{
						Title:            local.Title,
						Text:             "",
						CollectionID:     collectionID,
						ParentDocumentID: parentDocID,
						Publish:          true,
					}
					doc, err := client.Documents.Create(ctx, params)
					if err != nil {
						result.Errors = append(result.Errors, fmt.Errorf("create %s: %w", local.Title, err))
						continue
					}
					docID = doc.ID
					result.Created++
					fmt.Printf("Created: %s\n", local.Title)
				}
			}

			var childRemoteNodes []models.NavigationNode
			if remoteDoc != nil {
				childRemoteNodes = remoteDoc.Children
			}
			if err := pushChildren(ctx, client, collectionID, docID, local.Children, childRemoteNodes, opts, result); err != nil {
				return err
			}
		} else {
			if remoteDoc != nil {
				if opts.DryRun {
					fmt.Printf("[dry-run] Would update: %s\n", local.Title)
				} else {
					_, err := client.Documents.Update(ctx, models.DocumentUpdateParams{
						ID:    remoteDoc.ID,
						Title: local.Title,
						Text:  local.Content,
					})
					if err != nil {
						result.Errors = append(result.Errors, fmt.Errorf("update %s: %w", local.Title, err))
						continue
					}
					fmt.Printf("Updated: %s\n", local.Title)
				}
				result.Updated++
			} else {
				if opts.DryRun {
					fmt.Printf("[dry-run] Would create: %s\n", local.Title)
				} else {
					params := models.DocumentCreateParams{
						Title:            local.Title,
						Text:             local.Content,
						CollectionID:     collectionID,
						ParentDocumentID: parentDocID,
						Publish:          true,
					}
					_, err := client.Documents.Create(ctx, params)
					if err != nil {
						result.Errors = append(result.Errors, fmt.Errorf("create %s: %w", local.Title, err))
						continue
					}
					fmt.Printf("Created: %s\n", local.Title)
				}
				result.Created++
			}
		}
	}
	return nil
}

func findRemoteByTitle(title string, nodes []models.NavigationNode) *models.NavigationNode {
	for i := range nodes {
		if nodes[i].Title == title {
			return &nodes[i]
		}
	}
	return nil
}

func printDryRunTree(node *LocalNode, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	kind := "doc"
	if node.IsDir {
		kind = "dir"
	}
	fmt.Printf("%s[%s] %s\n", indent, kind, node.Title)
	for _, child := range node.Children {
		printDryRunTree(child, depth+1)
	}
}
