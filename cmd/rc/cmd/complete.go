package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Mark a reminder as completed",
	Long: `Mark a reminder as completed by its ID.

Examples:
  rc complete 1768773271812-7727a989
  rc complete 1768773271812-7727a989 --backend apple`,
	Args: cobra.ExactArgs(1),
	RunE: runComplete,
}

func init() {
	rootCmd.AddCommand(completeCmd)
}

func runComplete(cmd *cobra.Command, args []string) error {
	id := args[0]

	s, err := getStore()
	if err != nil {
		return fmt.Errorf("initializing store: %w", err)
	}

	if err := s.Complete(context.Background(), id); err != nil {
		return fmt.Errorf("completing reminder: %w", err)
	}

	fmt.Printf("Completed: %s\n", id)
	return nil
}
