package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var outputJSON bool

var usersCmd = &cobra.Command{
	Use:   "users [email]",
	Short: "List users or get a specific user's progress",
	Long:  "List all users and their training progress. Provide an email address to get details for a specific user.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return getUser(args[0])
		}
		return listUsers()
	},
}

func init() {
	usersCmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
}

func listUsers() error {
	users, err := client.ListUsers()
	if err != nil {
		return err
	}

	if outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(users)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "EMAIL\tPROGRESS\tCOMPLETE\tGROUP\tCREATED")
	for _, u := range users {
		group := ""
		if u.Group != nil {
			group = u.Group.Name
		}
		fmt.Fprintf(w, "%s\t%d%%\t%v\t%s\t%s\n", u.Email, u.Progress, u.Complete, group, u.CreatedAt)
	}
	return w.Flush()
}

func getUser(email string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	user, err := client.GetUser(email)
	if err != nil {
		return err
	}

	if outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(user)
	}

	fmt.Printf("Email:      %s\n", user.Email)
	fmt.Printf("Progress:   %d%%\n", user.Progress)
	fmt.Printf("Complete:   %v\n", user.Complete)
	if user.CompletionDate != nil {
		fmt.Printf("Completed:  %s\n", *user.CompletionDate)
	}
	fmt.Printf("Created:    %s\n", user.CreatedAt)
	if user.LastCommunication != nil {
		fmt.Printf("Last Email: %s\n", *user.LastCommunication)
	}
	if user.Group != nil {
		fmt.Printf("Group:      %s\n", user.Group.Name)
	}

	fmt.Println("\nExercises:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for exercise, status := range user.Exercises {
		fmt.Fprintf(w, "  %s\t%s\n", exercise, status)
	}
	return w.Flush()
}
