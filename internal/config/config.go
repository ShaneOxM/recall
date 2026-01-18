package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultDataDir  = ".recall"
	DefaultDataFile = "reminders.jsonl"
)

type Config struct {
	DataDir      string
	DataFile     string
	TodoistToken string
}

func Default() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		DataDir:      filepath.Join(home, DefaultDataDir),
		DataFile:     DefaultDataFile,
		TodoistToken: os.Getenv("TODOIST_API_TOKEN"),
	}
}

func (c *Config) DataPath() string {
	return filepath.Join(c.DataDir, c.DataFile)
}
