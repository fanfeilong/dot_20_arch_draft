---
name: d2a-report-build
description: Built-in d2a skill for d2a-report-build stage guidance and state updates.
---

# d2a-report-build

## Goal

Assemble the current repository d2a workspace into a report-ready package for local presentation.

## Instructions

1. Start by confirming the current repository context. Print the repo name, repo path, and d2a path before report work.
2. If the active repository is unknown, stop and ask the user which repository should be used.
3. After context is confirmed, call `d2a skill-state d2a-report-build --status started --stage report-prepared --phase analysis-generation --next-step "Refine the report summary and report artifacts." --summary "Started report-build work."`.
4. Treat this skill as the user-facing entry for the reporting stage inside Codex.
5. If report data is missing or stale, call `d2a report` before refining the report.
6. Read `.d2a/report/index.md` and `.d2a/report/data/*.json`.
7. Use `.d2a/docs/`, `.d2a/src/`, and `.d2a/tests/` as the content sources behind the report.
8. Keep the report focused on architecture, mini implementation, tests, and the teaching narrative.
9. Treat `.d2a/report/data/*.json` as the stable input contract for the future Vue app.
10. Refine `.d2a/report/index.md` or future report assets without changing the stage contract.
11. When this pass is complete, call `d2a skill-state d2a-report-build --status completed --stage report-ready --phase analysis-generation --next-step "Review the local report or run d2a serve." --summary "Completed report-build work."`.

## Turn-End Continuation Rule

1. At the end of every learner-facing turn, call `d2a skill-state` to persist the latest phase and progress before finishing the reply.
2. If the current phase is still in progress, record `--status progress` with updated question or challenge index.
3. End the reply with: `继续请使用 $d2a-step`.
4. If the current phase just completed, still emit the same continuation line so the learner can resume from the next state via `d2a-step`.

## Output

- Report summary
- Report data alignment
- Next presentation improvements
