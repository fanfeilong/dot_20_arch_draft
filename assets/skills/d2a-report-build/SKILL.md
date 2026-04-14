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

## Output

- Report summary
- Report data alignment
- Next presentation improvements
