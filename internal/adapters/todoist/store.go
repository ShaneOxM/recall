package todoist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shaneoxm/recall/internal/protocol"
)

const baseURL = "https://api.todoist.com/rest/v2"

var ErrNotFound = errors.New("task not found")

// Store implements protocol.Store using Todoist REST API.
type Store struct {
	token   string
	project string // optional project name, defaults to Inbox
	client  *http.Client
}

// New creates a new Todoist store with the given API token.
func New(token string, project string) *Store {
	return &Store{
		token:   token,
		project: project,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// todoistTask represents a Todoist task.
type todoistTask struct {
	ID          string   `json:"id"`
	Content     string   `json:"content"`
	Description string   `json:"description,omitempty"`
	Due         *dueDateObj `json:"due,omitempty"`
	Priority    int      `json:"priority,omitempty"` // 1=normal, 2, 3, 4=urgent
	Labels      []string `json:"labels,omitempty"`
	IsCompleted bool     `json:"is_completed"`
	CreatedAt   string   `json:"created_at"`
}

type dueDateObj struct {
	String   string `json:"string,omitempty"`
	Date     string `json:"date,omitempty"`
	Datetime string `json:"datetime,omitempty"`
}

type createTaskRequest struct {
	Content     string   `json:"content"`
	Description string   `json:"description,omitempty"`
	DueString   string   `json:"due_string,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

// Add creates a new task in Todoist.
func (s *Store) Add(ctx context.Context, reminder *protocol.Reminder) error {
	req := createTaskRequest{
		Content:     reminder.Title,
		Description: s.buildDescription(reminder),
		Labels:      reminder.Tags,
	}

	if reminder.Due != nil {
		req.DueString = reminder.Due.Format("2006-01-02 15:04")
	}

	if reminder.Priority > 0 {
		// Todoist: 1=normal, 4=urgent. Ours: 1=low, 3=high
		// Map: our 1->2, our 2->3, our 3->4
		req.Priority = reminder.Priority + 1
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	_, err = s.doRequest(ctx, "POST", "/tasks", body)
	return err
}

// Get retrieves a task by ID.
func (s *Store) Get(ctx context.Context, id string) (*protocol.Reminder, error) {
	data, err := s.doRequest(ctx, "GET", "/tasks/"+id, nil)
	if err != nil {
		return nil, err
	}

	var task todoistTask
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return s.toReminder(&task), nil
}

// List returns all active tasks from Todoist.
func (s *Store) List(ctx context.Context, filter *protocol.ListFilter) ([]*protocol.Reminder, error) {
	data, err := s.doRequest(ctx, "GET", "/tasks", nil)
	if err != nil {
		return nil, err
	}

	var tasks []todoistTask
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	var reminders []*protocol.Reminder
	for _, task := range tasks {
		r := s.toReminder(&task)
		if s.matchesFilter(r, filter) {
			reminders = append(reminders, r)
		}
	}

	return reminders, nil
}

// Update modifies an existing task.
func (s *Store) Update(ctx context.Context, reminder *protocol.Reminder) error {
	req := createTaskRequest{
		Content:     reminder.Title,
		Description: s.buildDescription(reminder),
		Labels:      reminder.Tags,
	}

	if reminder.Due != nil {
		req.DueString = reminder.Due.Format("2006-01-02 15:04")
	}

	if reminder.Priority > 0 {
		req.Priority = reminder.Priority + 1
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	_, err = s.doRequest(ctx, "POST", "/tasks/"+reminder.ID, body)
	return err
}

// Delete removes a task by ID.
func (s *Store) Delete(ctx context.Context, id string) error {
	_, err := s.doRequest(ctx, "DELETE", "/tasks/"+id, nil)
	return err
}

// Complete marks a task as completed.
func (s *Store) Complete(ctx context.Context, id string) error {
	_, err := s.doRequest(ctx, "POST", "/tasks/"+id+"/close", nil)
	return err
}

func (s *Store) doRequest(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode == 404 {
		return nil, ErrNotFound
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

func (s *Store) buildDescription(r *protocol.Reminder) string {
	var parts []string

	if r.Notes != "" {
		parts = append(parts, r.Notes)
	}

	if len(r.Links) > 0 {
		parts = append(parts, "")
		parts = append(parts, "Links:")
		for _, link := range r.Links {
			parts = append(parts, "- "+link)
		}
	}

	return strings.Join(parts, "\n")
}

func (s *Store) toReminder(task *todoistTask) *protocol.Reminder {
	r := &protocol.Reminder{
		ID:        task.ID,
		Title:     task.Content,
		Notes:     task.Description,
		Tags:      task.Labels,
		Completed: task.IsCompleted,
	}

	// Parse priority (Todoist 1=normal, 4=urgent -> ours 0=none, 3=high)
	if task.Priority > 1 {
		r.Priority = task.Priority - 1
	}

	// Parse due date
	if task.Due != nil {
		if task.Due.Datetime != "" {
			if t, err := time.Parse(time.RFC3339, task.Due.Datetime); err == nil {
				r.Due = &t
			}
		} else if task.Due.Date != "" {
			if t, err := time.Parse("2006-01-02", task.Due.Date); err == nil {
				r.Due = &t
			}
		}
	}

	// Parse created date
	if task.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, task.CreatedAt); err == nil {
			r.CreatedAt = t
		}
	}

	return r
}

func (s *Store) matchesFilter(r *protocol.Reminder, filter *protocol.ListFilter) bool {
	if filter == nil {
		return true
	}

	if !filter.IncludeCompleted && r.Completed {
		return false
	}

	if len(filter.Tags) > 0 {
		hasTag := false
		for _, ft := range filter.Tags {
			for _, rt := range r.Tags {
				if strings.EqualFold(ft, rt) {
					hasTag = true
					break
				}
			}
		}
		if !hasTag {
			return false
		}
	}

	if filter.Search != "" {
		search := strings.ToLower(filter.Search)
		if !strings.Contains(strings.ToLower(r.Title), search) &&
			!strings.Contains(strings.ToLower(r.Notes), search) {
			return false
		}
	}

	return true
}
