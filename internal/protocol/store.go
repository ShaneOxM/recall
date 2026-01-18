package protocol

import (
	"context"
	"time"
)

// Store defines the interface for reminder storage backends.
type Store interface {
	// Add creates a new reminder.
	Add(ctx context.Context, reminder *Reminder) error

	// Get retrieves a reminder by ID.
	Get(ctx context.Context, id string) (*Reminder, error)

	// List returns all reminders matching the filter.
	List(ctx context.Context, filter *ListFilter) ([]*Reminder, error)

	// Update modifies an existing reminder.
	Update(ctx context.Context, reminder *Reminder) error

	// Delete removes a reminder by ID.
	Delete(ctx context.Context, id string) error

	// Complete marks a reminder as completed.
	Complete(ctx context.Context, id string) error
}

// ListFilter specifies criteria for listing reminders.
type ListFilter struct {
	// IncludeCompleted includes completed reminders in results.
	IncludeCompleted bool

	// Tags filters to reminders with any of these tags.
	Tags []string

	// DueBefore filters to reminders due before this time.
	DueBefore *time.Time

	// DueAfter filters to reminders due after this time.
	DueAfter *time.Time

	// Search filters to reminders containing this text in title or notes.
	Search string
}
