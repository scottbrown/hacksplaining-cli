package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	remindSubject string
	remindMessage string
)

var remindCmd = &cobra.Command{
	Use:   "remind <email>",
	Short: "Send a training reminder email to a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		if err := validateEmail(email); err != nil {
			return err
		}
		statusCode, err := client.RemindUser(email, remindSubject, remindMessage)
		if err != nil {
			return err
		}

		switch statusCode {
		case http.StatusOK:
			fmt.Printf("Reminder sent to %s.\n", email)
		}
		return nil
	},
}

func init() {
	remindCmd.Flags().StringVar(&remindSubject, "subject", "", "Custom subject line for the reminder email")
	remindCmd.Flags().StringVar(&remindMessage, "message", "", "Custom message body for the reminder email")
}
