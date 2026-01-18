package protocol

import (
	"time"
)

// Reminder represents a reminder with rich context.
type Reminder struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Due         *time.Time `json:"due,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	Links       []string  `json:"links,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Priority    int       `json:"priority,omitempty"` // 0=none, 1=low, 2=medium, 3=high
	Completed   bool      `json:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewReminder creates a new reminder with the given title.
func NewReminder(title string) *Reminder {
	now := time.Now()
	return &Reminder{
		ID:        generateID(),
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Complete marks the reminder as completed.
func (r *Reminder) Complete() {
	now := time.Now()
	r.Completed = true
	r.CompletedAt = &now
	r.UpdatedAt = now
}

// AddLink adds a link to the reminder.
func (r *Reminder) AddLink(link string) {
	r.Links = append(r.Links, link)
	r.UpdatedAt = time.Now()
}

// AddTag adds a tag to the reminder.
func (r *Reminder) AddTag(tag string) {
	r.Tags = append(r.Tags, tag)
	r.UpdatedAt = time.Now()
}

// SetDue sets the due date for the reminder.
func (r *Reminder) SetDue(due time.Time) {
	r.Due = &due
	r.UpdatedAt = time.Now()
}

// SetNotes sets the notes for the reminder.
func (r *Reminder) SetNotes(notes string) {
	r.Notes = notes
	r.UpdatedAt = time.Now()
}

// SetPriority sets the priority for the reminder.
func (r *Reminder) SetPriority(priority int) {
	r.Priority = priority
	r.UpdatedAt = time.Now()
}
