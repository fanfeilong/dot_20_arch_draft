# d2a status And d2a-status

## Goal

Provide a reliable way to re-display the current workflow state at any time.

## `d2a status`

Add a command:

```text
d2a status
```

It should print:

- repository name
- repository path
- d2a path
- current stage
- last command
- last skill
- next step
- recent history summary

This command should be read-only.

## `d2a-status` Skill

Add a skill:

- `d2a-status`

Its purpose is to:

1. read `.d2a/state.json`
2. read recent entries from `.d2a/history.jsonl`
3. restate the current stage machine position
4. tell the user what the next best step is

## Why Both Are Needed

The command is useful for terminal-level checks.

The skill is useful inside Codex, where the user may want the current state restated in the same conversational surface as the next action.

## Status

Implemented:

1. `d2a status`
2. `.d2a/state.json`
3. `.d2a/history.jsonl`
4. `d2a-status` skill
5. automatic state updates from stage commands and `d2a skill-state`
