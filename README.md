# Recall

CLI-first reminders tool with rich context support. Let AI agents create reminders for you.

```bash
rc add "Call mom" --due tomorrow --note "Birthday next week" --tag family
rc list --today
```

## Why Recall?

AI coding agents (Claude Code, Cursor, etc.) can now leave you reminders when your brain runs wild during a session. Instead of forgetting that random thought, the agent captures it.

Recall bridges to your existing reminder systems:
- **Local** - JSONL file at `~/.recall/reminders.jsonl` (git-syncable)
- **Apple Reminders** - Native macOS Reminders app via AppleScript
- **Todoist** - Todoist REST API

## Installation

### From Source

```bash
git clone https://github.com/shaneoxm/recall.git
cd recall
make install
```

### Go Install

```bash
go install github.com/shaneoxm/recall/cmd/rc@latest
```

## Usage

### Add Reminders

```bash
# Basic reminder
rc add "Buy groceries"

# With due date
rc add "Pay rent" --due friday
rc add "Call dentist" --due tomorrow
rc add "Submit report" --due 2024-03-15

# With context
rc add "Review PR" --due monday --note "Check auth changes" --link "https://github.com/..." --tag work

# With priority
rc add "Urgent task" --due today --priority high
```

### List Reminders

```bash
# All pending reminders
rc list

# Filter by time
rc list --today
rc list --tomorrow
rc list --week

# Filter by tag
rc list --tag work
rc list --tag family

# Include completed
rc list --all
```

### Backend Selection

```bash
# Local JSONL (default)
rc add "Task" --backend local

# Apple Reminders (macOS only)
rc add "Task" --backend apple

# Todoist
export TODOIST_API_TOKEN="your-token"
rc add "Task" --backend todoist
```

## Configuration

### Environment Variables

| Variable | Description |
|----------|-------------|
| `TODOIST_API_TOKEN` | Todoist API token from [Developer Settings](https://todoist.com/app/settings/integrations/developer) |

### Example Setup

```bash
# Add to your shell profile (~/.zshrc, ~/.bashrc)
export TODOIST_API_TOKEN="your-token-here"

# Optional: Set default backend
alias rc="rc --backend todoist"
```

## Agent Integration

### Claude Code

Add to your `CLAUDE.md` or system prompt:

```markdown
When the user mentions something they need to remember or do later,
use the `rc` command to create a reminder:

rc add "Task description" --due <date> --note "Context" --tag <category>
```

### MCP Server (Coming Soon)

Recall can run as an MCP server for native AI tool integration.

## Data Storage

### Local (JSONL)

Reminders are stored in `~/.recall/reminders.jsonl`:

```json
{"id":"123-abc","title":"Call mom","due":"2024-01-15T09:00:00Z","notes":"Birthday next week","tags":["family"],"completed":false,"created_at":"2024-01-14T10:00:00Z"}
```

The JSONL format is:
- Human readable
- Git-syncable
- Append-friendly
- Corruption resistant

### Apple Reminders

Creates a "Recall" list in Apple Reminders. Reminders sync via iCloud.

### Todoist

Tasks are created in your Todoist Inbox. Use the Todoist app to organize into projects.

## Development

```bash
# Build
make build

# Run tests
make test

# Install locally
make install
```

## License

MIT License - see [LICENSE](LICENSE)

## Related

- [SaveContext](https://github.com/smonastero/savecontext) - Context persistence for AI coding agents
