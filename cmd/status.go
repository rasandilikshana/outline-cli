package cmd

import (
	"fmt"

	"outline-cli/internal/cli"
	"outline-cli/internal/config"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check connection status and instance info",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			cli.Output.Error("Config error: %s", err)
			return err
		}

		fmt.Printf("URL:      %s\n", cfg.URL)
		if cfg.APIKey != "" {
			fmt.Printf("API Key:  %s...%s\n", cfg.APIKey[:10], cfg.APIKey[len(cfg.APIKey)-4:])
		} else {
			fmt.Printf("API Key:  (not set)\n")
			return fmt.Errorf("API key not configured")
		}
		if cfg.Source != "" {
			fmt.Printf("Source:   %s\n", cfg.Source)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		info, err := client.Auth.Info(getContext())
		if err != nil {
			cli.Output.Error("Connection failed: %s", err)
			return err
		}

		fmt.Printf("Status:   Connected\n")
		fmt.Printf("User:     %s (%s)\n", info.User.Name, info.User.Email)
		fmt.Printf("Role:     %s\n", info.User.Role)
		fmt.Printf("Team:     %s\n", info.Team.Name)

		collections, _, err := client.Collections.List(getContext(), models.CollectionListParams{
			PaginationParams: models.PaginationParams{Limit: 1},
		})
		if err == nil {
			// Do a full count
			allColls, _, _ := client.Collections.List(getContext(), models.CollectionListParams{
				PaginationParams: models.PaginationParams{Limit: 100},
			})
			if allColls != nil {
				fmt.Printf("Collections: %d\n", len(allColls))
			} else {
				fmt.Printf("Collections: %d+\n", len(collections))
			}
		}

		cli.Output.Success("All checks passed.")
		return nil
	},
}
