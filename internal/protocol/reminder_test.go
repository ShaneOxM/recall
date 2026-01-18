package protocol

import (
	"strings"
	"testing"
	"time"
)

func TestNewReminder(t *testing.T) {
	title := "Test reminder"
	r := NewReminder(title)

	if r.Title != title {
		t.Errorf("expected title %q, got %q", title, r.Title)
	}
	if r.ID == "" {
		t.Error("expected non-empty ID")
	}
	if r.Completed {
		t.Error("expected Completed to be false")
	}
	if r.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if r.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestReminder_Complete(t *testing.T) {
	r := NewReminder("Test")
	before := r.UpdatedAt

	time.Sleep(time.Millisecond)
	r.Complete()

	if !r.Completed {
		t.Error("expected Completed to be true")
	}
	if r.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
	if !r.UpdatedAt.After(before) {
		t.Error("expected UpdatedAt to be updated")
	}
}

func TestReminder_AddLink(t *testing.T) {
	r := NewReminder("Test")
	link := "https://example.com"

	r.AddLink(link)

	if len(r.Links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(r.Links))
	}
	if r.Links[0] != link {
		t.Errorf("expected link %q, got %q", link, r.Links[0])
	}
}

func TestReminder_AddTag(t *testing.T) {
	r := NewReminder("Test")
	tag := "work"

	r.AddTag(tag)

	if len(r.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(r.Tags))
	}
	if r.Tags[0] != tag {
		t.Errorf("expected tag %q, got %q", tag, r.Tags[0])
	}
}

func TestReminder_SetDue(t *testing.T) {
	r := NewReminder("Test")
	due := time.Now().Add(24 * time.Hour)

	r.SetDue(due)

	if r.Due == nil {
		t.Fatal("expected Due to be set")
	}
	if !r.Due.Equal(due) {
		t.Errorf("expected due %v, got %v", due, *r.Due)
	}
}

func TestReminder_SetNotes(t *testing.T) {
	r := NewReminder("Test")
	notes := "Important notes"

	r.SetNotes(notes)

	if r.Notes != notes {
		t.Errorf("expected notes %q, got %q", notes, r.Notes)
	}
}

func TestReminder_SetPriority(t *testing.T) {
	r := NewReminder("Test")

	r.SetPriority(3)

	if r.Priority != 3 {
		t.Errorf("expected priority 3, got %d", r.Priority)
	}
}

func TestGenerateID_Format(t *testing.T) {
	id := generateID()

	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		t.Errorf("expected ID format 'timestamp-hex', got %q", id)
	}
	if len(parts[1]) != 8 {
		t.Errorf("expected 8 char hex suffix, got %q", parts[1])
	}
}

func TestGenerateID_Unique(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateID()
		if ids[id] {
			t.Errorf("duplicate ID generated: %s", id)
		}
		ids[id] = true
	}
}
