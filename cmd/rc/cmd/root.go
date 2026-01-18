package cmd

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rc",
	Short: "Recall - Reminders with context",
	Long: `Recall is a CLI-first reminders tool with rich context support.

Create reminders with notes, links, and tags. Sync across backends
like Apple Reminders, Todoist, or local JSONL storage.

Examples:
  rc add "Call mom" --due tomorrow --note "Birthday next week"
  rc add "Review PR" --due monday --link "https://github.com/..." --tag work
  rc list --today
  rc list --tag work
  rc fetch "Call mom" --full`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Load env vars from ~/.recall/.env
	if home, err := os.UserHomeDir(); err == nil {
		godotenv.Load(filepath.Join(home, ".recall", ".env"))
	}

	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default $HOME/.recall/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&backendFlag, "backend", "b", "local", "storage backend (local, apple, todoist)")
}
