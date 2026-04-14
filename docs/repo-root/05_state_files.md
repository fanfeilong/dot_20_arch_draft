# State Files

## Goal

Define the minimal persisted state under `.d2a/`.

## Required Files

```text
.d2a/
  state.json
  history.jsonl
  target.json
```

## `state.json`

This file should represent the latest known state.

Suggested fields:

```json
{
  "repo_name": "n8n",
  "repo_path": "/abs/path/to/n8n",
  "d2a_path": "/abs/path/to/n8n/.d2a",
  "current_stage": "architecture-in-progress",
  "last_command": "d2a analyze",
  "last_skill": "d2a-architecture-walkthrough",
  "next_step": "Complete docs/architecture/00_overview.md through 99_code_map.md",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

## `history.jsonl`

This file should append an event for every significant command and skill action.

Each line should record:

- timestamp
- actor type: `command` or `skill`
- actor name
- stage before
- stage after
- short summary

## `target.json`

This file remains the stable target-repository reference.

It should stay focused on target identity, while `state.json` tracks workflow progress.

## Design Rule

`state.json` answers "where are we now?"

`history.jsonl` answers "how did we get here?"
