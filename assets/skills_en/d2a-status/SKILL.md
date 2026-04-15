---
name: d2a-status
description: Built-in d2a skill for d2a-status stage guidance and state updates.
---

# d2a-status

## Goal

Re-display the current d2a workflow state from `.d2a/state.json` and `.d2a/history.jsonl`.

## Required Output Format (Same Envelope as d2a-step)

All non-`d2a-step` skills must reuse the same envelope format as `d2a-step`:

```text
==================================================
[<Layer-1> <N/Total> | <Layer-2>] <repo>
next: <next_skill> -> <next_file>
==================================================

<body of this skill>

--------------------------------------------------
done: <what this turn completed>
state: <current skeleton position> -> <next skeleton position> · Continue with $d2a-step
--------------------------------------------------
```

Rules:

1. Keep the opening to exactly two lines between separators.
2. Keep the ending to exactly two lines between separators.
3. Do not print the old multi-line start header list anymore (repo/path/stage/flow as separate lines).
4. If `next_file` is unknown, print `unknown` and keep the same two-line shape.


## Body Formatting Hard Rules

1. The body must use `- ` bullet lines, minimum 2 lines and maximum 4 lines.
2. Each point must occupy exactly one line; do not write long paragraphs.
3. Keep each line under 80 characters; split lines when longer.
4. Avoid Markdown emphasis in body text (for example `` `...` `` or `**...**`).

If the active repository is unknown, stop and ask the user which repository should be used.

## Human In Loop Marker Rule

When the current turn asks the user a question and waits for user input, the last line of the reply body must append:

`[human_in_loop]`

## Instructions

1. Start by confirming the current repository context. Put repo/path information inside the shared envelope format instead of printing a separate header list.
2. If the active repository is unknown, stop and ask the user which repository should be used.
3. After context is confirmed, call `d2a skill-state d2a-status --status started --phase analysis-generation --next-step "Read the latest state and recent history." --summary "Started status review."`.
4. Read `.d2a/state.json`.
5. Read recent entries from `.d2a/history.jsonl`.
6. Restate the current stage, last command, last skill, and next recommended step.
7. Keep the summary short and operational.
8. If the state file is missing or stale, tell the user which command should be run first.
9. When the status summary is complete, call `d2a skill-state d2a-status --status completed --phase analysis-generation --summary "Completed status review."`.

## Turn-End Continuation Rule

1. Before ending each turn, call `d2a skill-state` to persist the current phase and progress.
2. End the reply with: `Continue with $d2a-step`.


## Output

- Current stage
- Last command
- Last skill
- Next step
- Recent history summary
