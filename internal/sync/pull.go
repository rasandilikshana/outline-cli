package sync

import (
	"context"
	"fmt"
	"os"
	"outline-cli/internal/api"
	"outline-cli/internal/models"
	"path/filepath"
)

type PullOptions struct {
	DryRun bool
}

type PullResult struct {
	Downloaded int
	Errors     []error
}

func Pull(ctx context.Context, client *api.Client, collectionName string, localPath string, opts PullOptions) (*PullResult, error) {
	result := &PullResult{}

	coll, err := client.Collections.FindByName(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("find collection: %w", err)
	}
	if coll == nil {
		return nil, fmt.Errorf("collection not found: %s", collectionName)
	}

	fmt.Printf("Pulling collection: %s (%s)\n", coll.Name, coll.ID)

	nodes, err := client.Collections.Documents(ctx, coll.ID)
	if err != nil {
		return nil, fmt.Errorf("get document tree: %w", err)
	}

	if !opts.DryRun {
		if err := os.MkdirAll(localPath, 0755); err != nil {
			return nil, fmt.Errorf("create output dir: %w", err)
		}
	}

	if err := pullNodes(ctx, client, nodes, localPath, opts, result); err != nil {
		return result, err
	}

	return result, nil
}

func pullNodes(ctx context.Context, client *api.Client, nodes []models.NavigationNode, parentPath string, opts PullOptions, result *PullResult) error {
	for _, node := range nodes {
		doc, err := client.Documents.Info(ctx, node.ID)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("fetch %s: %w", node.Title, err))
			continue
		}

		if len(node.Children) > 0 {
			dirName := SanitizeFilename(node.Title)
			dirPath := filepath.Join(parentPath, dirName)

			if opts.DryRun {
				fmt.Printf("[dry-run] Would create dir: %s\n", dirPath)
			} else {
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					result.Errors = append(result.Errors, fmt.Errorf("mkdir %s: %w", dirPath, err))
					continue
				}
			}

			if doc.Text != "" {
				indexPath := filepath.Join(dirPath, "index.md")
				if opts.DryRun {
					fmt.Printf("[dry-run] Would write: %s\n", indexPath)
				} else {
					if err := os.WriteFile(indexPath, []byte(doc.Text), 0644); err != nil {
						result.Errors = append(result.Errors, fmt.Errorf("write %s: %w", indexPath, err))
					}
				}
				result.Downloaded++
			}

			if err := pullNodes(ctx, client, node.Children, dirPath, opts, result); err != nil {
				return err
			}
		} else {
			filename := SanitizeFilename(node.Title) + ".md"
			filePath := filepath.Join(parentPath, filename)

			if opts.DryRun {
				fmt.Printf("[dry-run] Would write: %s\n", filePath)
			} else {
				if err := os.WriteFile(filePath, []byte(doc.Text), 0644); err != nil {
					result.Errors = append(result.Errors, fmt.Errorf("write %s: %w", filePath, err))
					continue
				}
				fmt.Printf("Downloaded: %s\n", filePath)
			}
			result.Downloaded++
		}
	}
	return nil
}
