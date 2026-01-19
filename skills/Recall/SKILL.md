---
name: Recall
description: Create reminders for the user via CLI. USE WHEN user mentions something to remember, a task for later, a follow-up needed, OR when you notice something the user should not forget.
---

# Recall

Create reminders that sync to the user's preferred system (local, Apple Reminders, or Todoist).

## Quick Actions

For most requests, use these patterns directly:

| User Says | Do This |
|-----------|---------|
| "remind me to..." | `rc add "..." --due <date>` |
| "I need to remember..." | `rc add "..." --note "context"` |
| "add a todo for..." | `rc add "..." --tag work` |
| "show my reminders" | `rc list` or `rc list --today` |
| "what's due?" | `rc list --today` |

## Workflow Routing

| Workflow | Trigger | File |
|----------|---------|------|
| **AddReminder** | "remind me", "don't forget", "todo" | `Workflows/AddReminder.md` |
| **ListReminders** | "show reminders", "what's due" | `Workflows/ListReminders.md` |
| **CompleteReminder** | "done with", "finished", "complete" | `Workflows/CompleteReminder.md` |

## CLI Reference

```bash
# Add reminders
rc add "Task description" --due <date> --note "context" --tag <category> --priority <level>

# Due dates: today, tomorrow, monday, friday, 2024-03-15
# Priority: low, medium, high

# List reminders
rc list                 # All pending
rc list --today         # Due today
rc list --tomorrow      # Due tomorrow
rc list --week          # Due this week
rc list --tag work      # By tag
rc list --ids           # Show IDs for complete/delete

# Complete/Delete
rc complete <id>
rc delete <id>

# Backend selection (default: local)
rc add "Task" --backend apple    # Apple Reminders
rc add "Task" --backend todoist  # Todoist
```

## Examples

**Example 1: User mentions a follow-up**
```
User: "oh I need to call the dentist tomorrow"
> rc add "Call the dentist" --due tomorrow
> Created reminder for tomorrow.
```

**Example 2: Work task with context**
```
User: "remind me to review that PR before the meeting"
> rc add "Review PR" --due today --note "Before standup meeting" --tag work
> Created work reminder for today.
```

**Example 3: Agent notices something**
```
Agent notices user mentioned birthday next week
> rc add "Mom's birthday" --due "next saturday" --note "User mentioned during session" --tag personal
> Created reminder.
```

## Proactive Use

Create reminders when you notice:
- User mentions dates or deadlines
- User says "I should..." or "I need to..."
- Tasks that shouldn't be forgotten
- Follow-ups from the current session

## Do NOT

- Create reminders for things already tracked in SaveContext issues
- Ask permission for obvious reminder requests
- Create duplicate reminders for the same thing
