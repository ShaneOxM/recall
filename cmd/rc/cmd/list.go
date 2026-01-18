package cmd

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/shaneoxm/recall/internal/protocol"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List reminders",
	Long: `List reminders with optional filtering.

Examples:
  rc list                   # List all pending reminders
  rc list --today           # List reminders due today
  rc list --tag work        # List reminders tagged "work"
  rc list --all             # Include completed reminders`,
	RunE: runList,
}

var (
	listToday     bool
	listTomorrow  bool
	listWeek      bool
	listTags      []string
	listAll       bool
	listCompleted bool
	listShowIDs   bool
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&listToday, "today", false, "show only reminders due today")
	listCmd.Flags().BoolVar(&listTomorrow, "tomorrow", false, "show only reminders due tomorrow")
	listCmd.Flags().BoolVar(&listWeek, "week", false, "show reminders due this week")
	listCmd.Flags().StringSliceVarP(&listTags, "tag", "t", nil, "filter by tags")
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "include completed reminders")
	listCmd.Flags().BoolVar(&listCompleted, "completed", false, "show only completed reminders")
	listCmd.Flags().BoolVar(&listShowIDs, "ids", false, "show reminder IDs (for complete/delete)")
}

func runList(cmd *cobra.Command, args []string) error {
	s, err := getStore()
	if err != nil {
		return fmt.Errorf("initializing store: %w", err)
	}

	filter := &protocol.ListFilter{
		IncludeCompleted: listAll || listCompleted,
		Tags:             listTags,
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := today.AddDate(0, 0, 1)

	if listToday {
		filter.DueAfter = &today
		filter.DueBefore = &endOfDay
	} else if listTomorrow {
		tomorrow := today.AddDate(0, 0, 1)
		endOfTomorrow := tomorrow.AddDate(0, 0, 1)
		filter.DueAfter = &tomorrow
		filter.DueBefore = &endOfTomorrow
	} else if listWeek {
		endOfWeek := today.AddDate(0, 0, 7)
		filter.DueAfter = &today
		filter.DueBefore = &endOfWeek
	}

	reminders, err := s.List(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("listing reminders: %w", err)
	}

	// Filter for completed-only if requested
	if listCompleted {
		var completed []*protocol.Reminder
		for _, r := range reminders {
			if r.Completed {
				completed = append(completed, r)
			}
		}
		reminders = completed
	}

	if len(reminders) == 0 {
		fmt.Println("No reminders found.")
		return nil
	}

	// Sort by due date (nil dates last), then by created date
	sort.Slice(reminders, func(i, j int) bool {
		if reminders[i].Due == nil && reminders[j].Due == nil {
			return reminders[i].CreatedAt.Before(reminders[j].CreatedAt)
		}
		if reminders[i].Due == nil {
			return false
		}
		if reminders[j].Due == nil {
			return true
		}
		return reminders[i].Due.Before(*reminders[j].Due)
	})

	for _, r := range reminders {
		printReminder(r, listShowIDs)
	}

	return nil
}

func printReminder(r *protocol.Reminder, showID bool) {
	status := "[ ]"
	if r.Completed {
		status = "[x]"
	}

	dueStr := ""
	if r.Due != nil {
		dueStr = fmt.Sprintf(" (due: %s)", r.Due.Format("Mon Jan 2"))
	}

	priorityStr := ""
	if r.Priority > 0 {
		priorities := []string{"", "!", "!!", "!!!"}
		priorityStr = " " + priorities[r.Priority]
	}

	fmt.Printf("%s %s%s%s\n", status, r.Title, dueStr, priorityStr)
	if showID {
		fmt.Printf("    ID: %s\n", r.ID)
	}

	if r.Notes != "" {
		fmt.Printf("    Note: %s\n", r.Notes)
	}
	if len(r.Tags) > 0 {
		fmt.Printf("    Tags: %v\n", r.Tags)
	}
	if len(r.Links) > 0 {
		for _, link := range r.Links {
			fmt.Printf("    Link: %s\n", link)
		}
	}
}
