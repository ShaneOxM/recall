package cmd

import (
	"fmt"

	"github.com/shaneoxm/recall/internal/adapters/apple"
	"github.com/shaneoxm/recall/internal/adapters/jsonl"
	"github.com/shaneoxm/recall/internal/adapters/todoist"
	"github.com/shaneoxm/recall/internal/config"
	"github.com/shaneoxm/recall/internal/protocol"
)

var (
	store       protocol.Store
	backendFlag string
)

func getStore() (protocol.Store, error) {
	if store != nil {
		return store, nil
	}

	backend := backendFlag
	if backend == "" {
		backend = "local" // default
	}

	var err error
	switch backend {
	case "local", "jsonl":
		cfg := config.Default()
		store, err = jsonl.New(cfg.DataPath())
	case "apple", "reminders":
		store = apple.New("Recall")
	case "todoist":
		cfg := config.Default()
		if cfg.TodoistToken == "" {
			return nil, fmt.Errorf("TODOIST_API_TOKEN environment variable not set")
		}
		store = todoist.New(cfg.TodoistToken, "")
	default:
		return nil, fmt.Errorf("unknown backend: %s (use: local, apple, todoist)", backend)
	}

	if err != nil {
		return nil, err
	}
	return store, nil
}
