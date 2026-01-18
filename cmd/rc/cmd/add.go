package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/shaneoxm/recall/internal/protocol"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new reminder",
	Long: `Add a new reminder with optional context.

Examples:
  rc add "Call mom"
  rc add "Call mom" --due tomorrow --note "Birthday next week"
  rc add "Review PR" --due monday --link "https://github.com/..." --tag work
  rc add "Pay rent" --due friday --priority high`,
	Args: cobra.ExactArgs(1),
	RunE: runAdd,
}

var (
	addDue      string
	addNote     string
	addLinks    []string
	addTags     []string
	addPriority string
)

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addDue, "due", "d", "", "due date (e.g., tomorrow, monday, 2024-01-15)")
	addCmd.Flags().StringVarP(&addNote, "note", "n", "", "note or instructions")
	addCmd.Flags().StringSliceVarP(&addLinks, "link", "l", nil, "links (can be specified multiple times)")
	addCmd.Flags().StringSliceVarP(&addTags, "tag", "t", nil, "tags (can be specified multiple times)")
	addCmd.Flags().StringVarP(&addPriority, "priority", "p", "", "priority (low, medium, high)")
}

func runAdd(cmd *cobra.Command, args []string) error {
	title := args[0]
	reminder := protocol.NewReminder(title)

	if addNote != "" {
		reminder.SetNotes(addNote)
	}

	for _, link := range addLinks {
		reminder.AddLink(link)
	}

	for _, tag := range addTags {
		reminder.AddTag(tag)
	}

	if addPriority != "" {
		priority := parsePriority(addPriority)
		reminder.SetPriority(priority)
	}

	if addDue != "" {
		due, err := parseDue(addDue)
		if err != nil {
			return fmt.Errorf("invalid due date: %w", err)
		}
		reminder.SetDue(due)
	}

	s, err := getStore()
	if err != nil {
		return fmt.Errorf("initializing store: %w", err)
	}

	if err := s.Add(context.Background(), reminder); err != nil {
		return fmt.Errorf("saving reminder: %w", err)
	}

	fmt.Printf("Created reminder: %s (ID: %s)\n", reminder.Title, reminder.ID)
	if reminder.Due != nil {
		fmt.Printf("  Due: %s\n", reminder.Due.Format("Mon Jan 2, 2006 3:04 PM"))
	}
	if reminder.Notes != "" {
		fmt.Printf("  Note: %s\n", reminder.Notes)
	}
	if len(reminder.Links) > 0 {
		fmt.Printf("  Links: %v\n", reminder.Links)
	}
	if len(reminder.Tags) > 0 {
		fmt.Printf("  Tags: %v\n", reminder.Tags)
	}

	return nil
}

func parsePriority(s string) int {
	switch s {
	case "low", "1":
		return 1
	case "medium", "med", "2":
		return 2
	case "high", "3":
		return 3
	default:
		return 0
	}
}

func parseDue(s string) (time.Time, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())

	switch s {
	case "today":
		return today, nil
	case "tomorrow":
		return today.AddDate(0, 0, 1), nil
	case "monday", "mon":
		return nextWeekday(today, time.Monday), nil
	case "tuesday", "tue":
		return nextWeekday(today, time.Tuesday), nil
	case "wednesday", "wed":
		return nextWeekday(today, time.Wednesday), nil
	case "thursday", "thu":
		return nextWeekday(today, time.Thursday), nil
	case "friday", "fri":
		return nextWeekday(today, time.Friday), nil
	case "saturday", "sat":
		return nextWeekday(today, time.Saturday), nil
	case "sunday", "sun":
		return nextWeekday(today, time.Sunday), nil
	}

	// Try parsing as date
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"Jan 2",
		"January 2",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			// If year wasn't specified, use current year
			if t.Year() == 0 {
				t = t.AddDate(now.Year(), 0, 0)
			}
			return time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, now.Location()), nil
		}
	}

	return time.Time{}, fmt.Errorf("could not parse date: %s", s)
}

func nextWeekday(from time.Time, weekday time.Weekday) time.Time {
	days := int(weekday - from.Weekday())
	if days <= 0 {
		days += 7
	}
	return from.AddDate(0, 0, days)
}
