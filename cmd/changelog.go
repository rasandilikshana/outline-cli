package cmd

import (
	"fmt"
	"outline-cli/internal/changelog"
	"outline-cli/internal/cli"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Generate and publish changelogs from git history",
}

var changelogGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate changelog markdown from git history",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		includeAuthors, _ := cmd.Flags().GetBool("include-authors")
		repoPath, _ := cmd.Flags().GetString("repo")

		if from == "" {
			return fmt.Errorf("--from is required (git ref: tag, commit, branch)")
		}
		if to == "" {
			to = "HEAD"
		}

		opts := changelog.Options{
			From:           from,
			To:             to,
			IncludeAuthors: includeAuthors,
			RepoPath:       repoPath,
		}

		md, err := changelog.Generate(opts)
		if err != nil {
			return err
		}

		fmt.Print(md)
		return nil
	},
}

var changelogPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Generate changelog and push to Outline as a document",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		collectionID, _ := cmd.Flags().GetString("collection")
		title, _ := cmd.Flags().GetString("title")
		parentID, _ := cmd.Flags().GetString("parent")
		includeAuthors, _ := cmd.Flags().GetBool("include-authors")
		repoPath, _ := cmd.Flags().GetString("repo")

		if from == "" {
			return fmt.Errorf("--from is required (git ref: tag, commit, branch)")
		}
		if to == "" {
			to = "HEAD"
		}
		if collectionID == "" {
			return fmt.Errorf("--collection is required")
		}
		if title == "" {
			title = fmt.Sprintf("Changelog %s..%s", from, to)
		}

		opts := changelog.Options{
			From:           from,
			To:             to,
			IncludeAuthors: includeAuthors,
			RepoPath:       repoPath,
		}

		md, err := changelog.Generate(opts)
		if err != nil {
			return err
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		// Check if doc with same title exists (upsert)
		results, _, _ := client.Documents.Search(getContext(), models.SearchParams{
			Query:        title,
			CollectionID: collectionID,
		})
		for _, r := range results {
			if r.Document.Title == title && r.Document.CollectionID == collectionID {
				_, err := client.Documents.Update(getContext(), models.DocumentUpdateParams{
					ID:   r.Document.ID,
					Text: md,
				})
				if err != nil {
					return fmt.Errorf("update existing document: %w", err)
				}
				cli.Output.Success("Updated existing document: %s (%s)", title, r.Document.ID)
				return nil
			}
		}

		doc, err := client.Documents.Create(getContext(), models.DocumentCreateParams{
			Title:            title,
			Text:             md,
			CollectionID:     collectionID,
			ParentDocumentID: parentID,
			Publish:          true,
		})
		if err != nil {
			return err
		}
		cli.Output.Success("Created changelog document: %s (%s)", doc.Title, doc.ID)
		return nil
	},
}

func init() {
	changelogCmd.AddCommand(changelogGenerateCmd)
	changelogCmd.AddCommand(changelogPushCmd)

	for _, cmd := range []*cobra.Command{changelogGenerateCmd, changelogPushCmd} {
		cmd.Flags().String("from", "", "Starting git ref (tag, commit, branch)")
		cmd.Flags().String("to", "HEAD", "Ending git ref")
		cmd.Flags().Bool("include-authors", false, "Include commit authors")
		cmd.Flags().String("repo", ".", "Path to git repository")
	}

	changelogPushCmd.Flags().String("collection", "", "Target collection ID")
	changelogPushCmd.Flags().String("title", "", "Document title (default: Changelog <from>..<to>)")
	changelogPushCmd.Flags().String("parent", "", "Parent document ID")
}
