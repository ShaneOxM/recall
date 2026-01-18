package jsonl

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shaneoxm/recall/internal/protocol"
)

func TestStore_AddAndGet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()
	r := protocol.NewReminder("Test reminder")

	if err := store.Add(ctx, r); err != nil {
		t.Fatalf("failed to add reminder: %v", err)
	}

	got, err := store.Get(ctx, r.ID)
	if err != nil {
		t.Fatalf("failed to get reminder: %v", err)
	}

	if got.Title != r.Title {
		t.Errorf("expected title %q, got %q", r.Title, got.Title)
	}
}

func TestStore_List(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	r1 := protocol.NewReminder("First")
	r2 := protocol.NewReminder("Second")
	r2.AddTag("work")

	store.Add(ctx, r1)
	store.Add(ctx, r2)

	reminders, err := store.List(ctx, nil)
	if err != nil {
		t.Fatalf("failed to list reminders: %v", err)
	}

	if len(reminders) != 2 {
		t.Errorf("expected 2 reminders, got %d", len(reminders))
	}
}

func TestStore_ListWithTagFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	r1 := protocol.NewReminder("Work task")
	r1.AddTag("work")
	r2 := protocol.NewReminder("Personal task")
	r2.AddTag("personal")

	store.Add(ctx, r1)
	store.Add(ctx, r2)

	filter := &protocol.ListFilter{Tags: []string{"work"}}
	reminders, err := store.List(ctx, filter)
	if err != nil {
		t.Fatalf("failed to list reminders: %v", err)
	}

	if len(reminders) != 1 {
		t.Errorf("expected 1 reminder, got %d", len(reminders))
	}
	if reminders[0].Title != "Work task" {
		t.Errorf("expected 'Work task', got %q", reminders[0].Title)
	}
}

func TestStore_ListExcludesCompleted(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	r1 := protocol.NewReminder("Active")
	r2 := protocol.NewReminder("Done")
	r2.Complete()

	store.Add(ctx, r1)
	store.Add(ctx, r2)

	filter := &protocol.ListFilter{IncludeCompleted: false}
	reminders, err := store.List(ctx, filter)
	if err != nil {
		t.Fatalf("failed to list reminders: %v", err)
	}

	if len(reminders) != 1 {
		t.Errorf("expected 1 reminder, got %d", len(reminders))
	}
}

func TestStore_Complete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()
	r := protocol.NewReminder("Test")
	store.Add(ctx, r)

	if err := store.Complete(ctx, r.ID); err != nil {
		t.Fatalf("failed to complete reminder: %v", err)
	}

	got, _ := store.Get(ctx, r.ID)
	if !got.Completed {
		t.Error("expected reminder to be completed")
	}
}

func TestStore_Delete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()
	r := protocol.NewReminder("Test")
	store.Add(ctx, r)

	if err := store.Delete(ctx, r.ID); err != nil {
		t.Fatalf("failed to delete reminder: %v", err)
	}

	_, err = store.Get(ctx, r.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_Update(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()
	r := protocol.NewReminder("Original")
	store.Add(ctx, r)

	r.Title = "Updated"
	if err := store.Update(ctx, r); err != nil {
		t.Fatalf("failed to update reminder: %v", err)
	}

	got, _ := store.Get(ctx, r.ID)
	if got.Title != "Updated" {
		t.Errorf("expected title 'Updated', got %q", got.Title)
	}
}

func TestStore_ListWithDueFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()
	now := time.Now()

	r1 := protocol.NewReminder("Today")
	r1.SetDue(now)
	r2 := protocol.NewReminder("Next week")
	r2.SetDue(now.Add(7 * 24 * time.Hour))

	store.Add(ctx, r1)
	store.Add(ctx, r2)

	tomorrow := now.Add(24 * time.Hour)
	filter := &protocol.ListFilter{DueBefore: &tomorrow}
	reminders, err := store.List(ctx, filter)
	if err != nil {
		t.Fatalf("failed to list reminders: %v", err)
	}

	if len(reminders) != 1 {
		t.Errorf("expected 1 reminder, got %d", len(reminders))
	}
}

func TestStore_GetNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	_, err = store.Get(context.Background(), "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	reminders, err := store.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to list from empty store: %v", err)
	}

	if len(reminders) != 0 {
		t.Errorf("expected 0 reminders, got %d", len(reminders))
	}
}

func TestStore_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "reminders.jsonl")

	store, err := New(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	r := protocol.NewReminder("Test")
	if err := store.Add(context.Background(), r); err != nil {
		t.Fatalf("failed to add reminder: %v", err)
	}

	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
