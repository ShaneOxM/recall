package apple

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/shaneoxm/recall/internal/protocol"
)

var ErrNotFound = errors.New("reminder not found")

// Store implements protocol.Store using Apple Reminders via osascript.
type Store struct {
	listName string
}

// New creates a new Apple Reminders store.
// If listName is empty, uses "Recall" as the default list.
func New(listName string) *Store {
	if listName == "" {
		listName = "Recall"
	}
	return &Store{listName: listName}
}

// Add creates a new reminder in Apple Reminders.
func (s *Store) Add(ctx context.Context, reminder *protocol.Reminder) error {
	script := s.buildAddScript(reminder)
	_, err := s.runScript(ctx, script)
	return err
}

// Get retrieves a reminder by ID (title match for Apple Reminders).
func (s *Store) Get(ctx context.Context, id string) (*protocol.Reminder, error) {
	// Apple Reminders doesn't expose IDs easily via AppleScript
	// We search by title as a workaround
	reminders, err := s.List(ctx, &protocol.ListFilter{Search: id})
	if err != nil {
		return nil, err
	}
	for _, r := range reminders {
		if r.ID == id || r.Title == id {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

// List returns all reminders from the Apple Reminders list.
func (s *Store) List(ctx context.Context, filter *protocol.ListFilter) ([]*protocol.Reminder, error) {
	script := fmt.Sprintf(`
tell application "Reminders"
	set output to ""
	try
		set reminderList to list "%s"
		repeat with r in reminders of reminderList
			set rName to name of r
			set rBody to body of r
			if rBody is missing value then set rBody to ""
			set rCompleted to completed of r
			set rDueDate to ""
			try
				set rDueDate to due date of r as string
			end try
			set rPriority to priority of r
			set output to output & rName & "|||" & rBody & "|||" & rCompleted & "|||" & rDueDate & "|||" & rPriority & "
"
		end repeat
	end try
	return output
end tell`, s.listName)

	out, err := s.runScript(ctx, script)
	if err != nil {
		return nil, err
	}

	return s.parseReminders(out, filter), nil
}

// Update modifies an existing reminder.
func (s *Store) Update(ctx context.Context, reminder *protocol.Reminder) error {
	// For simplicity, we delete and re-add
	if err := s.Delete(ctx, reminder.Title); err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	return s.Add(ctx, reminder)
}

// Delete removes a reminder by title.
func (s *Store) Delete(ctx context.Context, id string) error {
	script := fmt.Sprintf(`
tell application "Reminders"
	try
		set reminderList to list "%s"
		repeat with r in reminders of reminderList
			if name of r is "%s" then
				delete r
				return "deleted"
			end if
		end repeat
	end try
	return "not found"
end tell`, s.listName, escapeAppleScript(id))

	out, err := s.runScript(ctx, script)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) == "not found" {
		return ErrNotFound
	}
	return nil
}

// Complete marks a reminder as completed.
func (s *Store) Complete(ctx context.Context, id string) error {
	script := fmt.Sprintf(`
tell application "Reminders"
	try
		set reminderList to list "%s"
		repeat with r in reminders of reminderList
			if name of r is "%s" then
				set completed of r to true
				return "completed"
			end if
		end repeat
	end try
	return "not found"
end tell`, s.listName, escapeAppleScript(id))

	out, err := s.runScript(ctx, script)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) == "not found" {
		return ErrNotFound
	}
	return nil
}

func (s *Store) buildAddScript(r *protocol.Reminder) string {
	props := []string{fmt.Sprintf(`name:"%s"`, escapeAppleScript(r.Title))}

	if r.Notes != "" {
		props = append(props, fmt.Sprintf(`body:"%s"`, escapeAppleScript(r.Notes)))
	}

	if r.Due != nil {
		props = append(props, fmt.Sprintf(`due date:date "%s"`, r.Due.Format("January 2, 2006 3:04:05 PM")))
	}

	if r.Priority > 0 {
		// Apple priority: 0=none, 1=high, 5=medium, 9=low (inverse of ours)
		applePriority := map[int]int{1: 9, 2: 5, 3: 1}[r.Priority]
		props = append(props, fmt.Sprintf(`priority:%d`, applePriority))
	}

	return fmt.Sprintf(`
tell application "Reminders"
	try
		set reminderList to list "%s"
	on error
		make new list with properties {name:"%s"}
		set reminderList to list "%s"
	end try
	tell reminderList
		make new reminder with properties {%s}
	end tell
end tell`, s.listName, s.listName, s.listName, strings.Join(props, ", "))
}

func (s *Store) runScript(ctx context.Context, script string) (string, error) {
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("osascript error: %w: %s", err, string(out))
	}
	return string(out), nil
}

func (s *Store) parseReminders(output string, filter *protocol.ListFilter) []*protocol.Reminder {
	var reminders []*protocol.Reminder
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|||")
		if len(parts) < 5 {
			continue
		}

		r := &protocol.Reminder{
			ID:        parts[0], // Use title as ID for Apple Reminders
			Title:     parts[0],
			Notes:     parts[1],
			Completed: parts[2] == "true",
			CreatedAt: time.Now(), // Apple doesn't expose creation date easily
			UpdatedAt: time.Now(),
		}

		// Parse due date
		if parts[3] != "" {
			if t, err := time.Parse("Monday, January 2, 2006 at 3:04:05 PM", parts[3]); err == nil {
				r.Due = &t
			}
		}

		// Parse priority (convert from Apple's scale)
		switch parts[4] {
		case "1":
			r.Priority = 3 // high
		case "5":
			r.Priority = 2 // medium
		case "9":
			r.Priority = 1 // low
		}

		// Apply filters
		if filter != nil {
			if !filter.IncludeCompleted && r.Completed {
				continue
			}
			if filter.Search != "" {
				search := strings.ToLower(filter.Search)
				if !strings.Contains(strings.ToLower(r.Title), search) &&
					!strings.Contains(strings.ToLower(r.Notes), search) {
					continue
				}
			}
		}

		reminders = append(reminders, r)
	}

	return reminders
}

func escapeAppleScript(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
