package cmd

import (
	"fmt"
	"net/http"

	"github.com/scottbrown/hacksplaining-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	addSubject string
	addMessage string
	addGroupID int
)

var addCmd = &cobra.Command{
	Use:   "add <email>",
	Short: "Add a user and send them an invitation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		if err := validateEmail(email); err != nil {
			return err
		}

		var req *api.AddUserRequest
		if addSubject != "" || addMessage != "" || addGroupID != 0 {
			req = &api.AddUserRequest{
				Subject: addSubject,
				Message: addMessage,
				GroupID: addGroupID,
			}
		}

		statusCode, err := client.AddUser(email, req)
		if err != nil {
			return err
		}

		switch statusCode {
		case http.StatusCreated:
			fmt.Printf("User %s created and invitation sent.\n", email)
		case http.StatusOK:
			fmt.Printf("User %s already exists.\n", email)
		}
		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addSubject, "subject", "", "Custom subject line for the invitation email")
	addCmd.Flags().StringVar(&addMessage, "message", "", "Custom message body for the invitation email")
	addCmd.Flags().IntVar(&addGroupID, "group-id", 0, "Group ID to assign the user to")
}
