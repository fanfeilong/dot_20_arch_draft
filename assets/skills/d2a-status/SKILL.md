---
name: d2a-status
description: Built-in d2a skill for d2a-status stage guidance and state updates.
---

# d2a-status

## Goal

Re-display the current d2a workflow state from `.d2a/state.json` and `.d2a/history.jsonl`.

## Instructions

1. Start by confirming the current repository context. Print the repo name, repo path, and d2a path before status review.
2. If the active repository is unknown, stop and ask the user which repository should be used.
3. After context is confirmed, call `d2a skill-state d2a-status --status started --phase analysis-generation --next-step "Read the latest state and recent history." --summary "Started status review."`.
4. Read `.d2a/state.json`.
5. Read recent entries from `.d2a/history.jsonl`.
6. Restate the current stage, last command, last skill, and next recommended step.
7. Keep the summary short and operational.
8. If the state file is missing or stale, tell the user which command should be run first.
9. When the status summary is complete, call `d2a skill-state d2a-status --status completed --phase analysis-generation --summary "Completed status review."`.

## Turn-End Continuation Rule

1. At the end of every learner-facing turn, call `d2a skill-state` to persist the latest phase and progress before finishing the reply.
2. If the current phase is still in progress, record `--status progress` with updated question or challenge index.
3. End the reply with: `继续请使用 $d2a-step`.
4. If the current phase just completed, still emit the same continuation line so the learner can resume from the next state via `d2a-step`.

## Output

- Current stage
- Last command
- Last skill
- Next step
- Recent history summary
