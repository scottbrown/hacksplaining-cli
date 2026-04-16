package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <email>",
	Short: "Remove a user from your license",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		if err := client.RemoveUser(email); err != nil {
			return err
		}
		fmt.Printf("User %s removed.\n", email)
		return nil
	},
}
