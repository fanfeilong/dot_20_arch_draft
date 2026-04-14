# Repo-Root d2a Model

## Goal

Simplify `d2a` around a single rule:

One Git repository equals one `d2a` working directory.

## User Workflow

1. Clone a target repository.
2. Enter the repository root.
3. Run `d2a init`.
4. Enter Codex in that repository.
5. Use `d2a-*` skills to drive the rest of the workflow.

## Core Rule

After initialization, the current repository root is the only d2a context.

There is no separate workspace or multi-target switching model in this design.

## Why

This removes several sources of complexity:

- no separate workspace root lookup
- no active target switching
- no nested multi-target directory model
- no ambiguity about which repository is active

The repository itself is the unit of work.

## State Persistence

The workflow should be treated as a persistent state machine stored under `.d2a/`.

This state machine should record:

- current stage
- last command
- last skill
- next step
- event history

The purpose is to ensure that state is never lost across Codex sessions or interrupted runs.

## Current Status

The implementation now uses repository-root terminology in command surfaces:

- `--repo <repo-dir>`
- `d2a repo: ...`
- `d2a repo path: ...`
- `d2a path: ...`

Persistent state-machine support is active:

- `.d2a/state.json`
- `.d2a/history.jsonl`
- `d2a status`
- `d2a-status`
