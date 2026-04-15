# State Files

## Goal

Define the minimal persisted state under `.d2a/`.

## Required Files

```text
.d2a/
  state.json
  history.jsonl
  target.json
  qa/
    <skill>.json
    <skill>.jsonl
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
  "last_skill": "d2a-arch-1-project-scope",
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

## `.d2a/qa/<skill>.json`

This file tracks the current confirmation-question cursor for one skill.

Suggested fields:

```json
{
  "skill": "d2a-arch-2-runtime-view",
  "stage": "architecture-in-progress",
  "phase": "confirmation-questions",
  "question_index": 2,
  "question_total": 4,
  "updated_at": "2026-01-01T00:00:00Z"
}
```

## `.d2a/qa/<skill>.jsonl`

This file appends one record per confirmation question turn:

- question prompt
- learner answer
- evaluation (`correct|partial|incorrect`)
- one-line explanation
- timestamp

## Separation Rule

`docs/architecture/*.md` is architecture output only.

Confirmation questions, answers, scoring, and turn-by-turn teaching traces must stay under `.d2a/qa/*`, not in architecture docs.

## Design Rule

`state.json` answers "where are we now?"

`history.jsonl` answers "how did we get here?"
