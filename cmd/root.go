package cmd

import (
	"fmt"
	"net/mail"
	"os"

	"github.com/scottbrown/hacksplaining-cli/internal/api"
	"github.com/scottbrown/hacksplaining-cli/internal/config"
	"github.com/spf13/cobra"
)

var client *api.Client

var rootCmd = &cobra.Command{
	Use:   "hacksplaining",
	Short: "CLI for the Hacksplaining API",
	Long:  "Interact with the Hacksplaining API to manage users and track security training progress.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "version" {
			return nil
		}
		apiKey, err := config.LoadAPIKey()
		if err != nil {
			return err
		}
		client = api.NewClient(apiKey)
		return nil
	},
	SilenceUsage: true,
}

func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email address %q", email)
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(remindCmd)
	rootCmd.AddCommand(versionCmd)
}
