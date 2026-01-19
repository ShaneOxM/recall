# ListReminders Workflow

For viewing existing reminders.

## When to Use

- User asks "what's on my list?"
- User wants to see today's tasks
- Before adding a reminder (to avoid duplicates)
- During session wrap-up

## Command Pattern

```bash
rc list [filters]
```

## Filters

| Flag | Description |
|------|-------------|
| `--today` | Due today or overdue |
| `--tomorrow` | Due tomorrow |
| `--week` | Due within 7 days |
| `--tag <name>` | Filter by tag |
| `--all` | Include completed |
| `--ids` | Show IDs (for complete/delete) |
| `--backend` | Which backend to query |

## Examples

**All pending reminders:**
```bash
rc list
```

**Today's reminders:**
```bash
rc list --today
```

**Work tasks this week:**
```bash
rc list --week --tag work
```

**Get IDs for completion:**
```bash
rc list --ids
```

**From specific backend:**
```bash
rc list --backend todoist
rc list --backend apple
```

## Output Format

```
[ ] Pay rent (due: Fri)
[ ] Review PR #123 (due: Mon) #work
[ ] Call mom (due: Sat) #family
```

With `--ids`:
```
[1768773271812-7727a989] Pay rent (due: Fri)
[1768773271813-8828b990] Review PR #123 (due: Mon) #work
```

## Best Practices

1. Check list before adding to avoid duplicates
2. Use `--today` for daily review
3. Use `--ids` before complete/delete operations
