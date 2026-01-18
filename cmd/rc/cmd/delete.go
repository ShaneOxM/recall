package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a reminder",
	Long: `Delete a reminder by its ID.

Examples:
  rc delete 1768773271812-7727a989
  rc delete 1768773271812-7727a989 --backend apple`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	id := args[0]

	s, err := getStore()
	if err != nil {
		return fmt.Errorf("initializing store: %w", err)
	}

	if err := s.Delete(context.Background(), id); err != nil {
		return fmt.Errorf("deleting reminder: %w", err)
	}

	fmt.Printf("Deleted: %s\n", id)
	return nil
}
