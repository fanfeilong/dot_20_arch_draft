# Skill Flow

## End-To-End Flow In Codex

After `d2a init <repo-dir>`, the intended flow is:

1. User enters Codex in the target repository root.
2. User invokes an analysis skill.
3. The analysis skill calls `d2a analyze` if the repository has not yet been prepared for that target.
4. The analysis skill fills `docs/architecture/*.md`.
5. User invokes a mini-implementation skill.
6. The implementation skill calls `d2a derive-mini` if implementation task files are not yet prepared.
7. The implementation skill fills `docs/implementation/*.md` and later writes `src/`.
8. User invokes a testing skill.
9. The testing skill calls `d2a test-mini` if the tests stage is not yet prepared.
10. The testing skill fills `tests/*` and later writes runnable tests.
11. User invokes a report skill.
12. The report skill calls `d2a report` if report data is stale or missing.
13. The report skill fills or refines the final report output.

## Recommended Skill Entry Points

### Analysis

- `$d2a-architecture-walkthrough`
- `$d2a-project-scope`
- `$d2a-runtime-view`
- `$d2a-core-objects`
- `$d2a-state-evolution`
- `$d2a-module-view`
- `$d2a-tradeoff-view`

### Implementation

Planned next skills:

- `$d2a-mini-scope`
- `$d2a-mini-design`
- `$d2a-mini-build`

### Testing

Planned next skills:

- `$d2a-mini-test`

### Reporting

Planned next skills:

- `$d2a-report-build`

## File Ownership By Skill Stage

- Analysis skills own `docs/architecture/*`
- Implementation skills own `docs/implementation/*` and later `src/*`
- Testing skills own `tests/*`
- Reporting skills own `report/*`

## Success Condition

The user should feel that they are working through a set of skills, while `d2a` commands quietly keep repository state in the right structural stage underneath.
