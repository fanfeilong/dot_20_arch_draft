# CLI And Workflow

## Planned Commands

The current CLI already supports:

- `d2a help`
- `d2a init <target-dir>`
- `d2a version`

The staged CLI is:

1. `d2a init <repo-dir>`
2. `d2a analyze <target-repo> [--repo <repo-dir>]`
3. `d2a derive-mini [--repo <repo-dir>] [--skip-challenge-reason <text>]`
4. `d2a test-mini [--repo <repo-dir>]`
5. `d2a report [--repo <repo-dir>]`
6. `d2a serve [--repo <repo-dir>]`

## Workflow Meaning

### `d2a init`

Create the repository-root `.d2a` structure, starter docs, and built-in skills.

### `d2a analyze`

Fill `docs/architecture/` by applying the six-element decomposition:

- boundary
- driver
- core objects
- state evolution
- cooperation
- constraints

### `d2a derive-mini`

Use the architecture docs to choose the smallest runnable implementation that still demonstrates the target project's core architecture idea.

### `d2a test-mini`

Add or run incremental integration tests for the mini implementation under `tests/`.

### `d2a report`

Generate structured report data and report content based on `docs/`, `src/`, and `tests/`.

### `d2a serve`

Start the local report static server.

## Skill Roles

The current built-in skills should map to the analysis stage:

- `d2a-project-scope`
- `d2a-runtime-view`
- `d2a-core-objects`
- `d2a-state-evolution`
- `d2a-module-view`
- `d2a-tradeoff-view`
- `d2a-architecture-walkthrough`

Current implementation-focused skills:

- `d2a-mini-scope`
- `d2a-mini-design`
- `d2a-mini-build`
- `d2a-mini-test`
- `d2a-report-build`
