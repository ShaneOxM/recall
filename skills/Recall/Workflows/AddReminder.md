# AddReminder Workflow

For creating new reminders via the `rc` CLI.

## When to Use

- User explicitly asks to be reminded
- User mentions something they need to do later
- Agent notices a task or follow-up the user should remember
- During session wrap-up for pending items

## Command Pattern

```bash
rc add "Task description" [options]
```

## Options

| Flag | Description | Examples |
|------|-------------|----------|
| `--due` | Due date | `today`, `tomorrow`, `monday`, `friday`, `2024-03-15` |
| `--note` | Additional context | `"Check auth changes"` |
| `--tag` | Category | `work`, `personal`, `family`, `health` |
| `--priority` | Importance | `low`, `medium`, `high` |
| `--link` | Related URL | `"https://github.com/..."` |
| `--backend` | Storage | `local`, `apple`, `todoist` |

## Examples

**Basic reminder:**
```bash
rc add "Buy groceries"
```

**With due date:**
```bash
rc add "Pay rent" --due friday
rc add "Submit report" --due 2024-03-15
```

**With full context:**
```bash
rc add "Review PR #123" --due monday --note "Check auth changes" --link "https://github.com/..." --tag work
```

**High priority:**
```bash
rc add "Urgent client call" --due today --priority high --tag work
```

**To specific backend:**
```bash
rc add "Doctor appointment" --due tomorrow --backend apple
```

## Best Practices

1. **Be specific** - "Call dentist about cleaning" not "Call dentist"
2. **Add context** - Use `--note` for why/what to remember
3. **Tag consistently** - Use common tags: work, personal, family, health
4. **Set realistic dates** - Use relative dates when possible

## Do NOT

- Create vague reminders like "remember the thing"
- Set all reminders to high priority
- Create reminders for immediate tasks (just do them)
- Duplicate reminders already in the system
