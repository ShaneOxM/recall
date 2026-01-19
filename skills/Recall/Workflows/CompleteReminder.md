# CompleteReminder Workflow

For marking reminders as done or deleting them.

## When to Use

- User says they finished a task
- User wants to remove a reminder
- During cleanup of old/irrelevant reminders

## Commands

```bash
# Mark as completed (keeps history)
rc complete <id>

# Delete permanently
rc delete <id>
```

## Workflow

1. **Get the ID first:**
```bash
rc list --ids
```

2. **Complete or delete:**
```bash
rc complete 1768773271812-7727a989
# or
rc delete 1768773271812-7727a989
```

## Backend-Specific IDs

| Backend | ID Format | Example |
|---------|-----------|---------|
| local | timestamp-hex | `1768773271812-7727a989` |
| todoist | numeric | `9929728458` |
| apple | title string | `"Call mom"` |

**Examples by backend:**

```bash
# Local
rc list --ids
rc complete 1768773271812-7727a989

# Todoist
rc list --ids --backend todoist
rc complete 9929728458 --backend todoist

# Apple (uses title as ID)
rc list --ids --backend apple
rc complete "Call mom" --backend apple
```

## Complete vs Delete

| Action | Use When |
|--------|----------|
| `complete` | Task is done, keep history |
| `delete` | Wrong entry, no longer relevant |

## Best Practices

1. Always `list --ids` first to get correct ID
2. Use `complete` for finished tasks (maintains history)
3. Use `delete` for mistakes or irrelevant items
4. Specify `--backend` if not using local default
