# d2a Repository-Root Workflow

## Goal

Turn `d2a` from a skill installer into a repository-root architecture workflow initializer.

The workflow should help a user:

1. Analyze an open-source project with `d2a-*` skills.
2. Save the analysis as small, numbered architecture docs.
3. Derive a runnable mini implementation that preserves the project's core architecture idea.
4. Grow integration tests alongside the mini implementation.
5. Present the combined `docs + src + tests` output in a local HTML report.

## Product Direction

`d2a` should no longer initialize only `.codex/skills` and similar agent directories.

It should initialize a full workspace for architecture decomposition:

- `docs/` for the analysis and implementation narrative
- `src/` for the mini implementation
- `tests/` for incremental integration tests
- `report/` for a local Vue-based report app

## Scope

This design targets an incremental implementation, not the final full product.

The near-term goal is:

1. Create the lab directory structure.
2. Create starter templates in `docs/`, `src/`, `tests/`, and `report/`.
3. Keep the current multi-agent skill installation behavior.

The later goal is:

1. Add commands that write architecture docs automatically.
2. Add commands that derive and refine the mini implementation.
3. Add commands that generate and serve the local report.

## Current Position

The current implementation already provides deterministic stage scaffolding for analysis, mini derivation, testing, and reporting.
