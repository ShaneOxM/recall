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

# Make available to AI agents (Claude Code, etc.)
sudo ln -sf ~/go/bin/rc /usr/local/bin/rc
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

# Show IDs (for complete/delete)
rc list --ids
```

### Complete & Delete

```bash
# Get IDs first
rc list --ids

# Mark as completed (local backend)
rc complete 1768773271812-7727a989

# Delete a reminder (local backend)
rc delete 1768773271812-7727a989

# With other backends
rc list --ids --backend todoist
rc complete 9929728458 --backend todoist
rc delete 9929728458 --backend todoist

rc list --ids --backend apple
rc complete "Call mom" --backend apple   # Apple uses title as ID
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

Create `~/.recall/.env` with your settings:

```bash
mkdir -p ~/.recall
echo 'TODOIST_API_TOKEN=your-token-here' > ~/.recall/.env
```

Get your Todoist token from [Developer Settings](https://todoist.com/app/settings/integrations/developer).

### Optional: Default Backend

```bash
# Add to your shell profile (~/.zshrc, ~/.bashrc)
alias rc="rc --backend todoist"
```

## Agent Integration

### Skill Package (Recommended)

Install the Recall skill so AI agents automatically know how to create reminders:

```bash
# Claude Code
cp -r skills/Recall ~/.claude/skills/

# Other agents - copy to your agent's skills directory
```

The skill teaches agents to:
- Create reminders when you mention tasks
- Use appropriate due dates and tags
- Proactively capture things you shouldn't forget

### Manual Setup

Alternatively, add to your `CLAUDE.md` or system prompt:

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

- [SaveContext](https://github.com/greenfieldlabs-inc/savecontext) - The OS for AI coding agents
