---
name: d2a-report-build
description: Built-in d2a skill for d2a-report-build stage guidance and state updates.
---

# d2a-report-build

## Goal

Assemble the current repository d2a workspace into a report-ready package for local presentation.

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
3. After context is confirmed, call `d2a skill-state d2a-report-build --status started --stage report-prepared --phase analysis-generation --next-step "Refine the report summary and report artifacts." --summary "Started report-build work."`.
4. Treat this skill as the user-facing entry for the reporting stage inside Codex.
5. If report data is missing or stale, call `d2a report` before refining the report.
6. Read `.d2a/report/index.md` and `.d2a/report/data/*.json`.
7. Use `docs/`, `.d2a/src/`, and `.d2a/tests/` as the content sources behind the report.
8. Keep the report focused on architecture, mini implementation, tests, and the teaching narrative.
9. Treat `.d2a/report/data/*.json` as the stable input contract for the future Vue app.
10. You must generate two-page brief artifacts: `report/brief.md` and `report/brief.html` (A4 print style).
11. Before finishing, you must explicitly run `d2a report` once so the latest report artifacts are materialized on disk (not only described in conversation).
12. If any DoD item is missing, or `d2a report` has not succeeded, do not mark this skill as completed.
13. When this pass is complete, call `d2a skill-state d2a-report-build --status completed --stage report-ready --phase analysis-generation --next-step "Run d2a serve to open the report." --summary "Completed report-build work and refreshed artifacts via d2a report."`.

## DoD (All Required)

1. Strict two-page A4 structure:
   - Page 1: one state-machine/architecture diagram + compact six-element table (boundary/driver/core objects/state machine/module cooperation/constraints).
   - Page 2: mini implementation brief (stack, 20%% slice, build summary, test evidence, intentional omissions).
2. If content is too long, compress it; never spill to a third page.
3. Keep it teachable and concise; avoid long narrative paragraphs.
4. Required output files:
   - `report/brief.md`
   - `report/brief.html`
   - `report/index.html`

## Turn-End Continuation Rule

1. Before ending each turn, call `d2a skill-state` to persist the current phase and progress.
2. End the reply with: `Continue with $d2a-step`.


## Output (Artifacts)

- `report/brief.md` (structured 2-page A4 brief)
- `report/brief.html` (printable 2-page A4 brief)
- `report/index.md` (report index)

## Persistence (.d2a)

- Skill progression and next-step routing are persisted through `d2a skill-state` into `.d2a/state.json` and `.d2a/history.jsonl`.
