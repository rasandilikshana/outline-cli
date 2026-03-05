package cmd

import (
	"fmt"

	"outline-cli/internal/cli"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var authTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test API credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		info, err := client.Auth.Info(getContext())
		if err != nil {
			cli.Output.Error("Authentication failed: %s", err)
			return fmt.Errorf("authentication failed")
		}
		cli.Output.Success("Authentication successful!")
		if outputFormat == "json" {
			printJSON(info)
		} else {
			fmt.Printf("User:  %s (%s)\n", info.User.Name, info.User.Email)
			fmt.Printf("Role:  %s\n", info.User.Role)
			fmt.Printf("Team:  %s\n", info.Team.Name)
		}
		return nil
	},
}

var authWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user identity",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		info, err := client.Auth.Info(getContext())
		if err != nil {
			return err
		}
		if outputFormat == "json" {
			printJSON(info)
		} else {
			fmt.Printf("%s <%s> [%s] @ %s\n", info.User.Name, info.User.Email, info.User.Role, info.Team.Name)
		}
		return nil
	},
}

func init() {
	authCmd.AddCommand(authTestCmd)
	authCmd.AddCommand(authWhoamiCmd)
}
