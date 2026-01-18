package jsonl

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/shaneoxm/recall/internal/protocol"
)

var ErrNotFound = errors.New("reminder not found")

// Store implements protocol.Store using JSONL file storage.
type Store struct {
	path string
	mu   sync.RWMutex
}

// New creates a new JSONL store at the given path.
func New(path string) (*Store, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating directory: %w", err)
	}

	return &Store{path: path}, nil
}

// Add creates a new reminder.
func (s *Store) Add(ctx context.Context, reminder *protocol.Reminder) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.appendReminder(reminder)
}

// Get retrieves a reminder by ID.
func (s *Store) Get(ctx context.Context, id string) (*protocol.Reminder, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reminders, err := s.readAll()
	if err != nil {
		return nil, err
	}

	for _, r := range reminders {
		if r.ID == id {
			return r, nil
		}
	}

	return nil, ErrNotFound
}

// List returns all reminders matching the filter.
func (s *Store) List(ctx context.Context, filter *protocol.ListFilter) ([]*protocol.Reminder, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reminders, err := s.readAll()
	if err != nil {
		return nil, err
	}

	if filter == nil {
		return reminders, nil
	}

	var result []*protocol.Reminder
	for _, r := range reminders {
		if s.matchesFilter(r, filter) {
			result = append(result, r)
		}
	}

	return result, nil
}

// Update modifies an existing reminder.
func (s *Store) Update(ctx context.Context, reminder *protocol.Reminder) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	reminders, err := s.readAll()
	if err != nil {
		return err
	}

	found := false
	for i, r := range reminders {
		if r.ID == reminder.ID {
			reminders[i] = reminder
			found = true
			break
		}
	}

	if !found {
		return ErrNotFound
	}

	return s.writeAll(reminders)
}

// Delete removes a reminder by ID.
func (s *Store) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	reminders, err := s.readAll()
	if err != nil {
		return err
	}

	newReminders := make([]*protocol.Reminder, 0, len(reminders))
	found := false
	for _, r := range reminders {
		if r.ID == id {
			found = true
			continue
		}
		newReminders = append(newReminders, r)
	}

	if !found {
		return ErrNotFound
	}

	return s.writeAll(newReminders)
}

// Complete marks a reminder as completed.
func (s *Store) Complete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	reminders, err := s.readAll()
	if err != nil {
		return err
	}

	for _, r := range reminders {
		if r.ID == id {
			r.Complete()
			return s.writeAll(reminders)
		}
	}

	return ErrNotFound
}

func (s *Store) appendReminder(reminder *protocol.Reminder) error {
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(reminder)
	if err != nil {
		return fmt.Errorf("marshaling reminder: %w", err)
	}

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("writing reminder: %w", err)
	}

	return nil
}

func (s *Store) readAll() ([]*protocol.Reminder, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var reminders []*protocol.Reminder
	seen := make(map[string]*protocol.Reminder)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var r protocol.Reminder
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			continue // Skip malformed lines
		}

		// Latest version wins (append-only semantics)
		seen[r.ID] = &r
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	for _, r := range seen {
		reminders = append(reminders, r)
	}

	return reminders, nil
}

func (s *Store) writeAll(reminders []*protocol.Reminder) error {
	tmpPath := s.path + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	for _, r := range reminders {
		data, err := json.Marshal(r)
		if err != nil {
			f.Close()
			os.Remove(tmpPath)
			return fmt.Errorf("marshaling reminder: %w", err)
		}
		if _, err := f.Write(append(data, '\n')); err != nil {
			f.Close()
			os.Remove(tmpPath)
			return fmt.Errorf("writing reminder: %w", err)
		}
	}

	if err := f.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err := os.Rename(tmpPath, s.path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("renaming temp file: %w", err)
	}

	return nil
}

func (s *Store) matchesFilter(r *protocol.Reminder, filter *protocol.ListFilter) bool {
	// Filter completed
	if !filter.IncludeCompleted && r.Completed {
		return false
	}

	// Filter by tags
	if len(filter.Tags) > 0 {
		hasTag := false
		for _, filterTag := range filter.Tags {
			for _, reminderTag := range r.Tags {
				if strings.EqualFold(reminderTag, filterTag) {
					hasTag = true
					break
				}
			}
			if hasTag {
				break
			}
		}
		if !hasTag {
			return false
		}
	}

	// Filter by due date
	if filter.DueBefore != nil && r.Due != nil {
		if r.Due.After(*filter.DueBefore) {
			return false
		}
	}
	if filter.DueAfter != nil && r.Due != nil {
		if r.Due.Before(*filter.DueAfter) {
			return false
		}
	}

	// Filter by search text
	if filter.Search != "" {
		search := strings.ToLower(filter.Search)
		if !strings.Contains(strings.ToLower(r.Title), search) &&
			!strings.Contains(strings.ToLower(r.Notes), search) {
			return false
		}
	}

	return true
}
